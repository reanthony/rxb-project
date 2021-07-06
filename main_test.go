package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFilms(t *testing.T) {
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

	// Check the response body is what we expect.
	expected := `[{"FilmID": 282,"Title": "Encounters Curtain","Rating": "NC-17","Description": "A Insightful Epistle of a Pastry Chef And a Womanizer who must Build a Boat in New Orleans","Category": "Drama","ReleaseYear": 2006,"Language": "English             ","ActorFName": "Alec","ActorLName": "Wayne","Length": 92,"RentalDuration": 5,"RentalRate": 0.99,"ReplacementCost": 20.99,"SpeacialFeatures": "{Trailers}"}]`
	//!strings.Contains(str, "-numinput")
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
