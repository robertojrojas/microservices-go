package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robertojrojas/microservices-go/pets/cats-service/api"
	"errors"
)


func main() {

	flag.Parse()

	errChan := make(chan error, 1)

	go func() {
		errChan <- api.StartServer()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	err := shutdownHook(signalChan, errChan)
	log.Fatal(err)


}


func shutdownHook(signalChan chan os.Signal, errChan chan error) error {
	for {
		select {
		case err := <-errChan:
			return err
		case s := <-signalChan:
			return errors.New(fmt.Sprintf("Captured %v. Exiting...", s))
		}
	}
}
