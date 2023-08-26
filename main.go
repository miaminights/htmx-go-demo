package main

import (
	"html/template"
	"htmx-go-demo/pkg/database"
	"htmx-go-demo/pkg/endpoints"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func main() {
	db_url := os.Getenv("DB_URL")

	if db_url == "" {
		local_tmp_db, _ := filepath.Abs("tmp/films.db")
		db_url = "file://" + local_tmp_db
	}

	err := database.Init(db_url)

	if err != nil {
		log.Fatalf("could not initialize db: %+v", err)
	}

	tmpl, err := template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatalf("could not initialize templates: %+v", err)
	}

	e := echo.New()
	e.Renderer = endpoints.NewTemplateRenderer(tmpl)

	// Handle static files (CSS, JS, Images)
	e.Static("/static", "static")

	// define handlers
	e.GET("/", endpoints.IndexHandler)
	e.POST("/add-film", endpoints.PostFilmHandler)
	e.PUT("/update-film/:id", endpoints.PostFilmHandler)
	e.DELETE("/remove-film/:id", endpoints.DeleteFilmHandler)
	e.GET("/film-edit-form/:id", endpoints.FilmEditFormHandler)
	e.GET("/get-film-item/:id", endpoints.FilmItemHandler)

	e.Logger.Fatal(e.Start(":8000"))
}
