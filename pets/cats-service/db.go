package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var dbURL string

// ConnectToDB connects to the database
func ConnectToDB() {
	conn, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db = conn
}

func readAllCats() (cats []*Cat, err error) {

	rows, err := db.Query("select * from cats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func createCat(cat *Cat) (err error) {

	stmt, err := db.Prepare("insert into cats set cat_name=?, cat_age=?, cat_type=?")
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

func readCat(id int64) (cat *Cat, err error) {

	cat = &Cat{}
	err = db.QueryRow(
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

func updateCat(cat *Cat) (err error) {

	stmt, err := db.Prepare("update cats set cat_name=?, cat_age=?, cat_type=? where cat_id = ?")
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

func deleteCat(id int64) (err error) {

	stmt, err := db.Prepare("delete from cats where cat_id=?")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)

	return
}
