package models

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DogMongoStore implements DogsDataStore
type DogMongoStore struct {
	DogsDataStore
	DBName         string
	CollectionName string
	Session        *mgo.Session
}

func NewDogMongoStore(mongoDBURL, DBName, CollectionName string) *DogMongoStore {
	session := connectToMongoDB(mongoDBURL)
	return &DogMongoStore{
		Session:        session,
		DBName:         DBName,
		CollectionName: CollectionName,
	}
}

func connectToMongoDB(mongoDBURL string) (session *mgo.Session) {
	session, err := mgo.Dial(mongoDBURL)
	if err != nil {
		log.Fatal(err)
	}

	return

}

func (dataStore *DogMongoStore) Disconnect() {
	dataStore.Session.Close()
}

func (dataStore *DogMongoStore) ReadAllDogs() (dogs []*Dog, err error) {

	dogsCollection := dataStore.Session.DB(dataStore.DBName).C(dataStore.CollectionName)
	var result []struct {
		ID   bson.ObjectId `bson:"_id"`
		Name string
		Age  int32
		Type string
	}
	err = dogsCollection.Find(bson.M{}).All(&result)
	if err != nil {
		return
	}

	for _, temp := range result {
		dog := &Dog{
			ID:   temp.ID.Hex(),
			Name: temp.Name,
			Age:  temp.Age,
			Type: temp.Type,
		}
		dogs = append(dogs, dog)
	}

	return
}

func (dataStore *DogMongoStore) CreateDog(dog *Dog) (err error) {
	dogsCollection := dataStore.Session.DB(dataStore.DBName).C(dataStore.CollectionName)
	err = dogsCollection.Insert(dog)
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
