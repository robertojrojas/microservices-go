package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/robertojrojas/microservices-go/pets/cats-service/models"
)

// CatsRoutesHandler represents the HTTP Handler methods
type CatsRoutesHandler interface {
	CreateHandler(w http.ResponseWriter, r *http.Request)
	ReadAllHandler(w http.ResponseWriter, r *http.Request)
	ReadHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
}

// CatsRoutes implements HTTP Handlers
type CatsRoutes struct {
	catsDBStore models.CatsDataStore
	CatsRoutesHandler
}

// NewCatsRoutes return a new CatsRoutesHandler
func NewCatsRoutes(dataStore models.CatsDataStore) CatsRoutesHandler {
	return &CatsRoutes{
		catsDBStore: dataStore,
	}
}

func (cr *CatsRoutes) CreateHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cat := &models.Cat{}
	err = json.Unmarshal(data, cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = cr.catsDBStore.CreateCat(cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cr *CatsRoutes) ReadAllHandler(w http.ResponseWriter, r *http.Request) {

	cats, err := cr.catsDBStore.ReadAllCats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	catsData, err := json.Marshal(cats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(catsData))

}

func (cr *CatsRoutes) ReadHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cat, err := cr.catsDBStore.ReadCat(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "cat not found for id: "+idStr)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	catData, err := json.Marshal(cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(catData))
}

func (cr *CatsRoutes) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cat := &models.Cat{}
	err = json.Unmarshal(data, cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	cat.ID = int64(id)

	err = cr.catsDBStore.UpdateCat(cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (cr *CatsRoutes) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	err = cr.catsDBStore.DeleteCat(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
