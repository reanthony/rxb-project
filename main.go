/*** Helpful Documentation
//https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
//https://golangbyexample.com/go-mod-sum-module/
//http://go-database-sql.org/overview.html
//https://chromium.googlesource.com/external/github.com/gorilla/mux/+/refs/tags/v1.2.0/README.md
//https://golang.org/pkg/database/sql/

Begin 06/28:
Setup Docker, Golang in WSL 2 Windows env
Research into Golang syntax, structure, mux

TODO (06/29):
Create:
- DB connection Func
- Test Connection
- Outline initial bullet point completion plans

TODO (06/30):
For

	A list of films
  - The user should be able to search the films by title (Do last)
  - The user should be able to filter the films by rating
  - The user should be able to filter the films by category

Create
- Film Struct
- Query Func to grab neccessary records
	- Array of type struct pointer to get and return query records
- getCategory and getRating Funcs
	- Run connect() and unnamed 'Query Func'
	- Map of route variables for rating and catagory to allow filtering
***/

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"database/sql"
	"fmt"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5555
	user     = "postgres"
	password = "postgres"
	dbname   = "dvdrental"
)

type Film struct {
	ID          int
	Title       string
	Rating      string
	Description string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", WelcomeHandler).Methods("GET")
	r.HandleFunc("/films", getFilms).Methods("GET")
	http.ListenAndServe(":8080", r)
}

//connect to postgres db for querying
func connect() *sql.DB {
	//format login details
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//open connection to postgres db using provided info
	db, err := sql.Open("postgres", psqlInfo)

	//if unsuccessful login, return panic and exit process
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	//check for connection
	if err != nil {
		fmt.Println("Error")
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

//get and return films from db... usage: curl http://localhost:8080/films
func getFilms(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//SQL query for each row of the DB
	rows, err := db.Query(
		`SELECT
			film_id,
			title,
			rating,
			description
		FROM film
		LIMIT 1`)
	if err != nil {
		log.Fatal(err)
	}

	//close row fetching once function has returned (defer)
	defer rows.Close()

	//create array of pointers of type Film
	var films []*Film

	//while rows.Next() {
	for rows.Next() {

		//allocate memory for db read
		flm := new(Film)

		//convert data read from DB to proper Go types
		rows.Scan(&flm.ID, &flm.Description, &flm.Rating, &flm.Title)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.Marshal(films)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

/*func testConnection(w http.ResponseWriter, r *http.Request) {
	db := connect()
	result, err := db.Query(`SELECT name FROM Film limit 10`)
	if err != nil {
		log.Fatal(err)
	}
	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}
*/

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(WelcomeResponse{Message: "Welcome to Mockbuster!!!!"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

type WelcomeResponse struct {
	Message string `json:"message"`
}
