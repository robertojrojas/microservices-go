package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robertojrojas/microservices-go/pets/cats-service/server"
)

func main() {

	flag.Parse()

	errChan := make(chan error, 1)

	go func() {
		errChan <- server.StartServer()
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
