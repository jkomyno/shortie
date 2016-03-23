package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "math/rand"
)

type ShortenedUrl struct {
	ShortUrl string
	Url      string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func loadFromShortUrl(shortUrl string) (*ShortenedUrl, error) {
    // Setup the database connection
    db, err := sql.Open("mysql", config.ConnectionString)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    // Prepare the statement for reading data
    stmtOut, err := db.Prepare("SELECT target FROM urls WHERE alias = ?")
    if err != nil {
        return nil, err
    }
    defer stmtOut.Close()

    // We're going to store the result here
    u := &ShortenedUrl{ShortUrl: shortUrl}

    // Perform the query on the DB
    err = stmtOut.QueryRow(shortUrl).Scan(&u.Url)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func loadFromUrl(url string) (*ShortenedUrl, error) {
    // Setup the database connection
    db, err := sql.Open("mysql", config.ConnectionString)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    // Prepare the statement for reading data
    stmtOut, err := db.Prepare("SELECT id FROM urls WHERE target = ?")
    if err != nil {
        return nil, err
    }
    defer stmtOut.Close()

    // We're going to store the result here
    var id int64

    // Perform the query on the DB
    err = stmtOut.QueryRow(url).Scan(&id)
    if err != nil {
        return nil, err
    }

    // Encode the url
    return &ShortenedUrl{Url: url, ShortUrl: generateURL(id)}, nil
}

func (u *ShortenedUrl) save() error {
    // Setup the database connection
    db, err := sql.Open("mysql", config.ConnectionString)
    if err != nil {
        return err
    }
    defer db.Close()

    // Insert into the DB
    res, err := db.Exec("INSERT INTO urls (target, alias) VALUES (?, ?)", u.Url, u.ShortUrl)
    if err != nil {
        return err
    }
    if res == nil {
        return nil
    }
    return nil
}
