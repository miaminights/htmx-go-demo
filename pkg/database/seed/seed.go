package main

import (
	"htmx-go-demo/pkg/database"
	"log"
	"path/filepath"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var filmSeedData = []database.Film{
	{Title: "The Godfather", Director: "Francis Ford Coppola"},
	{Title: "Blade Runner", Director: "Ridley Scott"},
	{Title: "The Thing", Director: "John Carpenter"},
}

func main() {
	local_tmp_db, _ := filepath.Abs("tmp/films.db")
	db_url := "file://" + local_tmp_db

	err := database.Init(db_url)

	if err != nil {
		log.Fatalf("could not initialize db: %+v", err)
	}

	for i := range filmSeedData {
		database.Db.Exec(
			"INSERT INTO films (title, director) VALUES (?, ?) ",
			filmSeedData[i].Title,
			filmSeedData[i].Director,
		)
	}
}
