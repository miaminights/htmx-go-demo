package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Film struct {
	Title    string
	Director string
	Id       string
}

func main() {
	fmt.Println("Go app...")

	films := map[string]map[string]Film{
		"Films": {
			"e9fb18e4-f501-4ca7-a3a0-6382c938ca28": {Title: "The Godfather", Director: "Francis Ford Coppola", Id: "e9fb18e4-f501-4ca7-a3a0-6382c938ca28"},
			"1d1a7cf4-4127-4f3c-bd6b-4769fd024cc5": {Title: "Blade Runner", Director: "Ridley Scott", Id: "1d1a7cf4-4127-4f3c-bd6b-4769fd024cc5"},
			"c309a87b-0918-488b-8d24-b5d3e439fb53": {Title: "The Thing", Director: "John Carpenter", Id: "c309a87b-0918-488b-8d24-b5d3e439fb53"},
		},
	}

	// handler function #1 - returns the index.html template, with film data
	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, films)
	}

	// handler function #2 - returns the template block with the newly added film, as an HTMX response
	postNewFilmHandler := func(w http.ResponseWriter, r *http.Request) {
		// Simulate request latency
		time.Sleep(1 * time.Second)

		// Get values from the form post request
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		// Create new random id for a new Film entry
		newId := uuid.New().String()

		newFilmEntry := Film{Title: title, Director: director, Id: newId}

		// Add new film entry to Films in memory
		films["Films"][newId] = newFilmEntry

		// Get html template and specifically use the block named "film-list-element" to render the new film entry
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", newFilmEntry)

		// Method for writing a new li manually without using the block syntax
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
	}

	// handler function #3 - deletes film from list based on the id
	deleteFilmHandler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		// Get id from query param
		id := r.URL.Query().Get("id")

		// Remove film from Films list in memory
		delete(films["Films"], id)

		// Send ok status code for Htmx to delete the element in the list
		w.WriteHeader(200)
	}

	// define handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add-film", postNewFilmHandler)
	http.HandleFunc("/remove-film", deleteFilmHandler)

	// Handle static files (CSS, JS, Images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8000", nil))

}
