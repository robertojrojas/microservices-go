package models

// Dog data struct
type Dog struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
	Type string `json:"type"`
}

// DogsDataStore represents interface to manage dogs
type DogsDataStore interface {
	ReadAllDogs() (dogs []*Dog, err error)
	CreateDog(dog *Dog) (err error)
	ReadDog(id string) (dog *Dog, err error)
	UpdateDog(dog *Dog) (err error)
	DeleteDog(id string) (err error)
}
