package main

import (
	"io"
	"net/http"
	"regexp"
	"log"
	"os"
	"github.com/jessevdk/go-flags"
	"github.com/fluent/fluent-logger-golang/fluent"
)

var config struct {
	DefaultURL string `short:"d" long:"default-url" required:"true" env:"DEFAULT_URL"`
	BaseURL string `short:"u" long:"base-url" required:"true" env:"BASE_URL"`
	BindingPort string `short:"p" long:"port" required:"true" env:"PORT"`
	ConnectionString string `short:"c" long:"db" required:"true" env:"CONNECTION_STRING"`
	AuthenticationSecret string `short:"s" long:"secret" required:"true" env:"AUTH_SECRET"`
	FluentPort int `long:"fluent-port" env:"FLUENT_PORT"`
	FluentHost string `long:"fluent-host" env:"FLUENT_HOST"`
	FluentTag string `long:"fluent-tag" env:"FLUENT_TAG"`
}

var validURL = regexp.MustCompile("(https?://([-\\w\\.]+)+(:\\d+)?(/([\\w/_\\.]*(\\?\\S+)?)?)?)")

/*
 * Logging
 */

var (
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func InitLogging(
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime)
}

func LogRequest(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        Info.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
        handler.ServeHTTP(w, r)
    })
}

/*
 * Beacon
 */

func logPageView(u *ShortenedUrl) {
	// Setup the logger
	logger, err := fluent.New(fluent.Config{FluentPort: config.FluentPort, FluentHost: config.FluentHost})
  	if err != nil {
    	Warning.Println("Failed to connect to Fluent")
		return
  	}
	defer logger.Close()

	// Build the log data
	var data = map[string]string{
    	"a":  "3",
    	"u": u.Url,
		"s": u.ShortUrl,
  	}

	// Send to fluent
  	error := logger.Post(config.FluentTag, data)
  	if error != nil {
    	Warning.Printf("Failed to send data to fluentd")
  	}
}

/*
 * Web Handlers
 */

func rootHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	// Validate the url
	if !validURL.MatchString(url) {
		// Return a bad request
		Warning.Printf("%s Bad Input: %s\n", r.RemoteAddr, url)
		http.Error(w, "Bad URL", http.StatusBadRequest)
		return
	}

	// Authenticate the request
	if ! authenticateRequest(r, url) {
		// Return an unauthorized request
		Warning.Printf("%s Unauthorized: %s\n", r.RemoteAddr, url)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Save out and return
	u := &ShortenedUrl{Url: url}
	err := u.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, config.BaseURL+u.ShortUrl, http.StatusCreated)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[len("/"):]
	url := r.FormValue("url")

	if shortUrl == "" && len(url) > 0 {
		rootHandler(w, r)
		return
	}

	u, err := loadFromShortUrl(shortUrl)

	if err != nil {
		http.Redirect(w, r, config.DefaultURL, http.StatusFound)
		return
	}

	// Log the page view
	logPageView(u)

	// Redirect to target url
	http.Redirect(w, r, u.Url, http.StatusFound)
}

func main() {
	// Load configuration
	_, err := flags.Parse(&config)
	if err != nil {
    	os.Exit(1)
	}

	// Setup logging
	InitLogging(os.Stdout, os.Stdout, os.Stderr)
	Info.Println("Listening on port :"+config.BindingPort)

	// Start the web server
	mainHandler := http.HandlerFunc(viewHandler)
	http.Handle("/", LogRequest(mainHandler))
	http.ListenAndServe(":"+config.BindingPort, nil)
}
