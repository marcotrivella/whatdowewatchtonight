package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"

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

	//loadFilms()
}

func loadFilms() {
	mappaCategorie = make(map[string][]Movie)
	films := []Movie{}
	log.Println("Query su film")
	err := db.Select(&films, "SELECT * FROM films.film ORDER BY film_date ASC")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Inizio a inserirli nella mappa in memoria")
	for _, movie_index := range films {
		movie_categorie := strings.Split(movie_index.Categorie, ",")
		for _, categoria := range movie_categorie {
			if categoria != "Foreign" {
				mappaCategorie[categoria] = append(mappaCategorie[categoria], movie_index)
			}
		}
	}
	log.Println("Finito di completare la mappa dei film in memoria")

	log.Println("Test film casuali")
	categorie := make([]string, 0)
	/*categorie = append(categorie, "Horror")*/
	//categorie = append(categorie, "Avventura")
	categorie = append(categorie, "Western")
	movies := getRandomFilms(categorie, "2016-01-01")
	for _, v := range movies {
		fmt.Println(v.Title + " - " + v.Date)
	}
}

func getCredentialsFromFile() Credentials {
	return readCredentials()
}

func readCredentials() Credentials {
	raw, err := ioutil.ReadFile(CREDENTIALS_PATH)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Credentials
	json.Unmarshal(raw, &c)
	return c
}

func getRandomFilms(categorie []string, date string) []Movie {
	movies := []Movie{}
	changes := true
	mappaCategorie_Indici := make(map[string][]int)
	if len(categorie) < 1 { //Non ci sono filtri selezionati --> li metto tutti
		categorie = getCategorie()
	}
	if date != "" {
		mappaCategorie_Indici = getAvailableIndex(categorie, date) //Prend gli indici della mappa dei film in memoria filtrata per data > a data inserita
	}

	for i := 0; i < MAX_FILM_PAGE; i++ {
		categoria_random := categorie[random(0, len(categorie)-1)]
		var index int
		if date != "" {
			if len(mappaCategorie_Indici[categoria_random]) < 1 { //Non ci sono film di una determinata categoria usciti dopo l'anno filtrato
				continue
			}
			index_temp := random(0, len(mappaCategorie_Indici[categoria_random])-1)

			if mappaCategorie_Indici[categoria_random][index_temp] == -1 { //Controllo se ho gia usato il film all'indice index
				for j := 0; j < len(mappaCategorie_Indici[categoria_random]); j++ {
					if mappaCategorie_Indici[categoria_random][j] != -1 {
						changes = true
						index = mappaCategorie_Indici[categoria_random][j]
						mappaCategorie_Indici[categoria_random][j] = -1 //Setto -1 nella mappa degli indici filtrati per sapere quali ho gia utilizzato
						break
					}
					changes = false
				}
				if !changes {
					return movies //Ho utilizzato tutti i film filtrati da una determinata data --> ritorno il vettore di film
				}
			} else {
				index = mappaCategorie_Indici[categoria_random][index_temp]
				mappaCategorie_Indici[categoria_random][index_temp] = -1 //Setto -1 nella mappa degli indici filtrati per sapere quali ho gia utilizzato
			}
		} else {
			index = random(0, len(mappaCategorie[categoria_random])-1)
		}
		movies = append(movies, mappaCategorie[categoria_random][index])
	}
	return movies
}

func getAvailableIndex(categorie []string, date string) map[string][]int {
	mappa := make(map[string][]int)
	for _, categoria := range categorie {
		array_index := make([]int, 0)
		for j := 0; j < len(mappaCategorie[categoria]); j++ {
			if strings.Compare(mappaCategorie[categoria][j].Date, date) >= 0 {
				array_index = append(array_index, j)
			}
		}
		mappa[categoria] = array_index
	}
	return mappa
}

func getCategorie() []string {
	categorie := make([]string, 0)
	for k, _ := range mappaCategorie {
		categorie = append(categorie, k)
	}

	return categorie
}

func random(min int, max int) int {
	if min == 0 && max == 0 {
		return 0
	}
	s2 := rand.NewSource(contatoreSeed)
	r2 := rand.New(s2)
	contatoreSeed++
	return r2.Intn(max-min) + min
}
