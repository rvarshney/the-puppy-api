package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"thepuppyapi/database"
)

func Index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("web/static/templates/index.html")
	randomPuppy := database.GetRandomPuppy("", "")
	t.Execute(res, randomPuppy.ImageUrl)
}

func Puppy(res http.ResponseWriter, req *http.Request) {
	// Parse GET request
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Grab GET form data
	params := req.Form

	// Parse out the breed, if present
	var breed []string
	breed = params["breed"]
	if breed == nil {
		breed = []string{""}
	}

	// Parse out the image type, if present
	var imageType []string
	imageType = params["type"]
	if imageType == nil {
		imageType = []string{""}
	}

	var format []string
	format = params["format"]
	if format == nil {
		format = []string{"json"}
	}

	log.Printf("web.Puppy: breed=%v type=%v format=%v", breed[0], imageType[0], format[0])

	// Get a random puppy with the constraints
	randomPuppy := database.GetRandomPuppy(breed[0], imageType[0])

	// If image source, redirect to the puppy image
	if format[0] == "src" {
		http.Redirect(res, req, randomPuppy.ImageUrl, http.StatusFound)
	}

	// Return JSON
	response := map[string]string{"puppy_url": randomPuppy.ImageUrl}
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func main() {
	// Set up database connection
	log.Println("Setting up DB...")
	database.InitPuppyDB()

	// Set up URL routes
	http.HandleFunc("/", Index)
	http.HandleFunc("/puppy", Puppy)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Start the server
	log.Println("Listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
