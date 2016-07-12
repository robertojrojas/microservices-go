package executor

import (
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/net/context"

	catModels "github.com/robertojrojas/microservices-go/pets/cats-service/models"
	"github.com/robertojrojas/microservices-go/pets/pets-service/services"
)

// RPCExecutor is in charge of performing all the RPC requests.
type RPCExecutor struct {
	// Cat Service
	CatService *services.CatService
	// Bird Service
	// Dog Service
}

type RPCResult struct {
	Cats []*catModels.Cat `json:"cats"`
}

func NewRPCExecutor(catService *services.CatService) *RPCExecutor {
	return &RPCExecutor{
		CatService: catService,
	}
}

func performRequest(service services.RPCService, ctx context.Context, resultChan chan *services.RPCResponse,
	errChan chan error) {

	rpcRequest := &services.RPCRequest{
		Ctx: ctx,
	}
	rpcResponse, err := service.RPC(rpcRequest)
	if err != nil {
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

	rpcCount := 0

	// Call GetAllCats (HTTP)
	go performRequest(RPCexecutor.CatService, ctx, resultChan, errChan)
	rpcCount++

	// Call GetAllBirds (gRPC)
	//rpcCount++

	// Call GetAllDogs (AMQP)
	//rpcCount++

	// If any of them fails, stop all other RPCs and return error
	// Keep an eye on a overall Timeout
	timeoutChan := time.NewTimer(2 * time.Second)
	pets := map[string]*services.RPCResponse{}
	result = &RPCResult{}
	for {
		select {
		case err = <-errChan:
			break
		case rpcResponse := <-resultChan:
			pets[rpcResponse.Key] = rpcResponse
			rpcCount--
		case <-timeoutChan.C:
			err = errors.New("RPCs timeout :(")
			break
		}

		if rpcCount == 0 {
			populateCats(pets, result)
			break
		}
	}

	return

}

func populateCats(pets map[string]*services.RPCResponse, result *RPCResult) (err error) {

	cats := []*catModels.Cat{}
	err = json.Unmarshal(pets[services.CatServiceKey].Data, &cats)
	if err != nil {
		return
	}
	result.Cats = cats
	return

}
