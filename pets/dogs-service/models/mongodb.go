package models

import (
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

func NewDogMongoStore(mongoDBURL, DBName, CollectionName string) (*DogMongoStore, error) {
	session, err := connectToMongoDB(mongoDBURL)
	if err != nil {
		return nil, err
	}
	return &DogMongoStore{
		Session:        session,
		DBName:         DBName,
		CollectionName: CollectionName,
	}, nil
}

func connectToMongoDB(mongoDBURL string) (session *mgo.Session, err error) {
	session, err = mgo.Dial(mongoDBURL)
	return

}

func (dataStore *DogMongoStore) Disconnect() {
	dataStore.Session.Close()
}

func (dataStore *DogMongoStore) ReadAllDogs() (dogs []*Dog, err error) {

	dogsCollection := dataStore.Session.DB(dataStore.DBName).C(dataStore.CollectionName)
	err = dogsCollection.Find(bson.M{}).All(&dogs)
	if err != nil {
		return
	}

	return
}

func (dataStore *DogMongoStore) CreateDog(dog *Dog) (err error) {

	dogsCollection := dataStore.Session.DB(dataStore.DBName).C(dataStore.CollectionName)
	err = dogsCollection.Insert(dog)
	return
}

func (dataStore *DogMongoStore) ReadDog(id string) (dog *Dog, err error) {

	dogsCollection := dataStore.Session.DB(dataStore.DBName).C(dataStore.CollectionName)
	dog = &Dog{}
	err = dogsCollection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(dog)
	return

}

func (dataStore *DogMongoStore) UpdateDog(dog *Dog) (err error) {
	return
}

func (dataStore *DogMongoStore) DeleteDog(id string) (err error) {
	return
}
