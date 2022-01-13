package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	//a.initializeRoutes()
}

func (a *App) Run(address string) {

}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}

	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	var count, start int

	query := r.URL.Query()
	countQueryString, countPresent := query["count"]
	startQueryString, startPresent := query["start"]

	if !countPresent || len(countQueryString) == 0 {
		count = 1
	} else {
		count, _ = strconv.Atoi(countQueryString[0])
	}

	if !startPresent || len(startQueryString) == 0 {
		start = 0
	} else {
		start, _ = strconv.Atoi(countQueryString[0])
	}
	// count, _ := strconv.Atoi(r.URL.Query()["count"][0])
	// start, _ := strconv.Atoi(r.URL.Query()["start"][0])

	if count > 10 {
		count = 1
	}

	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
