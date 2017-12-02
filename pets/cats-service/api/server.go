package api

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var serverHostPort string

func init() {
	flag.StringVar(&serverHostPort, "http", ":8091", "Host and port server listens on")
}

// StartServer configures and starts API Server
func StartServer() error {

	config := getConfig()
	router := mux.NewRouter()
	catsDB, err := NewCatsDB()
	if err != nil {
		return err
	}

	SetupCatsRoutes(catsDB, router)

	http.Handle("/", router)

	log.Printf("Listening on [%s]....\n", serverHostPort)
	return http.ListenAndServe(serverHostPort, nil)
}
