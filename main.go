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
	films := map[string]map[string]map[string]string{
		"Films": {
			"e9fb18e4-f501-4ca7-a3a0-6382c938ca28": {"Title": "The Godfather", "Director": "Francis Ford Coppola", "Id": "e9fb18e4-f501-4ca7-a3a0-6382c938ca28"},
			"1d1a7cf4-4127-4f3c-bd6b-4769fd024cc5": {"Title": "Blade Runner", "Director": "Ridley Scott", "Id": "1d1a7cf4-4127-4f3c-bd6b-4769fd024cc5"},
			"c309a87b-0918-488b-8d24-b5d3e439fb53": {"Title": "The Thing", "Director": "John Carpenter", "Id": "c309a87b-0918-488b-8d24-b5d3e439fb53"},
		},
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, films)
	}

	postNewFilmHandler := func(w http.ResponseWriter, r *http.Request) {
		// Simulate request latency
		time.Sleep(1 * time.Second)

		// Get values from the form post request
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		// Create new random id for a new Film entry
		newId := uuid.New().String()

		newFilmEntry := map[string]string{"Title": title, "Director": director, "Id": newId}

		// Add new film entry to Films in memory
		films["Films"][newId] = newFilmEntry

		// Get html template and specifically use the block named "film-list-element" to render the new film entry
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", newFilmEntry)

		// Method for writing a new li manually without using the block syntax
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
	}

	deleteFilmHandler := func(w http.ResponseWriter, r *http.Request) {
		// Simulate request latency
		time.Sleep(1 * time.Second)
		// Get id from query param
		id := r.URL.Query().Get("id")

		// Remove film from Films list in memory
		delete(films["Films"], id)

		// Send ok status code for Htmx to delete the element in the list
		w.WriteHeader(200)
	}

	updateFilmHandler := func(w http.ResponseWriter, r *http.Request) {
		// Simulate request latency
		time.Sleep(1 * time.Second)

		// Get id from query param
		id := r.URL.Query().Get("id")

		// Get values from the form post request
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		film := films["Films"][id]

		// Modify fields with form data
		film["Title"] = title
		film["Director"] = director

		// Render the updated film with the list element block
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", film)
	}

	filmEditFormHandler := func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		film := films["Films"][id]
		htmlStr := fmt.Sprintf(`
			<li id="film-%[1]s" hx-swap="outerHTML" class="list-group-item bg-primary text-white">
				<form id='edit-film' hx-put='/update-film?id=%[1]s' hx-target='#film-%[1]s'>
						<div>
								<label for='film-title'>Title</label>
								<input type='text' name='title' id='film-title' class='form-control' value='%[2]s' />
						</div>
						<div>
								<label for='film-director'>Director</label>
								<input type='text' name='director' id='film-director' class='form-control' value='%[3]s' />
						</div>
						<button type='submit' class='btn btn-light'>
								Submit
								<span class='spinner-border spinner-border-sm edit-film-indicator' id='edit-film-spinner' role='status' aria-hidden='true'></span>
						</button>
						<button hx-get='/get-film-item?id=%[1]s' hx-trigger='click' hx-target='#film-%[1]s' class='btn btn-secondary'>
								Cancel
						</button>
				</form>
			</li>`, id, film["Title"], film["Director"])
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}

	filmItemHandler := func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		film := films["Films"][id]

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", film)
	}

	// define handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add-film", postNewFilmHandler)
	http.HandleFunc("/remove-film", deleteFilmHandler)
	http.HandleFunc("/update-film", updateFilmHandler)
	http.HandleFunc("/film-edit-form", filmEditFormHandler)
	http.HandleFunc("/get-film-item", filmItemHandler)

	// Handle static files (CSS, JS, Images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8000", nil))

}
