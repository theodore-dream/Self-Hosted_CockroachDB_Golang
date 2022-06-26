package main

import (
    "context"
    "log"
    "fmt"
    "github.com/jackc/pgx/v4"
)

type Album struct {
    ID     int64
    Title  string
    Artist string
    Price  float32
}

func main() {
	conn, err := pgx.Connect(context.Background(), "postgresql://root@10.166.133.153:26257/album?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close(context.Background())
    fmt.Print("Complete")

albums, err := albumsByArtist("John Coltrane")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Albums found: %v\n", albums)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
    // An albums slice to hold data from returned rows.
    var albums []Album

    rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
    if err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    defer rows.Close()
    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var alb Album
        if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        albums = append(albums, alb)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    return albums, nil
}


