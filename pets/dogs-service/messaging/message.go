package messaging

/*
    ReadAllDogs() (dogs []*Dog, err error)
	CreateDog(dog *Dog) (err error)
	ReadDog(id string) (dog *Dog, err error)
	UpdateDog(dog *Dog) (err error)
    DeleteDog(id string) (err error)
*/

type RPCMessageType int

const (
	ReadAllMessage RPCMessageType = iota
	CreateMessage
	ReadMessage
	UpdateMessage
	DeleteMessage
	ErrorMessage
	UnknownMessage
)

type RPCMessage struct {
	Type RPCMessageType
	Data []byte
}
