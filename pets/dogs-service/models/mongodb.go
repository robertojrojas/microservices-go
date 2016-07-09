package models

// DogMongoStore implements DogsDataStore
type DogMongoStore struct {
	DogsDataStore
}

func (dataStore *DogMongoStore) ReadAllDogs() (dogs []*Dog, err error) {
	return
}

func (dataStore *DogMongoStore) CreateDog(dog *Dog) (err error) {
	return
}

func (dataStore *DogMongoStore) ReadDog(id string) (dog *Dog, err error) {
	return
}

func (dataStore *DogMongoStore) UpdateDog(dog *Dog) (err error) {
	return
}

func (dataStore *DogMongoStore) DeleteDog(id string) (err error) {
	return
}
