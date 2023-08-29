package template

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ParseTemplates() (*template.Template, error) {
	// cleanRoot := filepath.Clean("")

	tmpl := template.New("")

	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error when walking templates directory: %+v", err)
			return err
		}

		if strings.Contains(path, ".html") {
			log.Println(path)
			_, err := tmpl.ParseFiles(path)

			if err != nil {
				log.Printf("could not parse tempaltes: %+v", err)
				return err
			}
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
