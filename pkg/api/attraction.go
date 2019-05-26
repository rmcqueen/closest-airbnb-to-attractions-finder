package api

import (
    "strings"
    "github.com/codingsince1985/geo-golang"
    "github.com/codingsince1985/geo-golang/openstreetmap"
)

type Attraction struct {
    Name string `json:"name"`
    City string `json:"city"`
    StateOrProvinceName string `json:"state_or_province_name"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func (attraction *Attraction) MergeAttractionNameCityAndState() string {
    return strings.Join([]string{attraction.Name, attraction.City, attraction.StateOrProvinceName}, ", ")
}

func (attraction *Attraction) GeocodeAttraction() *geo.Location {
    mergedAttraction := attraction.MergeAttractionNameCityAndState()
    geocoder := openstreetmap.Geocoder()
    location, _ := geocoder.Geocode(mergedAttraction)

    return location;
}

