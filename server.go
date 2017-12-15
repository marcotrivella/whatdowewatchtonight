package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

const (
	TEMPLATE_PATH    = "static/html"
	CREDENTIALS_PATH = "static/data/credentials.json"
	MAX_FILM_PAGE    = 50
)

type Context struct {
	Movies []Movie
}

var db *sqlx.DB
var mappaCategorie map[string][]Movie
var contatoreSeed = int64(0)

func main() {
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.HandleFunc("/", templateHandler)

	INIT_DB()

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func templateHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		serveTemplate(w, req, "index", Context{})
	} else if req.URL.Path == "/" {

	} else {
		return
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request, file_name string, context Context) {
	lp := filepath.Join(TEMPLATE_PATH, "layout.html")
	fp := filepath.Join(TEMPLATE_PATH, file_name+".html")

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout.html", context); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
