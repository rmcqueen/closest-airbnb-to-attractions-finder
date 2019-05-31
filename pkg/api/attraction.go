package api

import (
	"log"
	"strings"

	"github.com/codingsince1985/geo-golang"
)

// Attraction holds the information related to the geographic location of
// a real-world location (i.e, a famous restaurant.)
type Attraction struct {
	Name                string  `json:"name"`
	City                string  `json:"city"`
	StateOrProvinceName string  `json:"state_or_province_name"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}

// MissingAttractionKeyIdentifierError indicates a key identifying piece for the attractio is missing.
type MissingAttractionKeyIdentifierError struct {
	message string
}

func (e *MissingAttractionKeyIdentifierError) Error() string {
	return e.message
}

// MergeAttractionNameCityAndState combines the name of the attraction, city, and state separated by a comma and space.
func (attraction *Attraction) MergeAttractionNameCityAndState() (string, error) {
	joinedString := strings.Join([]string{attraction.Name, attraction.City, attraction.StateOrProvinceName}, ", ")
	if strings.Contains(joinedString, ", , ") {
		return "", &MissingAttractionKeyIdentifierError{"Missing one of: attraction name, city name, or city name."}
	} else if strings.HasPrefix(joinedString, ", ") {
		return "", &MissingAttractionKeyIdentifierError{"Missing city name."}
	} else if strings.HasSuffix(joinedString, ", ") {
		return "", &MissingAttractionKeyIdentifierError{"Missing all of: attraction name, or city name, or state name."}
	} else {
		return joinedString, nil
	}
}

// GeocodeAttraction takes the conjoined attraction's name and obtains the lat/lng coordinates for it.
// A conjoined address is ATTRACTION_NAME, CITY, STATE. Country is omitted for now.
func (attraction *Attraction) GeocodeAttraction(geocoder geo.Geocoder) (*geo.Location, error) {
	mergedAttraction, err := attraction.MergeAttractionNameCityAndState()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	location, _ := geocoder.Geocode(mergedAttraction)

	return location, nil
}
