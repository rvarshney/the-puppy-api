package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

type Puppy struct {
	ImageUrl  string
	ImageType string
	Breed     string
}

var db *sql.DB

func InitPuppyDB() {
	if db == nil {
		// Get DB settings from OS env
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPwd := os.Getenv("DB_PWD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		dbStr := "host='" + dbHost + "' user='" + dbUser + "' password='" + dbPwd + "' dbname='" + dbName + "' port='" + dbPort + "'"

		var err error
		db, err = sql.Open("postgres", dbStr)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(db)
}

func GetRandomPuppy(breed string, imageType string) *Puppy {
	// Perform LIKE query to pick a random puppy. Yes, this is pretty darn inefficient.
	rows, err := db.Query(`SELECT image_url, image_type, breed FROM public.thepuppyapi_puppies WHERE breed LIKE $1 AND image_type LIKE $2 ORDER BY random() LIMIT 1`, "%"+strings.ToLower(breed)+"%", "%"+strings.ToLower(imageType)+"%")
	if err != nil {
		log.Fatal(err)
	}

	// Extract results
	var resultImageUrl string
	var resultImageType string
	var resultBreed string
	for rows.Next() {
		err = rows.Scan(&resultImageUrl, &resultImageType, &resultBreed)
	}
	return &Puppy{ImageUrl: resultImageUrl, ImageType: resultImageType, Breed: resultBreed}
}

func AddPuppy(puppy *Puppy) {
	// Attempts to perform an insert into the puppy database
	log.Printf("database.AddPuppy: Trying to add image_url=%v image_type=%v breed=%v", puppy.ImageUrl, puppy.ImageType, puppy.Breed)
	_, err := db.Query(`INSERT INTO public.thepuppyapi_puppies (image_url, image_type, breed) VALUES ($1, $2, $3)`, puppy.ImageUrl, puppy.ImageType, strings.ToLower(puppy.Breed))
	if err != nil {
		log.Printf("database.AddPuppy: %v", err)
	}
}
