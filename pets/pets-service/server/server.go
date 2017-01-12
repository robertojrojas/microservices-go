package server

import (
	"encoding/json"
	"log"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/robertojrojas/microservices-go/pets/pets-service/executor"
	"github.com/robertojrojas/microservices-go/pets/pets-service/services"
	"github.com/robertojrojas/microservices-go/pets/pets-service/cache"
	"net/http"
	"os"
	"io"
	"time"
)

const cacheKey = "PETS-SVC:ALLPETS"
const cacheTTL = 30 * time.Second


type PetsServer struct {
	rpcExecutor *executor.RPCExecutor
        cache 	*cache.RedisCache
}

type config struct {
	catsServiceURI  string
	birdsServiceURI string
	rabbitMQURI     string
	rabbitMQQueue   string
	redisURI        string
}

// StartServer configures and starts API Server
func StartServer(serverHostPort string) error {
	appConfig := getConfig()

	redisCache, err := cache.NewRedisCache(appConfig.redisURI)
	if err != nil {
		return err
	}

	catService := &services.CatService{
		URL: appConfig.catsServiceURI,
	}
	birdsService := &services.BirdService{
		ServiceAddress: appConfig.birdsServiceURI,
	}

	dogsService := &services.DogsService{
		ServiceAddress: appConfig.rabbitMQURI,
		RPCQueue:       appConfig.rabbitMQQueue,
	}
	rpcExecutor := executor.NewRPCExecutor(catService, birdsService, dogsService)

	petsServer := &PetsServer{
		rpcExecutor: rpcExecutor,
		cache: redisCache,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/pets", petsServer.petsHandler).Methods("GET")

	http.Handle("/", router)

	return http.ListenAndServe(serverHostPort, nil)
}

func (s *PetsServer) petsHandler(w http.ResponseWriter, r *http.Request) {
	err := getAllPets(s, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAllPets(s *PetsServer, w io.Writer) (error) {

	// For now we ignore errors from cache since
	// the read source of the data is the rpcExecutor
	allPetsJSON, _ := s.cache.Get(cacheKey)

	if allPetsJSON == "" {
		log.Println("No Pets found in cache. Querying sources....")
		results, err := s.rpcExecutor.GetAllPets()
		if err != nil {
			return err
		}

		data, err := json.Marshal(&results)
		if err != nil {
			return err
		}
		allPetsJSON = string(data)
		err = s.cache.Store(cacheKey, allPetsJSON, cacheTTL)
		if err != nil {
			return err
		}
	} else {
		log.Println("Pets found in cache")
	}

	fmt.Fprintf(w, "%s\n", allPetsJSON)

	return nil
}

func getConfig() *config {

	envConfig := config{}

	envConfig.rabbitMQURI = os.Getenv("RABBITMQ_URI")
	if envConfig.rabbitMQURI == "" {
		envConfig.rabbitMQURI = "amqp://guest:guest@localhost:5672/"
	}

	envConfig.rabbitMQQueue = os.Getenv("RABBITMQ_QUEUE")
	if envConfig.rabbitMQQueue == "" {
		envConfig.rabbitMQQueue = "dog_service_rpc_queue"
	}

	envConfig.redisURI = os.Getenv("REDIS_URI")
	if envConfig.redisURI == "" {
		envConfig.redisURI = ":6379"
	}

	envConfig.catsServiceURI = os.Getenv("CATS_SERVICE_URI")
	if envConfig.catsServiceURI == "" {
		envConfig.catsServiceURI = "http://cats-service-svc:8091/api/cats"
	}

	envConfig.birdsServiceURI = os.Getenv("BIRDS_SERVICE_URI")
	if envConfig.birdsServiceURI == "" {
		envConfig.birdsServiceURI = "birds-service-svc:8092"
	}

	return &envConfig
}
