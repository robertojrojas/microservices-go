package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type RouterHandler interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	Methods(methods ...string) *mux.Route
}

// CatsRoutesHandler represents the HTTP Handler methods
type CatsRoutesHandler interface {
	CreateHandler(w http.ResponseWriter, r *http.Request) error
	ReadAllHandler(w http.ResponseWriter, r *http.Request) error
	ReadHandler(w http.ResponseWriter, r *http.Request) error
	UpdateHandler(w http.ResponseWriter, r *http.Request) error
	DeleteHandler(w http.ResponseWriter, r *http.Request) error
}

// CatsRoutes implements HTTP Handlers
type CatsRoutes struct {
	catsDBStore CatsDataStore
	CatsRoutesHandler
}

// SetupCatsRoutes return a new CatsRoutesHandler
func SetupCatsRoutes(dataStore CatsDataStore, router RouterHandler) CatsRoutesHandler {
	cr := &CatsRoutes{
		catsDBStore: dataStore,
	}

	router.HandleFunc("/api/cats", errorHandler(cr.ReadAllHandler)).Methods("GET")
	router.HandleFunc("/api/cats", errorHandler(cr.CreateHandler)).Methods("POST")
	router.HandleFunc("/api/cats/{id:[0-9]+}", errorHandler(cr.ReadHandler)).Methods("GET")
	router.HandleFunc("/api/cats/{id:[0-9]+}", errorHandler(cr.UpdateHandler)).Methods("PUT")
	router.HandleFunc("/api/cats/{id:[0-9]+}", errorHandler(cr.DeleteHandler)).Methods("DELETE")

	return cr

}

func (cr *CatsRoutes) CreateHandler(w http.ResponseWriter, r *http.Request) error {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	cat := &Cat{}
	err = json.Unmarshal(data, cat)
	if err != nil {
		return err
	}

	err = cr.catsDBStore.CreateCat(cat)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func (cr *CatsRoutes) ReadAllHandler(w http.ResponseWriter, r *http.Request) error {

	cats, err := cr.catsDBStore.ReadAllCats()
	if err != nil {
		return err
	}
	catsData, err := json.Marshal(cats)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(catsData))

	return nil

}

func (cr *CatsRoutes) ReadHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	cat, err := cr.catsDBStore.ReadCat(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return notFound{}
		} else {
			return err
		}
	}

	catData, err := json.Marshal(cat)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(catData))

	return nil
}

func (cr *CatsRoutes) UpdateHandler(w http.ResponseWriter, r *http.Request) error {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	cat := &Cat{}
	err = json.Unmarshal(data, cat)
	if err != nil {
		return err
	}
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	cat.ID = int64(id)

	err = cr.catsDBStore.UpdateCat(cat)
	return err

}

func (cr *CatsRoutes) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	err = cr.catsDBStore.DeleteCat(id)
	return err
}

/*
  The following code was copied from: https://github.com/campoy/todo/blob/master/server/server.go#L59
  thanks @francesc
*/

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}
