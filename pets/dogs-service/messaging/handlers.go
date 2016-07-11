package messaging

import (
	"encoding/json"
	"log"

	"github.com/robertojrojas/microservices-go/pets/dogs-service/models"
)

type MessageHandler struct {
	dataStore *models.DogMongoStore
}

var unknownMessage *RPCMessage

func init() {
	unknownMessage = &RPCMessage{
		Type: UnknownMessage,
		Data: []byte("Unknown Message received"),
	}
}

func NewMessageHandler(dataStore *models.DogMongoStore) *MessageHandler {

	return &MessageHandler{
		dataStore: dataStore,
	}
}

func (messageHandler *MessageHandler) ReadAllDogs(rpcIn *RPCMessage) (rpcOut *RPCMessage, err error) {

	dogs, err := messageHandler.dataStore.ReadAllDogs()
	if err != nil {
		return
	}

	dogsData, err := json.Marshal(dogs)
	if err != nil {
		return
	}

	rpcOut = &RPCMessage{
		Type: ReadAllMessage,
		Data: dogsData,
	}

	return
}

func (messageHandler *MessageHandler) CreateDog(rpcIn *RPCMessage) (rpcOut *RPCMessage, err error) {

	dog := &models.Dog{}
	err = json.Unmarshal(rpcIn.Data, dog)
	if err != nil {
		return
	}

	err = messageHandler.dataStore.CreateDog(dog)
	if err != nil {
		return
	}

	dogData, err := json.Marshal(dog)
	if err != nil {
		return
	}

	rpcOut = &RPCMessage{
		Type: CreateMessage,
		Data: dogData,
	}

	return
}

func (messageHandler *MessageHandler) ReadDog(rpcIn *RPCMessage) (rpcOut *RPCMessage, err error) {

	dog, err := messageHandler.dataStore.ReadDog(string(rpcIn.Data))
	if err != nil {
		return
	}

	dogData, err := json.Marshal(dog)
	if err != nil {
		return
	}

	rpcOut = &RPCMessage{
		Type: CreateMessage,
		Data: dogData,
	}

	return
}

func (messageHandler *MessageHandler) UpdateDog(rpcIn *RPCMessage) (rpcOut *RPCMessage, err error) {

	return
}

func (messageHandler *MessageHandler) DeleteDog(rpcIn *RPCMessage) (rpcOut *RPCMessage, err error) {

	return
}

// RouteMessage routes message to appropriate message handler.
// When errors occur returns an ErrorMessage
func (messageHandler *MessageHandler) RouteMessage(message []byte) (rpcOut *RPCMessage) {

	rpcIn := &RPCMessage{}
	err := json.Unmarshal(message, rpcIn)
	if err != nil {
		rpcOut = createErrorMessage(err)
		return
	}

	//TODO: refactor to use a more flexible mapping
	//      and avoid DRY
	switch rpcIn.Type {
	case ReadAllMessage:
		log.Printf("handling ReadAllMessage...")
		rpcOut, err = messageHandler.ReadAllDogs(rpcIn)
		if err != nil {
			rpcOut = createErrorMessage(err)
		}
	case CreateMessage:
		log.Printf("handling CreateMessage...")
		rpcOut, err = messageHandler.CreateDog(rpcIn)
		if err != nil {
			rpcOut = createErrorMessage(err)
		}
	case ReadMessage:
		log.Printf("handling ReadMessage...")
		rpcOut, err = messageHandler.ReadDog(rpcIn)
		if err != nil {
			rpcOut = createErrorMessage(err)
		}
	case UpdateMessage:
		log.Printf("handling UpdateMessage...")
		rpcOut, err = messageHandler.UpdateDog(rpcIn)
		if err != nil {
			rpcOut = createErrorMessage(err)
		}
	case DeleteMessage:
		log.Printf("handling DeleteMessage...")
		rpcOut, err = messageHandler.DeleteDog(rpcIn)
		if err != nil {
			rpcOut = createErrorMessage(err)
		}
	default:
		log.Printf("handling UnknownMessage...")
		rpcOut = unknownMessage
	}

	log.Printf("returning %#v\n", rpcOut.Type)

	return
}

func createErrorMessage(err error) (errorMessage *RPCMessage) {

	errorMessage = &RPCMessage{
		Type: ErrorMessage,
		Data: []byte(err.Error()),
	}

	return
}
