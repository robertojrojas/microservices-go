package main

// Cat data
type Cat struct {
	ID   int64  `json:"id"    db:"cat_id"`
	Name string `json:"name"  db:"cat_name"`
	Age  int    `json:"age"   db:"cat_age"`
	Type string `json:"type"  db:"cat_type"`
}
