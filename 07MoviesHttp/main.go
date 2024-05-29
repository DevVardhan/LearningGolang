package main

import (
	"encoding/json"

	"errors"

	"fmt"

	"log"

	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

type Director struct {
	Name string `json:"name"`
	Age  int8   `json:"age"`
}

type Movie struct {
	Name     string    `json:"name"`
	Id       string    `json:"id"`
	Rating   float32   `json:"rating" `
	Director *Director `json:"director"`
}

var movies = []Movie{
	{Name: "Movie1", Id: generateUniqueID(), Rating: 8.4, Director: &Director{Name: "John", Age: 84}},
	{Name: "Movie2", Id: generateUniqueID(), Rating: 9.0, Director: &Director{Name: "John", Age: 84}},
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func addMovies(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie
	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		log.Fatal(errors.New("cant add movies"))
	}
	if newMovie.Id == "" {
		newMovie.Id = generateUniqueID()
	}
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(newMovie)
}

func generateUniqueID() string {
	return uuid.New().String()
}

func getMovieByName(name string) (int, *Movie, error) {
	for index, movie := range movies {
		if movie.Name == name {
			return index, &movie, nil // Return a pointer to the matching movie
		}
	}
	return 0, nil, errors.New("movie not found") // Return a more specific error
}

func MovieByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)  // use of gorilla mux since http , pathvalue not working correctly - > r.pathvalue("name")
	name := params["name"] // Extract name from path segment

	if name == "" {
		// Handle missing name parameter (e.g., return 404 Not Found)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Movie not found in request")
		return
	}

	_, movie, err := getMovieByName(name)
	if err != nil {
		// Handle error (e.g., log error, return 404 Not Found)
		if err.Error() == "movie not found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Movie not found in database")
		} else {
			// Handle other potential errors (optional)
			log.Println("Error getting movie by name:", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie) // Encode the retrieved movie
}

func Run(router *mux.Router) {
	fmt.Println("Listening at localhost:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(errors.New("cant start the server"))
	}
}

func main() {
	router := mux.NewRouter()
	defer Run(router)
	router.HandleFunc("/list", getMovies).Methods("GET")
	router.HandleFunc("/add", addMovies).Methods("POST")
	router.HandleFunc("/getby/{name}", MovieByName).Methods("GET")

	router.HandleFunc("/delete/{name}", func(w http.ResponseWriter, r *http.Request) {
		parms := mux.Vars(r)
		name := parms["name"]
		index, _, err := getMovieByName(name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Name not found in database")
		}

		movies := append(movies[:index], movies[index+1:]...)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(movies)

	}).Methods("DELETE")

}
