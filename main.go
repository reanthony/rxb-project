/***

Eland Anthony
RxB Backend Development Project - Mockbuster REST API

***/

/*
Useful queries for testing:

const commentQueryInsert = `INSERT INTO comment(
comment,
customer_id,
film_id)
VALUES('ThisIsATest','20','100')`

const commentQueryDropTable = `DROP TABLE comment`
*/

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

const CreateCommentTableIfNotExists = `CREATE TABLE IF NOT EXISTS comment (
		comment_id SERIAL PRIMARY KEY,
		comment VARCHAR NOT NULL,
		customer_id INT NOT NULL,
		film_id INT NOT NULL)`

type Film struct {
	FilmID           int
	Title            string  `json:"Title,omitempty"`
	Rating           string  `json:"Rating,omitempty"`
	Description      string  `json:"Description,omitempty"`
	Category         string  `json:"Category,omitempty"`
	ReleaseYear      int     `json:"ReleaseYear,omitempty"`
	Language         string  `json:"Language,omitempty"`
	ActorFName       string  `json:"ActorFName,omitempty"`
	ActorLName       string  `json:"ActorLName,omitempty"`
	Length           int     `json:"Length,omitempty"`
	RentalDuration   int     `json:"RentalDuration,omitempty"`
	RentalRate       float32 `json:"RentalRate,omitempty"`
	ReplacementCost  float32 `json:"ReplacementCost,omitempty"`
	SpeacialFeatures string  `json:"SpeacialFeatures,omitempty"`
	CustomerId       int     `json:"CustomerId,omitempty"`
	CommentId        int     `json:"CommentId,omitempty"`
	Comment          string  `json:"comment,omitempty"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", WelcomeHandler).Methods("GET")
	r.HandleFunc("/films", GetFilms).Methods("GET")
	r.HandleFunc("/films/ratings/{rating}", getRating).Methods("GET")
	r.HandleFunc("/films/categories/{category}", getCategory).Methods("GET")
	r.HandleFunc("/films/titles/{title}", GetFilmInfo).Methods("GET")
	r.HandleFunc("/films/postcomment", postComment).Methods("POST")
	r.HandleFunc("/films/{film_id}/comment/{customer_id}", getComment).Methods("GET")

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
		panic(err)
	}

	fmt.Println("")
	fmt.Println("Successfully connected!")
	fmt.Println("")

	return db
}

//get and return films and details from db
func GetFilms(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//SQL query for each row of the DB
	rows, err := db.Query(
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.title,
			f.rating,
			f.description,
			f.release_year,
			f.Length,
			a.first_name,
			a.last_name,
			l.name,
			c.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, category c, film_category fc, language l, actor a, film_actor fa
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			l.language_id = f.language_id AND
			fa.film_id = f.film_id AND
			a.actor_id = fa.actor_id`)

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
		rows.Scan(&flm.FilmID, &flm.Title, &flm.Rating, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.ActorFName, &flm.ActorLName, &flm.Language, &flm.Category, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpeacialFeatures)

		//add film to array

		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.MarshalIndent(films, "", "	")
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
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.title,
			f.rating,
			f.description,
			f.release_year,
			f.Length,
			a.first_name,
			a.last_name,
			l.name,
			c.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, category c, film_category fc, language l, actor a, film_actor fa
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			l.language_id = f.language_id AND
			fa.film_id = f.film_id AND
			a.actor_id = fa.actor_id AND
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
		rows.Scan(&flm.FilmID, &flm.Title, &flm.Rating, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.ActorFName, &flm.ActorLName, &flm.Language, &flm.Category, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpeacialFeatures)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.MarshalIndent(films, "", "	")
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
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.title,
			f.rating,
			f.description,
			f.release_year,
			f.Length,
			a.first_name,
			a.last_name,
			l.name,
			c.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, category c, film_category fc, language l, actor a, film_actor fa
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			l.language_id = f.language_id AND
			fa.film_id = f.film_id AND
			a.actor_id = fa.actor_id AND
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
		rows.Scan(&flm.FilmID, &flm.Title, &flm.Rating, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.ActorFName, &flm.ActorLName, &flm.Language, &flm.Category, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpeacialFeatures)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.MarshalIndent(films, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func GetFilmInfo(w http.ResponseWriter, r *http.Request) {

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
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.title,
			f.rating,
			f.description,
			f.release_year,
			f.Length,
			a.first_name,
			a.last_name,
			l.name,
			c.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, category c, film_category fc, language l, actor a, film_actor fa
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			l.language_id = f.language_id AND
			fa.film_id = f.film_id AND
			a.actor_id = fa.actor_id AND f.title=$1`, title)

	if err != nil {
		log.Fatal(err)
	}

	//close row fetching once function has returned (defer)
	defer rows.Close()

	//create array of pointers of type Film
	var films []*Film

	for rows.Next() {

		//allocate memory for db read
		flm := new(Film)

		//convert data read from DB to proper Go types
		rows.Scan(&flm.FilmID, &flm.Title, &flm.Rating, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.ActorFName, &flm.ActorLName, &flm.Language, &flm.Category, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpeacialFeatures)

		//add film to array
		films = append(films, flm)
	}

	//marshall format to json
	res, _ := json.MarshalIndent(films, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func postComment(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//ceate comment table if not exists
	db.Exec(CreateCommentTableIfNotExists)

	//cannot use mux here with a post
	//must use a decoder to translate input from json

	//create instance of Film struct to store inputs
	var cmt Film

	//decode post args point it to cmt
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cmt)
	if err != nil {
		log.Fatal(err)
	}

	//assign vars needed to the given post args
	comment := cmt.Comment
	customer_id := cmt.CustomerId
	film_id := cmt.FilmID

	//must initialize comment_id despite its serial assignment
	comment_id := 0

	//excecute insert query to insert the given info to the comment table
	//query must return... return the new comment id which should incrememnt by 1 for each comment
	//due to its 'seriel' attribute
	err = db.QueryRow(
		`INSERT INTO comment(
			comment,
			customer_id,
			film_id)
		VALUES($1,$2,$3) RETURNING comment_id`, comment, customer_id, film_id).Scan(&comment_id)

	if err != nil {
		log.Fatal(err)
	}

	//marshall format to json
	res, _ := json.Marshal(WelcomeResponse{Message: "Your Comment has been posted!"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func getComment(w http.ResponseWriter, r *http.Request) {

	//connect to postgres db. look for 'Successfully Connected' output
	db := connect()

	//get each route variable
	vars := mux.Vars(r)

	//using the [string]string map, get string where route var = {rating} to use in query
	film_id := vars["film_id"]
	customer_id := vars["customer_id"]

	//SQL query for each row of the DB
	rows, err := db.Query(
		`SELECT DISTINCT ON (c.comment_id)
			f.film_id,
			f.title,
			c.comment,
			cu.customer_id,
			c.comment_id
		FROM film f, comment c, customer cu
		WHERE c.film_id=f.film_id AND c.customer_id=cu.customer_id 
		AND c.film_id=$1 AND c.customer_id=$2`, film_id, customer_id)

	if err != nil {
		log.Fatal(err)
	}

	//create array of pointers of type Film
	var comment []*Film
	for rows.Next() {

		//allocate memory for db read
		cmt := new(Film)

		//convert data read from DB to proper Go types
		rows.Scan(&cmt.FilmID, &cmt.Title, &cmt.Comment, &cmt.CustomerId, &cmt.CommentId)

		//add comment to array
		comment = append(comment, cmt)
	}

	//marshall format to json
	res, _ := json.MarshalIndent(comment, "", "		")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := json.MarshalIndent(WelcomeResponse{Message: "Welcome to Mockbuster! See Console for usage information."}, "\n", "	")

	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("View all films and key information: curl http://localhost:8080/films")
	fmt.Println("View all films and key information by category: curl http://localhost:8080/films/categories/[desired category]")
	fmt.Println("View all films and key information by rating: curl http://localhost:8080/films/ratings/[desired rating]")
	fmt.Println("View all information for a desired film name: curl http://localhost:8080/films/titles/[desired title[title name should deliminate words via underscores]]")
	fmt.Println("Post a comment to a fil using its Film ID and your customer ID: curl -X POST -d '{\"comment\":\"[desired comment]]\", \"ID\":[desired film id], \"CustomerId\":[your customer id]}' http://localhost:8080/films/postcomment")
	fmt.Println("Fetch all of your comments for a film using its Film ID and your Customer ID: curl http://localhost:8080/films/[desired film id]/comment/[your customer id]")
	fmt.Println("")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

type WelcomeResponse struct {
	Message string `json:"message"`
}
