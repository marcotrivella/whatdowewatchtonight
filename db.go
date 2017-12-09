package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Movie struct {
	ID           int    `db:"film_id"`
	Asin         string `db:"film_asin"`
	AmazonURL    string `db:"film_pageurl"`
	Actors       string `db:"film_actors"`
	Publisher    string `db:"film_publisher"`
	Title        string `db:"film_title"`
	Categorie    string `db:"film_categorie"`
	Manufacturer string `db:"film_manufacturer"`
	Locandina    string `db:"film_locandina"`
	Description  string `db:"film_description"`
	IDTMD        int    `db:"film_idtmdb"`
	Date         string `db:"film_date"`
	Director     string `db:"film_director"`
	/*
		Animazione
		Avventura
		Azione
		Dramma
		Fantasy
		Commedia
		Famiglia
		Romance
		Fantascienza
		Crime
		Foreign
		Thriller
		Musica
		Mistero
		Western
		Storia
		Guerra
		Horror
	*/
}

type Credentials struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	DBName   string `json:"DBName"`
}

func INIT_DB() {
	credentials := getCredentialsFromFile()
	var err error

	db, err = sqlx.Connect("mysql", credentials.User+":"+credentials.Password+"@tcp(127.0.0.1:3306)/"+credentials.DBName)
	if err != nil {
		log.Fatalln(err)
	}

	loadFilms()
}

func loadFilms() {
	mappaCategorie = make(map[string][]Movie)
	films := []Movie{}
	log.Println("Query su film")
	err := db.Select(&films, "SELECT * FROM films.film")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Inizio a inserirli nella mappa in memoria")
	for _, movie_index := range films {
		movie_categorie := strings.Split(movie_index.Categorie, ",")
		for _, categoria := range movie_categorie {
			mappaCategorie[categoria] = append(mappaCategorie[categoria], movie_index)
		}
	}
	log.Println("Finito di completare la mappa dei film in memoria")
}

func getCredentialsFromFile() Credentials {
	return readCredentials()
}

func readCredentials() Credentials {
	raw, err := ioutil.ReadFile(CREDENTIALS_PATH + "/credentials.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Credentials
	json.Unmarshal(raw, &c)
	return c
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
