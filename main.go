/*
 * GitHub: github.com/qianjinxu
 * Email: xuqianjinchn@gmail.com
 * Bio: https://jin.bio
 */

package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		// User:                 "",
		// Passwd:               "",
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "192.168.2.1:43306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}
	var err error
	// db, err = sql.Open("mysql", "username:password@tcp(192.168.2.1:43306)/recordings")
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("MySQL Connected!")

	selectArtist, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", selectArtist)

	selectID, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", selectID)

	insertID, err := addAlbum(Album{
		Title:  "DJ",
		Artist: "Qianjin Xu",
		Price:  999.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Insert ID: %v\n", insertID)
}

func albumsByArtist(artist string) ([]Album, error) {
	stmt, err := db.Prepare("SELECT * FROM album WHERE artist = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(artist)
	if err != nil {
		return nil, fmt.Errorf("Artist: %q %v", artist, err)
	}
	defer rows.Close()
	var albums []Album
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("Artist: %q %v", artist, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Artist: %q %v", artist, err)
	}
	return albums, nil
}

func albumByID(id int64) (Album, error) {
	stmt, err := db.Prepare("SELECT * FROM album WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var alb Album
	if err := stmt.QueryRow(id).Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("ID: %d (Empty set)", id)
		}
		return alb, fmt.Errorf("ID: %d (%v)", id, err)
	}
	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("Insert ID: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Insert ID: %v", err)
	}
	return id, nil
}
