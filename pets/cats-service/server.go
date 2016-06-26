package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var serverHostPort string

func init() {

	flag.StringVar(&dbURL, "dbURL", "root:my-secret-pw@tcp(:13306)/cats_db", "DB URL including database")
	flag.StringVar(&serverHostPort, "serverHostPort", ":8091", "Host and port server listens on")
}

// StartServer configures and starts API Server
func StartServer() {

	routes := mux.NewRouter()
	routes.HandleFunc("/api/cats", readAllHandler).Methods("GET")
	routes.HandleFunc("/api/cats", createHandler).Methods("POST")
	routes.HandleFunc("/api/cats/{id:[0-9]+}", readHandler).Methods("GET")
	routes.HandleFunc("/api/cats/{id:[0-9]+}", updateHandler).Methods("PUT")
	routes.HandleFunc("/api/cats/{id:[0-9]+}", deleteHandler).Methods("DELETE")
	http.Handle("/", routes)

	log.Printf("Connecting to DB[%s]....\n", dbURL)
	ConnectToDB()

	log.Printf("Listening on [%s]....\n", serverHostPort)
	http.ListenAndServe(serverHostPort, nil)
}
