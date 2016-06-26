package models

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type catsDB struct {
	CatsDataStore
	db *sql.DB
}

// NewCatsDB returns a new CatsDataStore
func NewCatsDB(dbURL string) CatsDataStore {
	db := connectToDB(dbURL)
	return &catsDB{
		db: db,
	}
}

// ConnectToDB connects to the database
func connectToDB(dbURL string) (db *sql.DB) {
	log.Printf("Connecting to DB[%s]....\n", dbURL)
	conn, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return conn

}

func (dbWrapper *catsDB) ReadAllCats() (cats []*Cat, err error) {

	rows, err := dbWrapper.db.Query("select * from cats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cats = []*Cat{}
	for rows.Next() {
		cat := &Cat{}
		err = rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Age,
			&cat.Type,
		)
		if err != nil {
			cats = nil
			return
		}
		cats = append(cats, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (dbWrapper *catsDB) CreateCat(cat *Cat) (err error) {

	stmt, err := dbWrapper.db.Prepare("insert into cats set cat_name=?, cat_age=?, cat_type=?")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(cat.Name, cat.Age, cat.Type)
	if err != nil {
		return
	}
	cat.ID, err = res.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (dbWrapper *catsDB) ReadCat(id int64) (cat *Cat, err error) {

	cat = &Cat{}
	err = dbWrapper.db.QueryRow(
		"select cat_id, cat_name, cat_age, cat_type from cats where cat_id=?", id).
		Scan(
		&cat.ID,
		&cat.Name,
		&cat.Age,
		&cat.Type)
	switch {
	case err == sql.ErrNoRows:
		cat = nil
	}

	return
}

func (dbWrapper *catsDB) UpdateCat(cat *Cat) (err error) {

	stmt, err := dbWrapper.db.Prepare("update cats set cat_name=?, cat_age=?, cat_type=? where cat_id = ?")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(cat.Name, cat.Age, cat.Type, cat.ID)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected != 1 {
		err = errors.New("Unable to update cat")
	}

	return
}

func (dbWrapper *catsDB) DeleteCat(id int64) (err error) {

	stmt, err := dbWrapper.db.Prepare("delete from cats where cat_id=?")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)

	return
}
