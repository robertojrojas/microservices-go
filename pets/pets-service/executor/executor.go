package executor

import (
	"errors"
	"log"
	"time"

	"golang.org/x/net/context"

	birdsModels "github.com/robertojrojas/microservices-go/pets/birds-service/models"
	catModels "github.com/robertojrojas/microservices-go/pets/cats-service/api"
	dogsModels "github.com/robertojrojas/microservices-go/pets/dogs-service/models"
	"github.com/robertojrojas/microservices-go/pets/pets-service/services"
)

// RPCExecutor is in charge of performing all the RPC requests.
type RPCExecutor struct {
	CatService *services.CatService

	BirdsService *services.BirdService

	DogsService *services.DogsService
}

type RPCResult struct {
	Cats  []*catModels.Cat    `json:"cats"`
	Birds []*birdsModels.Bird `json:"birds"`
	Dogs  []*dogsModels.Dog   `json:"dogs"`
}

func NewRPCExecutor(catService *services.CatService, birdsService *services.BirdService,
	dogsService *services.DogsService) *RPCExecutor {
	return &RPCExecutor{
		CatService:   catService,
		BirdsService: birdsService,
		DogsService:  dogsService,
	}
}

func performRequest(service services.RPCService, ctx context.Context, resultChan chan *services.RPCResponse,
	errChan chan error) {

	rpcRequest := &services.RPCRequest{
		Ctx: ctx,
	}
	rpcResponse, err := service.RPC(rpcRequest)
	if err != nil {
		log.Printf("RPC error calling [%T] [%#v] %s\n", service, rpcResponse, err)
		errChan <- err
		return
	}
	resultChan <- rpcResponse

}

// GetAllPets calls the various pet RPCs in parallel,
// collects all the results and returns them.
// If there is an error, or time report and return.
func (RPCexecutor RPCExecutor) GetAllPets() (result *RPCResult, err error) {

	// Calls a set of RPCs in parallel
	resultChan := make(chan *services.RPCResponse)
	errChan := make(chan error)
	ctx := context.Background()
	ctx, cancelContext := context.WithCancel(ctx)

	rpcCount := 0

	// Call GetAllCats (HTTP)
	go performRequest(RPCexecutor.CatService, ctx, resultChan, errChan)
	rpcCount++

	// Call GetAllBirds (gRPC)
	go performRequest(RPCexecutor.BirdsService, ctx, resultChan, errChan)
	rpcCount++

	// Call GetAllDogs (AMQP)
	go performRequest(RPCexecutor.DogsService, ctx, resultChan, errChan)
	rpcCount++

	// If any of them fails, stop all other RPCs and return error
	// Keep an eye on a overall Timeout
	timeoutChan := time.NewTimer(2 * time.Second)
	pets := map[string]*services.RPCResponse{}
	result = &RPCResult{}
	for {
		select {
		case err = <-errChan:
			cancelContext()
			break
		case rpcResponse := <-resultChan:
			pets[rpcResponse.Key] = rpcResponse
			rpcCount--
		case <-timeoutChan.C:
			err = errors.New("RPCs timeout :(")
			cancelContext()
			break
		}

		if rpcCount == 0 {
			err = populateCats(pets, result)
			if err != nil {
				return
			}
			err = populateBirds(pets, result)
			if err != nil {
				return
			}
			err = populateDogs(pets, result)
			if err != nil {
				return
			}
			break
		}
	}

	return

}

func populateCats(pets map[string]*services.RPCResponse, result *RPCResult) (err error) {

	//TODO: Yep, this is probably a bad idea. Need to make it better
	result.Cats = pets[services.CatServiceKey].Data.([]*catModels.Cat)
	return

}

func populateBirds(pets map[string]*services.RPCResponse, result *RPCResult) (err error) {

	//TODO: Yep, this is probably a bad idea. Need to make it better
	result.Birds = pets[services.BirdServiceKey].Data.([]*birdsModels.Bird)
	return

}

func populateDogs(pets map[string]*services.RPCResponse, result *RPCResult) (err error) {

	//TODO: Yep, this is probably a bad idea. Need to make it better
	result.Dogs = pets[services.DogsServiceKey].Data.([]*dogsModels.Dog)
	return

}
