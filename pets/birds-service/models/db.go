package models

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type birdsDB struct {
	BirdsDataStore
	db *sql.DB
}

// NewBirdsDB returns a new BirdsDataStore
func NewBirdsDB(dbURL string) BirdsDataStore {
	db := connectToDB(dbURL)
	return &birdsDB{
		db: db,
	}
}

// ConnectToDB connects to the database
func connectToDB(dbURL string) (db *sql.DB) {
	log.Printf("Connecting to DB[%s]....\n", dbURL)
	theDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	//Connection pool sample configuration
	theDB.SetMaxIdleConns(2)
	theDB.SetMaxOpenConns(10)

	err = theDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return theDB

}

func (dbWrapper *birdsDB) ReadAllBirds() (birds []*BirdRecord, err error) {

	rows, err := dbWrapper.db.Query("select bird_id, bird_name, bird_age, bird_type from birds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	birds = []*BirdRecord{}
	for rows.Next() {
		bird := &BirdRecord{}
		err = rows.Scan(
			&bird.ID,
			&bird.Name,
			&bird.Age,
			&bird.Type,
		)
		if err != nil {
			birds = nil
			return
		}
		birds = append(birds, bird)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (dbWrapper *birdsDB) CreateBird(bird *BirdRecord) (err error) {

	stmt, err := dbWrapper.db.Prepare("insert into birds set bird_name=?, bird_age=?, bird_type=$1")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(bird.Name, bird.Age, bird.Type)
	if err != nil {
		return
	}
	bird.ID, err = res.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (dbWrapper *birdsDB) ReadBird(id int64) (bird *BirdRecord, err error) {

	bird = &BirdRecord{}
	err = dbWrapper.db.QueryRow(
		"select bird_id, bird_name, bird_age, bird_type from birds where bird_id=$1", id).
		Scan(
		&bird.ID,
		&bird.Name,
		&bird.Age,
		&bird.Type)
	switch {
	case err == sql.ErrNoRows:
		bird = nil
	}

	return
}

func (dbWrapper *birdsDB) UpdateCat(bird *BirdRecord) (err error) {

	stmt, err := dbWrapper.db.Prepare("update birds set bird_name=?, bird_age=?, bird_type=$1 where bird_id = $2")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(bird.Name, bird.Age, bird.Type, bird.ID)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected != 1 {
		err = errors.New("Unable to update bird")
	}

	return
}

func (dbWrapper *birdsDB) DeleteBird(id int64) (err error) {

	stmt, err := dbWrapper.db.Prepare("delete from birds where bird_id=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)

	return
}
