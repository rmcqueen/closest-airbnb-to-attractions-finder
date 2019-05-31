package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"../pkg/api"
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

// AttractionsResponse demonstrates the components involved for API responses.
type AttractionsResponse struct {
	SuccessfulAttractions []api.Attraction `json:"successful_attractions"`
	FailedAttractions     []api.Attraction `json:"failed_attractions"`
	ClosestNeighborhood   api.Neighborhood `json:"closest_neighborhood"`
}

func server() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	http.HandleFunc("/attractions", handler)
	server()

}

func handler(w http.ResponseWriter, r *http.Request) {
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}

	var attractions []api.Attraction
	err = json.Unmarshal(jsn, &attractions)
	if err != nil {
		log.Fatal("Decoding error", err)
	}

	var responseAttractions AttractionsResponse

	var neighborhoods []api.Neighborhood
	geocoder := openstreetmap.Geocoder()
	for _, attraction := range attractions {
		attractionLocation, _ := attraction.GeocodeAttraction(geocoder)

		if attractionLocation == nil {
			responseAttractions.FailedAttractions = append(responseAttractions.FailedAttractions, attraction)
			continue
		}

		attraction.Latitude = attractionLocation.Lat
		attraction.Longitude = attractionLocation.Lng
		responseAttractions.SuccessfulAttractions = append(responseAttractions.SuccessfulAttractions, attraction)
		neighborhood, err := api.FindNeighborhoodContainingAttraction(attraction)
		if err != nil {
			log.Fatal(err)
			continue
		}

		neighborhoods = append(neighborhoods, neighborhood)
	}

	// TODO: for each neighborhood retrieved for each attraction, find the closest neighborhood to them all
	api.FindBestNeighborhood(neighborhoods)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseAttractions)
}
