package api


import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"os"
)

var serverHostPort string
const mysqlDBURI_key = "MYSQL_DB_URI"
const defaultMYSQLURI = "root:my-secret-pw@tcp(:3306)/cats_db"

func init() {
	flag.StringVar(&serverHostPort, "http", ":8091", "Host and port server listens on")
}

type config struct {
	mySQLDBURI string
}

// StartServer configures and starts API Server
func StartServer() error {

	config := getConfig()
	router := mux.NewRouter()
	catsDB := NewCatsDB(config.mySQLDBURI)

	SetupCatsRoutes(catsDB, router)

	http.Handle("/", router)

	log.Printf("Listening on [%s]....\n", serverHostPort)
	return http.ListenAndServe(serverHostPort, nil)
}

func getConfig() *config {

	envConfig := config{}

	envConfig.mySQLDBURI = os.Getenv(mysqlDBURI_key)
	if envConfig.mySQLDBURI == "" {
		envConfig.mySQLDBURI = defaultMYSQLURI
	}

	return &envConfig
}
