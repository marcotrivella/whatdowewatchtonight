package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//Movie struttura per film
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

//Credentials struttura per le credenziali
type Credentials struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	DBName   string `json:"DBName"`
}

func initDB() {
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
	err := db.Select(&films, "SELECT * FROM films.film ORDER BY film_date ASC")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Inizio a inserirli nella mappa in memoria")
	for _, movieIndex := range films {
		movieCategorie := strings.Split(movieIndex.Categorie, ",")
		for _, categoria := range movieCategorie {
			if categoria != "Foreign" {
				mappaCategorie[categoria] = append(mappaCategorie[categoria], movieIndex)
			}
		}
	}
	log.Println("Finito di completare la mappa dei film in memoria")
}

func getCredentialsFromFile() Credentials {
	return readCredentials()
}

func readCredentials() Credentials {
	raw, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Credentials
	json.Unmarshal(raw, &c)
	return c
}

func getRandomFilms(w http.ResponseWriter, req *http.Request) {
	categorie := strings.Split(req.FormValue("categorie"), ",")
	date := req.FormValue("date")
	movies := getRandomFilmsFiltered(categorie, date)
	moviesJSON, _ := json.Marshal(movies)
	fmt.Fprintf(w, string(moviesJSON))
	return
}

func getRandomFilmsFiltered(categorie []string, date string) []Movie {
	movies := []Movie{}
	changes := true
	mappaCategorieIndici := make(map[string][]int)
	if len(categorie) < 1 { //Non ci sono filtri selezionati --> li metto tutti
		categorie = getCategorie()
	}
	if date != "" {
		mappaCategorieIndici = getAvailableIndex(categorie, date) //Prend gli indici della mappa dei film in memoria filtrata per data > a data inserita
	}

	for i := 0; i < maxFilmPage; i++ {
		categoriaRandom := categorie[random(0, len(categorie)-1)]
		var index int
		if date != "" {
			if len(mappaCategorieIndici[categoriaRandom]) < 1 { //Non ci sono film di una determinata categoria usciti dopo l'anno filtrato
				continue
			}
			indexTemp := random(0, len(mappaCategorieIndici[categoriaRandom])-1)

			if mappaCategorieIndici[categoriaRandom][indexTemp] == -1 { //Controllo se ho gia usato il film all'indice index
				for j := 0; j < len(mappaCategorieIndici[categoriaRandom]); j++ {
					if mappaCategorieIndici[categoriaRandom][j] != -1 {
						changes = true
						index = mappaCategorieIndici[categoriaRandom][j]
						mappaCategorieIndici[categoriaRandom][j] = -1 //Setto -1 nella mappa degli indici filtrati per sapere quali ho gia utilizzato
						break
					}
					changes = false
				}
				if !changes {
					return movies //Ho utilizzato tutti i film filtrati da una determinata data --> ritorno il vettore di film
				}
			} else {
				index = mappaCategorieIndici[categoriaRandom][indexTemp]
				mappaCategorieIndici[categoriaRandom][indexTemp] = -1 //Setto -1 nella mappa degli indici filtrati per sapere quali ho gia utilizzato
			}
		} else {
			index = random(0, len(mappaCategorie[categoriaRandom])-1)
		}
		movies = append(movies, mappaCategorie[categoriaRandom][index])
	}
	return movies
}

func getAvailableIndex(categorie []string, date string) map[string][]int {
	mappa := make(map[string][]int)
	for _, categoria := range categorie {
		arrayIndex := make([]int, 0)
		for j := 0; j < len(mappaCategorie[categoria]); j++ {
			if strings.Compare(mappaCategorie[categoria][j].Date, date) >= 0 {
				arrayIndex = append(arrayIndex, j)
			}
		}
		mappa[categoria] = arrayIndex
	}
	return mappa
}

func getCategorie() []string {
	categorie := make([]string, 0)
	for k := range mappaCategorie {
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
