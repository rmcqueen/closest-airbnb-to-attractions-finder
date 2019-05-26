package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "../pkg/api"
)

type Neighborhood struct {
    Name string `json:"name"`
    City string `json:"city_name"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

type AttractionsResponse struct {
    SuccessfulAttractions []api.Attraction `json:"successful_attractions"`
    FailedAttractions []api.Attraction `json:"failed_attractions"`
    ClosestNeighborhood Neighborhood
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
    for _, attraction := range attractions {
        attractionLocation := attraction.GeocodeAttraction()

        if attractionLocation == nil {
            responseAttractions.FailedAttractions = append(responseAttractions.FailedAttractions, attraction)
        }

        attraction.Latitude = attractionLocation.Lat
        attraction.Longitude = attractionLocation.Lng
        responseAttractions.SuccessfulAttractions = append(responseAttractions.SuccessfulAttractions, attraction)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(responseAttractions)
}

