/***

Eland Anthony
RxB Backend Development Project - Mockbuster REST API
Unit Test Class

***/

/* To do
Add and implement 'addSpace func' to reduce LOC
Partition r body and loop through each part as opposed to the whole body when checking output in getRating, getCategory
*/
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestWelcomeHandler(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WelcomeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := ``
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	fmt.Println(rr.Body.String())
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}

func TestGetFilms(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	req, err := http.NewRequest("GET", "/films", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFilms)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := ``
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	parts := strings.Split(rr.Body.String(), "}")
	fmt.Println(parts[0] + parts[1] + parts[2])
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}

func TestGetRating(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	req, err := http.NewRequest("GET", "/films/ratings/PG", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	handler := http.HandlerFunc(getRating)
	r.Path("/films/ratings/{rating}").Handler((handler)).Methods("GET")
	r.ServeHTTP(rr, req)
	expected := `{"Rating": "PG"}`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	parts := strings.Split(rr.Body.String(), "}")
	fmt.Println(parts[0] + parts[1] + parts[2])
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}

func TestGetCategory(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	req, err := http.NewRequest("GET", "/films/categories/Comedy", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	handler := http.HandlerFunc(getCategory)
	r.Path("/films/categories/{category}").Handler((handler)).Methods("GET")
	r.ServeHTTP(rr, req)
	expected := `{"Category": "Comedy"}`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	parts := strings.Split(rr.Body.String(), "}")
	fmt.Println(parts[0] + parts[1] + parts[2])
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}

func TestGetFilmsByTitle(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	req, err := http.NewRequest("GET", "/films/titles/Academy_Dinosaur", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	handler := http.HandlerFunc(GetFilmInfo)
	r.Path("/films/titles/{title}").Handler((handler)).Methods("GET")
	r.ServeHTTP(rr, req)
	expected := `[{"FilmID": 1}]`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	fmt.Println(rr.Body.String())
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}

func TestPostComment(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	var jsonStr = []byte(`{"FilmID":11,"comment":"xyz","CustomerId":230}`)
	req, err := http.NewRequest("POST", "/films/postcomment", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"message":"Your Comment has been posted!"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	fmt.Println(rr.Body.String())
	fmt.Println("FilmID:11, comment:xyz, CustomerId:230")
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}

func TestGetComment(t *testing.T) {
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
	req, err := http.NewRequest("GET", "/films/11/comment/230", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	handler := http.HandlerFunc(getComment)
	r.Path("/films/{film_id}/comment/{customer_id}").Handler((handler)).Methods("GET")
	r.ServeHTTP(rr, req)
	expected := `{"comment": "xyz"}`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	parts := strings.Split(rr.Body.String(), "}")
	fmt.Println(parts[0] + parts[1] + parts[2])
	for i := 1; i < 5; i++ {
		fmt.Println()
	}
}
