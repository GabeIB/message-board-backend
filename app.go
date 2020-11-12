// app.go

package main

import (
	"fmt"
	"log"
	"database/sql"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App is a structure that holds pointers to the http request multiplexer and the database
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialize establishes a connection to the database and initializes API endpoints.
func (a *App) Initialize(dbUname, dbPass, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUname, dbPass, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//Run starts the server on a given port.
func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

//respondWithJSON responds to an http request with a JSON formatted response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

//respondWithError responds to an http request with error code and error message.
//In production code, where security is of greater concern, I would be much more careful about limiting the amount of error information given to the user.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

//authenticates checks the username and password of the basic authrization of an http request.
//In production code, I would have a database of usernames and password hashes.
func authenticate(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if ok {
		if username == "admin" && password == "back-challenge" {
			return true
		}
	}
	return false
}

//getMessage retrieves a specific message identified by id
func (a *App) getMessage(w http.ResponseWriter, r *http.Request) {
	if authenticate(r){
		vars := mux.Vars(r)
		id := vars["id"]
		m := message{ID: id}
		if err := m.getMessage(a.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
			    respondWithError(w, http.StatusNotFound, "Product not found")
			default:
			    respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		respondWithJSON(w, http.StatusOK, m)
	}else{
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
	}
}

//getMessages sends all messages to the client in JSON format.
func (a *App) getMessages(w http.ResponseWriter, r *http.Request) {
	if authenticate(r){
		messages, err := getMessages(a.DB)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, messages)
	}else{
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
	}
}

//createMessage adds a message to the database.
//createMessage does not require authentication.
func (a *App) createMessage(w http.ResponseWriter, r *http.Request) {
	var m message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if m.Name == "" || m.Email == "" || m.Text == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := m.createMessage(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, m)
}

//updateMessage looks up a message by id in database and modifies fields to match arguments
func (a *App) updateMessage(w http.ResponseWriter, r *http.Request) {
	if authenticate(r){
		vars := mux.Vars(r)
		id := vars["id"]
		m := message{ID: id}
		if err := m.updateMessage(a.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
			    respondWithError(w, http.StatusNotFound, "Product not found")
			default:
			    respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		respondWithJSON(w, http.StatusOK, m)
	}else{
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
	}
}

//initializeRoutes adds all API endpoints to the HTTP request multiplexer.
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/messages", a.getMessages).Methods("GET")
	a.Router.HandleFunc("/message", a.createMessage).Methods("POST")
	a.Router.HandleFunc("/message/{id:[0-9a-fA-f-]+}", a.getMessage).Methods("GET")
	a.Router.HandleFunc("/message/{id:[0-9a-fA-f-]+}", a.updateMessage).Methods("PUT")
}
