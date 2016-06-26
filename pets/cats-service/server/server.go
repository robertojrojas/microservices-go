package server

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robertojrojas/microservices-go/pets/cats-service/models"
	"github.com/robertojrojas/microservices-go/pets/cats-service/routes"
)

var serverHostPort string
var dbURL *string

func init() {

	dbURL = flag.String("dbURL", "root:my-secret-pw@tcp(:3306)/cats_db", "DB URL including database")
	flag.StringVar(&serverHostPort, "serverHostPort", ":8091", "Host and port server listens on")

}

// StartServer configures and starts API Server
func StartServer() {

	router := mux.NewRouter()
	catsDB := models.NewCatsDB(*dbURL)
	catsRoutes := routes.NewCatsRoutes(catsDB)
	router.HandleFunc("/api/cats", catsRoutes.ReadAllHandler).Methods("GET")
	router.HandleFunc("/api/cats", catsRoutes.CreateHandler).Methods("POST")
	router.HandleFunc("/api/cats/{id:[0-9]+}", catsRoutes.ReadHandler).Methods("GET")
	router.HandleFunc("/api/cats/{id:[0-9]+}", catsRoutes.UpdateHandler).Methods("PUT")
	router.HandleFunc("/api/cats/{id:[0-9]+}", catsRoutes.DeleteHandler).Methods("DELETE")
	http.Handle("/", router)

	log.Printf("Listening on [%s]....\n", serverHostPort)
	http.ListenAndServe(serverHostPort, nil)
}
