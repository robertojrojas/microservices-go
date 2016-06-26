package models

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
