package main

import (
	"code.google.com/p/google-api-go-client/customsearch/v1"
	"code.google.com/p/google-api-go-client/googleapi/transport"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
	"thepuppyapi/database"
)

func GetPuppySearchHandle(breed string) *customsearch.CseListCall {
	// Set up the API transport
	transport := &transport.APIKey{
		Key:       os.Getenv("SEARCH_API_KEY"),
		Transport: http.DefaultTransport,
	}

	// Get HTTP client
	httpClient := &http.Client{Transport: transport}

	// Instantiate custom search service
	searchService, err := customsearch.New(httpClient)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Set up the search parameters for the service
	searchCx := os.Getenv("SEARCH_CX")
	searchExactTerms := "puppy"
	searchType := "image"
	searchImageSize := "large"
	searchQuery := breed + " puppy pictures"
	searchStartIndex := int64(30)

	// Get search handle
	searchHandle := searchService.Cse.List(searchQuery).Cx(searchCx).ExactTerms(searchExactTerms).SearchType(searchType).Start(searchStartIndex).ImgSize(searchImageSize)
	return searchHandle
}

func GetPuppySearchResults(breed string) []*customsearch.Result {
	// Get puppy search handle
	searchHandle := GetPuppySearchHandle(breed)
	if searchHandle != nil {
		search, err := searchHandle.Do()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		return search.Items
	}
	return nil
}

func AddPuppySearchResults(breed string, searchResults []*customsearch.Result) {
	if searchResults == nil {
		return
	}

	log.Printf("finder.AddPuppySearchResults: Found %v search results", len(searchResults))

	for _, result := range searchResults {
		imageUrl := result.Link
		imageType := strings.Split(result.Mime, "/")[1]
		database.AddPuppy(&database.Puppy{ImageUrl: imageUrl, ImageType: imageType, Breed: breed})
	}
}

func main() {
	database.InitPuppyDB()

	breed := "labrador"
	searchResults := GetPuppySearchResults(breed)
	AddPuppySearchResults(breed, searchResults)
}
