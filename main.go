package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Film struct {
	Title    string
	Director string
	Id       int
}

func main() {
	fmt.Println("Go app...")

	films := map[string]map[string]Film{
		"Films": {
			"1": {Title: "The Godfather", Director: "Francis Ford Coppola", Id: 1},
			"2": {Title: "Blade Runner", Director: "Ridley Scott", Id: 2},
			"3": {Title: "The Thing", Director: "John Carpenter", Id: 3},
		},
	}

	// handler function #1 - returns the index.html template, with film data
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, films)
	}

	// handler function #2 - returns the template block with the newly added film, as an HTMX response
	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		filmLength := len(films["Films"])
		newId := filmLength + 1
		newIdKey := fmt.Sprint(newId)
		newFilmEntry := Film{Title: title, Director: director, Id: newId}
		films["Films"][newIdKey] = newFilmEntry
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", newFilmEntry)
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		id := r.URL.Query().Get("id")

		// Remove film from Films list in memory
		delete(films["Films"], id)

		// Send ok status code for Htmx to delete the element in the list
		w.WriteHeader(200)
	}

	// define handlers
	http.HandleFunc("/", h1)
	http.HandleFunc("/add-film", h2)
	http.HandleFunc("/remove-film", h3)

	// Handle static files (CSS, JS, Images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8000", nil))

}
