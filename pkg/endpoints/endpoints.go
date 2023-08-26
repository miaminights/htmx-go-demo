package endpoints

import (
	"htmx-go-demo/pkg/database"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type IndexPage struct {
	Films []database.Film
}

func IndexHandler(c echo.Context) error {
	var films []database.Film
	var err error

	films, err = database.GetFilms()

	if err != nil {
		c.Logger().Errorf("Unable to retrieve films: %+v", err)
		return c.String(http.StatusInternalServerError, "Unable to get films")
	}

	return c.Render(http.StatusOK, "index.html", IndexPage{
		Films: films,
	})
}

func FilmItemHandler(c echo.Context) error {
	// Get id from path
	id := c.Param("id")

	film, err := database.GetFilm(id)

	if err != nil {
		c.Logger().Errorf("Unable to retrieve film: %+v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "film-list-item", film)
}

func PostNewFilmHandler(c echo.Context) error {
	// Simulate request latency
	time.Sleep(1 * time.Second)

	// Get values from the form post request
	title := c.FormValue("title")
	director := c.FormValue("director")

	film := database.Film{
		Title:    title,
		Director: director,
		Id:       -1,
	}

	errors, err := film.Save()

	if err != nil {
		c.Logger().Errorf("Unable to save film: %+v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	if len(errors) > 0 {
		c.Logger().Errorf("missing required fields")
	}

	return c.Render(http.StatusOK, "film-list-item", film)

	// Get html template and specifically use the block named "film-list-item" to render the new film entry
	// tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// tmpl.ExecuteTemplate(w, "film-list-item", newFilmEntry)

	// Method for writing a new li manually without using the block syntax
	// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
	// tmpl, _ := template.New("t").Parse(htmlStr)
}

func DeleteFilmHandler(c echo.Context) error {
	// Simulate request latency
	time.Sleep(1 * time.Second)

	// Get id from path
	id := c.Param("id")

	err := database.DeleteFilm(id)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Send ok status code for Htmx to delete the element in the list
	return c.NoContent(http.StatusOK)
}

func PostFilmHandler(c echo.Context) error {
	// Simulate request latency
	time.Sleep(1 * time.Second)

	// Get id from path
	id, err := strconv.Atoi(c.Param("id"))

	if id == 0 {
		id = -1
	}

	// Get values from the form post request
	title := c.FormValue("title")
	director := c.FormValue("director")

	film := database.Film{
		Title:    title,
		Director: director,
		Id:       id,
	}

	errors, err := film.Save()

	if err != nil {
		c.Logger().Errorf("Unable to save film: %+v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	if len(errors) > 0 {
		c.Logger().Errorf("missing required fields")
	}

	// Render the updated film with the list element block
	return c.Render(http.StatusOK, "film-list-item", film)
}

func FilmEditFormHandler(c echo.Context) error {
	// Get id from path
	id := c.Param("id")
	film, err := database.GetFilm(id)

	if err != nil {
		c.Logger().Errorf("Unable to retrieve film: %+v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "edit-film-form", film)
}
