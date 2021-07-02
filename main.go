/*** Eland Anthony
	RxB Backend Development Project - Mockbuster REST API

Helpful Documentation
//https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
//https://golangbyexample.com/go-mod-sum-module/
//http://go-database-sql.org/overview.html
//https://chromium.googlesource.com/external/github.com/gorilla/mux/+/refs/tags/v1.2.0/README.md
//https://golang.org/pkg/database/sql/
//https://golang.org/pkg/encoding/json/

TODO:	Date: (06/28) (Begin):
Setup Docker, Golang in WSL 2 Windows env
Research into Golang syntax, structure, mux

TODO:	Date: (06/29):
Create:
- DB connection Func
- Test Connection
- Outline initial bullet point completion plans

TODO:	Date (06/30):
For:
	A list of films
	- The user should be able to search the films by title (Do last)
	- The user should be able to filter the films by rating
	- The user should be able to filter the films by category
Create:
- Film Struct
- Query Func to grab neccessary records
	- Array of type struct pointers to get and return query records
- getCategory and getRating Funcs
	- Run connect()
	- Map of route variables for rating and catagory to allow filtering
	- Run query based on route var

TODO:	Date (07/01) (To be completed 07/02?):
For:
	Satisfy the remaining two requirements involving customer comments:
		- The ability to add a customer comment to a film
		- The ability to retrieve customer comments for a film
Notes:
	- A customer table exists in schema, but a comment table does not
	- When querying for details of a single film, the query should (probably) pull in most or all columns in the film table
Create:
	- A comment table related to the customer table. Dont add to preexisting tables to reduce complexity.
	- A POST and GET function handler to post comment for a film and retrieve comments for a film
	- Append usage information to WelcomeHandler func message
Modify:
	- Append more/all of the remaining film columns to the getFilmInfo query
	- Print the json info in a more 'pretty' way
	Testing, readability, and reliability:
		- Add test cases to each func
		- Be sure comments are uniform and present
		- Double check queries, stress testing

***/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	Category    string
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", WelcomeHandler).Methods("GET")
	//r.HandleFunc("/films", testConnection).Methods("GET")
	r.HandleFunc("/films", getFilms).Methods("GET")
	r.HandleFunc("/films/ratings/{rating}", getRating).Methods("GET")
	r.HandleFunc("/films/categories/{category}", getCategory).Methods("GET")
	r.HandleFunc("/films/titles/{title}", getFilmInfo).Methods("GET")

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
		WHERE film_id = '1000'
		LIMIT 10`)

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
		rows.Scan(&flm.ID, &flm.Title, &flm.Rating, &flm.Description)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.Marshal(films)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func getRating(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//get each route variable
	vars := mux.Vars(r)

	//using the [string]string map, get string where route var = {rating} to use in query
	ratings := vars["rating"]

	//SQL query for each row of the DB
	rows, err := db.Query(
		`SELECT
			film_id,
			title,
			rating,
			description
		FROM film
		WHERE 
		rating=$1`, ratings)

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
		rows.Scan(&flm.ID, &flm.Title, &flm.Rating, &flm.Description)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.Marshal(films)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func getCategory(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//get each route variable
	vars := mux.Vars(r)

	//using the [string]string map, get string where route var = {category} to use in query
	categories := vars["category"]

	//SQL query for each row of the DB
	rows, err := db.Query(
		`SELECT DISTINCT
			f.film_id,
			f.title,
			f.rating,
			f.description,
			c.name
		FROM film f, category c, film_category fc
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			c.name=$1`, categories)

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
		rows.Scan(&flm.ID, &flm.Title, &flm.Rating, &flm.Description, &flm.Category)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.Marshal(films)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func getFilmInfo(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//get each route variable
	vars := mux.Vars(r)

	//using the [string]string map, get string where route var = {title} to use in query
	title := vars["title"]

	//replace underscores with spaces from title to match formating of 'title' in the DB
	title = strings.Replace(title, "_", " ", -1)

	//SQL query for each row of the DB
	rows, err := db.Query(
		`SELECT
			film_id,
			title,
			rating,
			description
		FROM film
		WHERE title=$1`, title)

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
		rows.Scan(&flm.ID, &flm.Title, &flm.Rating, &flm.Description)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.Marshal(films)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

/*
func testConnection(w http.ResponseWriter, r *http.Request) {
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
	res, _ := json.Marshal(WelcomeResponse{Message: "Welcome to Mockbuster!"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

type WelcomeResponse struct {
	Message string `json:"message"`
}
