package database

import (
	"fmt"
)

type Film struct {
	Title    string
	Director string
	Id       int
}

type ErrorMap = map[string]string

func (f *Film) Validate() ErrorMap {
	var errors ErrorMap = make(ErrorMap)

	if f.Title == "" {
		errors["title"] = "Title is required"
	}

	if f.Director == "" {
		errors["director"] = "Director is required"
	}

	return errors
}

func GetFilms() ([]Film, error) {
	res, err := Db.Query("SELECT * from films")

	if err != nil {
		return nil, fmt.Errorf("Unable to query films: %+v", err)
	}

	var films []Film = make([]Film, 0)

	for res.Next() {
		var id int
		var title string
		var director string

		err := res.Scan(&id, &title, &director)

		if err != nil {
			return nil, fmt.Errorf("Unable to scan db row: %+v", err)
		}

		films = append(films, Film{
			Title:    title,
			Director: director,
			Id:       id,
		})
	}

	res.Close()

	return films, nil
}

func GetFilm(id string) (*Film, error) {
	res, err := Db.Query("SELECT * from films WHERE id = $1", id)

	if err != nil {
		return nil, fmt.Errorf("Unable to query film: %+v", err)
	}

	res.Next()
	var _id int
	var title string
	var director string

	err = res.Scan(&_id, &title, &director)

	if err != nil {
		return nil, fmt.Errorf("Unable to scan db row: %+v", err)
	}

	res.Close()

	return &Film{
		Title:    title,
		Director: director,
		Id:       _id,
	}, nil
}

func DeleteFilm(id string) error {
	_, err := Db.Exec("DELETE from films WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("Unable to delete film: %+v", err)
	}

	return nil
}

func (f *Film) Save() (ErrorMap, error) {
	errors := f.Validate()

	if len(errors) > 0 {
		return errors, nil
	}

	var err error

	if f.Id == -1 {
		_, err = Db.Exec(`INSERT INTO films (title, director) VALUES ($1, $2)`, f.Title, f.Director)
	} else {
		_, err = Db.Exec(`UPDATE films SET title = $1, director = $2 WHERE id = $3`, f.Title, f.Director, f.Id)
	}

	return errors, err
}
