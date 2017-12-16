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
	templatePath    = "static/html"
	credentialsPath = "static/data/credentials.json"
	maxFilmPage     = 50
)

//Context da passare al template html
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

	initDB()

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func templateHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		serveTemplate(w, req, "index", Context{Movies: getRandomFilmsFiltered(make([]string, 0), "")})
	} else if req.URL.Path == "/films" {
		getRandomFilms(w, req)
	} else {
		return
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request, fileName string, context Context) {
	lp := filepath.Join(templatePath, "layout.html")
	fp := filepath.Join(templatePath, fileName+".html")

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
