package models

// BirdRecord data
type BirdRecord struct {
	ID   int64  `json:"id"    db:"bird_id"`
	Name string `json:"name"  db:"bird_name"`
	Age  int32  `json:"age"   db:"bird_age"`
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

// BirdTypeFromString converts a string into a Bird_BirdType
func BirdTypeFromString(birdTypeStr string) (birdType Bird_BirdType) {

	birdType = Bird_UNKNOWN

	switch birdTypeStr {

	case "BLACKBIRD":
		birdType = Bird_BLACKBIRD
	case "BLACKBIRDSCHICKADEE":
		birdType = Bird_BLACKBIRDSCHICKADEE
	case "CHICKADEESCROW":
		birdType = Bird_CHICKADEESCROW
	case "CROWSDOVE":
		birdType = Bird_CROWSDOVE
	case "DOVESDUCK":
		birdType = Bird_DOVESDUCK
	case "DUCKSFINCH":
		birdType = Bird_DUCKSFINCH
	case "FINCHESFLYCATCHER":
		birdType = Bird_FINCHESFLYCATCHER
	case "FLYCATCHERSGAMEBIRD":
		birdType = Bird_FLYCATCHERSGAMEBIRD
	case "GAMEBIRDSGULL":
		birdType = Bird_GAMEBIRDSGULL
	case "GULLSHAWK":
		birdType = Bird_GULLSHAWK
	case "HAWKSHERON":
		birdType = Bird_HAWKSHERON
	case "HERONSHUMMINGBIRD":
		birdType = Bird_HERONSHUMMINGBIRD
	case "HUMMINGBIRDSKINGFISHER":
		birdType = Bird_HUMMINGBIRDSKINGFISHER
	case "KINGFISHERSNUTHATCH":
		birdType = Bird_KINGFISHERSNUTHATCH
	case "NUTHATCHESOWL":
		birdType = Bird_NUTHATCHESOWL
	case "OWLSSHOREBIRD":
		birdType = Bird_OWLSSHOREBIRD
	case "SHOREBIRDSSPARROW":
		birdType = Bird_SHOREBIRDSSPARROW
	case "SPARROWSSWALLOW":
		birdType = Bird_SPARROWSSWALLOW
	case "SWALLOWSTHRUSH":
		birdType = Bird_SWALLOWSTHRUSH
	case "THRUSHESWARBLER":
		birdType = Bird_THRUSHESWARBLER
	case "WARBLERSWOODPECKER":
		birdType = Bird_WARBLERSWOODPECKER
	case "WOODPECKERSWREN":
		birdType = Bird_WOODPECKERSWREN
	case "WRENS":
		birdType = Bird_WRENS

	}

	return
}

// GetBirdTypeStringFromType converts Bird_BirdType to string
func GetBirdTypeStringFromType(birdType Bird_BirdType) (birdTypeStr string) {

	birdTypeStr = "UNKNOWN"
	switch birdType {

	case Bird_BLACKBIRD:
		birdTypeStr = "BLACKBIRD"
	case Bird_BLACKBIRDSCHICKADEE:
		birdTypeStr = "BLACKBIRDSCHICKADEE"
	case Bird_CHICKADEESCROW:
		birdTypeStr = "CHICKADEESCROW"
	case Bird_CROWSDOVE:
		birdTypeStr = "CROWSDOVE"
	case Bird_DOVESDUCK:
		birdTypeStr = "DOVESDUCK"
	case Bird_DUCKSFINCH:
		birdTypeStr = "DUCKSFINCH"
	case Bird_FINCHESFLYCATCHER:
		birdTypeStr = "FINCHESFLYCATCHER"
	case Bird_FLYCATCHERSGAMEBIRD:
		birdTypeStr = "FLYCATCHERSGAMEBIRD"
	case Bird_GAMEBIRDSGULL:
		birdTypeStr = "GAMEBIRDSGULL"
	case Bird_GULLSHAWK:
		birdTypeStr = "GULLSHAWK"
	case Bird_HAWKSHERON:
		birdTypeStr = "HAWKSHERON"
	case Bird_HERONSHUMMINGBIRD:
		birdTypeStr = "HERONSHUMMINGBIRD"
	case Bird_HUMMINGBIRDSKINGFISHER:
		birdTypeStr = "HUMMINGBIRDSKINGFISHER"
	case Bird_KINGFISHERSNUTHATCH:
		birdTypeStr = "KINGFISHERSNUTHATCH"
	case Bird_NUTHATCHESOWL:
		birdTypeStr = "NUTHATCHESOWL"
	case Bird_OWLSSHOREBIRD:
		birdTypeStr = "OWLSSHOREBIRD"
	case Bird_SHOREBIRDSSPARROW:
		birdTypeStr = "SHOREBIRDSSPARROW"
	case Bird_SPARROWSSWALLOW:
		birdTypeStr = "SPARROWSSWALLOW"
	case Bird_SWALLOWSTHRUSH:
		birdTypeStr = "SWALLOWSTHRUSH"
	case Bird_THRUSHESWARBLER:
		birdTypeStr = "THRUSHESWARBLER"
	case Bird_WARBLERSWOODPECKER:
		birdTypeStr = "WARBLERSWOODPECKER"
	case Bird_WOODPECKERSWREN:
		birdTypeStr = "WOODPECKERSWREN"
	case Bird_WRENS:
		birdTypeStr = "WRENS"

	}

}
