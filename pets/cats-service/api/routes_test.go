package api


import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"errors"
	"strings"
)



func TestShouldReadAllCats(t *testing.T) {

	cats := []*Cat {
		&Cat{
			ID: 1,
			Age: 1,
			Name: "catty",
			Type: "funny",
		},
		&Cat{
			ID: 2,
			Age: 2,
			Name: "catty2",
			Type: "funny2",
		},
	}
	dataStore := &testCatsDataStore{
		cats:cats,
	}
	cr := &CatsRoutes{
		catsDBStore: dataStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/cats", errorHandler(cr.ReadAllHandler))

	writer := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/api/cats", nil)

	if err != nil {
		t.Error(err.Error())
	}

	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Errorf("Got error status code: %d ", writer.Code)
	}

	var receivedCats []*Cat
	json.Unmarshal(writer.Body.Bytes(), &receivedCats)
	if receivedCats == nil {
		t.Error("Expects cats")
	}

	if len(cats) != len(receivedCats) {
		t.Errorf("Expected %d cats, but instead got %d cats", len(cats), len(receivedCats))
	}


}

func TestReadAllFails(t *testing.T) {

	const errorMessage string = "Unable to get cats"

	dataStore := &testCatsDataStore{
		err: errors.New(errorMessage),
	}
	cr := &CatsRoutes{
		catsDBStore: dataStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/cats", errorHandler(cr.ReadAllHandler))

	writer := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/api/cats", nil)

	if err != nil {
		t.Error(err.Error())
	}
	mux.ServeHTTP(writer, request)


	if writer.Code != http.StatusInternalServerError {
		t.Errorf("Got error status code: %d ", writer.Code)
	}

	errStr := strings.Trim(string(writer.Body.Bytes()), "\n")
	if errStr != errorMessage {
		t.Errorf("Expected %q, but instead got %q", errorMessage, errStr)
	}

}



func TestShouldCreateCat(t *testing.T) {
	dataStore := &testCatsDataStore{

	}
	cr := &CatsRoutes{
		catsDBStore: dataStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/cats", errorHandler(cr.CreateHandler))

	writer := httptest.NewRecorder()
	json := strings.NewReader(`
		{
		   "cat_id": 1,
		   "cat_name": "test_cat",
		   "cat_age": 3,
		   "cat_type": "funny"
		}`)
	request, err := http.NewRequest("POST", "/api/cats", json)

	if err != nil {
		t.Error(err.Error())
	}
	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusCreated {
		t.Errorf("Expected error code %d, but instead got %d", http.StatusCreated, writer.Code)
	}

}

func TestCreateCatFails(t *testing.T) {

	const errorMessage string = "Unable to get cats"
	dataStore := &testCatsDataStore{
		err:errors.New(errorMessage),
	}
	cr := &CatsRoutes{
		catsDBStore: dataStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/cats", errorHandler(cr.CreateHandler))

	writer := httptest.NewRecorder()
	json := strings.NewReader(`
		{
		   "cat_id": 1,
		   "cat_name": "test_cat",
		   "cat_age": 3,
		   "cat_type": "funny"
		}`)
	request, err := http.NewRequest("POST", "/api/cats", json)

	if err != nil {
		t.Error(err.Error())
	}
	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusInternalServerError {
		t.Errorf("Got error status code: %d ", writer.Code)
	}

	errStr := strings.Trim(string(writer.Body.Bytes()), "\n")
	if errStr != errorMessage {
		t.Errorf("Expected %q, but instead got %q", errorMessage, errStr)
	}

}

func TestReadHandler_OK(t *testing.T) {

}

func TestUpdateHandler_OK(t *testing.T) {

}


func TestDeleteHandler_OK(t *testing.T) {

}

type testCatsDataStore struct {
	cat *Cat
	cats []*Cat
	err error
}


func (dbWrapper *testCatsDataStore) ReadAllCats() (cats []*Cat, err error) {
	return dbWrapper.cats, dbWrapper.err
}

func (dbWrapper *testCatsDataStore) ReadCat(id int64) (cat *Cat, err error) {
	return dbWrapper.cat, dbWrapper.err
}

func (dbWrapper *testCatsDataStore) CreateCat(cat *Cat) (err error) {
	return dbWrapper.err
}

func (dbWrapper *testCatsDataStore) UpdateCat(cat *Cat) (err error) {
	return dbWrapper.err
}
func (dbWrapper *testCatsDataStore) DeleteCat(id int64) (err error) {
	return dbWrapper.err
}