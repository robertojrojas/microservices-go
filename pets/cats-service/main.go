package main

import (
	"flag"

	"github.com/robertojrojas/microservices-go/pets/cats-service/server"
)

func main() {

	flag.Parse()
	server.StartServer()
}
