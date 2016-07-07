package models

// BirdRecord data
type BirdRecord struct {
	ID   int64  `json:"id"    db:"bird_id"`
	Name string `json:"name"  db:"bird_name"`
	Age  int    `json:"age"   db:"bird_age"`
	Type string `json:"type"  db:"bird_type"`
}

// BirdsDataStore represents interface to manage birds
type BirdsDataStore interface {
	ReadAllBirds() (birds []*BirdRecord, err error)
	CreateBird(bird *BirdRecord) (err error)
	ReadBird(id int64) (bird *BirdRecord, err error)
	UpdateBird(bird *BirdRecord) (err error)
	DeleteBird(id int64) (err error)
}