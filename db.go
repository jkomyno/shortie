package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type ShortenedUrl struct {
	ShortUrl string
	Url      string
}

func loadFromShortUrl(shortUrl string) (*ShortenedUrl, error) {
    // Convert the shortUrl to a primary key
    pk := decodeURL(shortUrl)

    // Setup the database connection
    db, err := sql.Open("mysql", config.ConnectionString)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    // Prepare the statement for reading data
    stmtOut, err := db.Prepare("SELECT target FROM urls WHERE id = ?")
    if err != nil {
        return nil, err
    }
    defer stmtOut.Close()

    // We're going to store the result here
    u := &ShortenedUrl{ShortUrl: shortUrl}

    // Perform the query on the DB
    err = stmtOut.QueryRow(pk).Scan(&u.Url)
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

    // See if we already have this url in the DB
    u2, err := loadFromUrl(u.Url)
    if err == nil {
        // We found it so return it
        u.ShortUrl = u2.ShortUrl
        return nil
    }

    // Insert into the DB
    res, err := db.Exec("INSERT INTO urls (target) VALUES (?)", u.Url)

    // Get the inserted ID
    id, err := res.LastInsertId()
    if err != nil {
        return err
    }

    // Generate the shortURL from the ID
    u.ShortUrl = generateURL(id)
    return nil
}
