package api

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Cat data
type Cat struct {
	ID   int64  `json:"id"    db:"cat_id"`
	Name string `json:"name"  db:"cat_name"`
	Age  int    `json:"age"   db:"cat_age"`
	Type string `json:"type"  db:"cat_type"`
}

// CatsDataStore represents interface to manage cats
type CatsDataStore interface {
	ReadAllCats() (cats []*Cat, err error)
	CreateCat(cat *Cat) (err error)
	ReadCat(id int64) (cat *Cat, err error)
	UpdateCat(cat *Cat) (err error)
	DeleteCat(id int64) (err error)
}

type catsDB struct {
	CatsDataStore
	db *sql.DB
}

// NewCatsDB returns a new CatsDataStore
func NewCatsDB() (CatsDataStore, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	return &catsDB{
		db: db,
	}, nil
}

const (
	customTLS = "customTLS"
)

type config struct {
	mySQLDBHost    string
	mySQLDBPort    string
	mySQLDBUser    string
	mySQLDBPass    string
	mySQLDBName    string
	mySQLDBTLSCert string
}

func getConfig() *config {

	envConfig := config{}
	envConfig.mySQLDBHost = os.Getenv("MYSQL_HOST")
	envConfig.mySQLDBPort = os.Getenv("MYSQL_PORT")
	envConfig.mySQLDBName = os.Getenv("MYSQL_DBNAME")
	envConfig.mySQLDBUser = os.Getenv("MYSQL_USERNAME")
	envConfig.mySQLDBPass = os.Getenv("MYSQL_PASSWORD")
	envConfig.mySQLDBTLSCert = os.Getenv("MYSQL_CERT")

	return &envConfig
}

func buildMYSQLConnString(envConfig *config) string {
	mysqlURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		envConfig.mySQLDBUser,
		envConfig.mySQLDBPass,
		envConfig.mySQLDBHost,
		envConfig.mySQLDBPort,
		envConfig.mySQLDBName,
	)
	if config.mySQLDBTLSCert != "" {
		mysqlURL = fmt.Sprintf("%s?tls=%s", mysqlURL, customTLS)
	}
}

func setupTLS(envConfig *config) {
	if config.mySQLDBTLSCert != "" {
		rootCertPool := getSSLCert(config)
		mysql.RegisterTLSConfig(customTLS, &tls.Config{
			RootCAs: rootCertPool,
		})
	}
}

func getDBConnection() (*sql.DB, error){
	config := getConfig()
	setupTLS(config)
	dbURL := buildMYSQLConnString(config)
	log.Printf("Connecting to DB[%s]....\n", dbURL)
	theDB, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
}

// ConnectToDB connects to the database
func connectToDB() (db *sql.DB, error) {
	theDB, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	//Connection pool sample configuration
	theDB.SetMaxIdleConns(2)
	theDB.SetMaxOpenConns(10)

	err = theDB.Ping()
	if err != nil {
		return nil, err
	}

	return theDB

}

func (dbWrapper *catsDB) ReadAllCats() (cats []*Cat, err error) {

	rows, err := dbWrapper.db.Query("select cat_id, cat_name, cat_age, cat_type from cats")
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
