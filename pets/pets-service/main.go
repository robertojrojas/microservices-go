package main

import (
	"log"

	"flag"
	"fmt"
	"github.com/robertojrojas/microservices-go/pets/pets-service/server"
	"os"
	"os/signal"
	"syscall"
)

var serverHostPort string

func init() {
	flag.StringVar(&serverHostPort, "http", ":8094", "Host and port server listens on")
}

func main() {
	flag.Parse()

	errChan := make(chan error, 1)

	go func() {
		errChan <- server.StartServer(serverHostPort)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exciting...", s))
			os.Exit(0)
		}
	}

}
