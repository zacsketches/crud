// CRUD is an application that provides a golang API for a movie
// database.  This app is based on the hackernoon tutorial available
// here, https://hackernoon.com/build-restful-api-in-go-and-mongodb-5e7f2ec4be94
// However, I've replaced the use of mgo go driver for the officials
// mongo-go-driver.
//
// API spec
// 		Get     /movies			Get list of movies
//		Get     /movies/:id		Find a movie by its ID
//		Post    /movies			Create a new movie
//		Put	    /movies			Update an existing movie
//		Delete  /movies			Delete an existing movie
//
// Some other links for where to go next include:
// https://www.alexedwards.net/blog/a-recap-of-request-handling
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
// https://gowebexamples.com/forms/
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial-part-1-connecting-using-bson-and-crud-operations
//
// Over the course of creating this example I found that Stack Overflow lacks
// specific examples for the use of Atlas and the offical mongo-go-driver so I wrote
// two worked examples
// https://stackoverflow.com/questions/55564562/what-is-the-bson-syntax-for-set-in-updateone-for-the-official-mongo-go-driver
// https://stackoverflow.com/questions/55554772/how-do-you-connect-to-mongodb-atlas-with-the-official-mongo-go-driver

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"github.com/zacsketches/crud/dao"
	"github.com/zacsketches/crud/models"
)

var db = dao.MoviesDAO{
	User:     "zac-admin",
	Server:   "zacs-garden-47y0p.mongodb.net",
	Database: "test",
}

func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
	movies, err := db.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movies)
}

func FindMovieEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := db.FindByID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	res, err := db.Insert(movie)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	movie.ID = res.InsertedID.(primitive.ObjectID)
	respondWithJson(w, http.StatusCreated, movie)
}

func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Println(movie)
	res, err := db.Update(movie)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("%+v\n", res)
	if res.MatchedCount < 1 {
		respondWithJson(w, http.StatusOK, map[string]string{"status": "no op"})
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	res, err := db.Delete(movie)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if res.DeletedCount < 1 {
		respondWithJson(w, http.StatusOK, map[string]string{"result": "no op"})
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", FindMovieEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
