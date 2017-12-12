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

	loadFilms()
}

func loadFilms() {
	/*movies := []Movie{}
	query := "SELECT * FROM films.film"
	if !date.IsZero() {
		query += "WHERE film_date = " + date.String()
	}
	return movies*/
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

	log.Println("Test film casuali di animazione, Avventura e western")
	categorie := make([]string, 0)
	/*categorie = append(categorie, "Horror")
	categorie = append(categorie, "Avventura")*/
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
	raw, err := ioutil.ReadFile(CREDENTIALS_PATH + "/credentials.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Credentials
	json.Unmarshal(raw, &c)
	return c
}

func getRandomFilms(categorie []string, date string) []Movie {
	index_start := 0
	movies := []Movie{}
	mappaCategorie_Indici := make(map[string][]int)
	if len(categorie) < 1 { //Non ci sono filtri selezionati --> li metto tutti
		for k, _ := range mappaCategorie {
			categorie = append(categorie, k)
		}
	}
	if date != "" {
		mappaCategorie_Indici = getAvailableIndex(categorie, date)
	}
	/*for k, _ := range mappaCategorie_Indici {
		for _, v := range mappaCategorie_Indici[k] {
			fmt.Println(k + " - " + strconv.Itoa(v))
		}
	}*/
	for i := 0; i < MAX_FILM_PAGE; i++ {
		categoria_random := categorie[random(0, len(categorie)-1)]
		index_max := len(mappaCategorie[categoria_random]) - 1
		var index int
		if date != "" {
			if len(mappaCategorie_Indici[categoria_random]) < 1 { //Non ci sono film di una determinata categoria usciti dopo l'anno filtrato
				continue
			}
			index_start = mappaCategorie_Indici[categoria_random][0]
			index = random(0, index_max-index_start)
			/*for _, v := range mappaCategorie_Indici[categoria_random] {
				fmt.Println(v)
			}
			fmt.Println()*/
			if mappaCategorie_Indici[categoria_random][index] == 0 {
				for j := 0; j < len(mappaCategorie_Indici[categoria_random]); j++ {
					if mappaCategorie_Indici[categoria_random][j] != 0 {
						index = mappaCategorie_Indici[categoria_random][j]
						break
					}
				}
			} else {
				mappaCategorie_Indici[categoria_random][index] = 0
			}
		} else {
			index = random(index_start, index_max)
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

func random(min int, max int) int {
	if min == 0 && max == 0 {
		return 0
	}
	s2 := rand.NewSource(contatoreSeed)
	r2 := rand.New(s2)
	contatoreSeed++
	return r2.Intn(max-min) + min
}
