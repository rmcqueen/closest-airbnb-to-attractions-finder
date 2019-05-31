package api

import (
	"testing"

	"github.com/codingsince1985/geo-golang"
)

// StubGeocoder mocks the external calls to a third-party library to ensure no RPCs are done
// during testing.
type StubGeocoder struct {
}

func (sg StubGeocoder) Geocode(address string) (*geo.Location, error) {
	return &geo.Location{Lat: -64.07703, Lng: -3.76949}, nil
}

func (sg StubGeocoder) ReverseGeocode(lat, lng float64) (*geo.Address, error) {
	return &geo.Address{}, nil
}

func Geocoder() geo.Geocoder {
	return StubGeocoder{}
}

func TestMergeAttractionNameCityAndState_allEmptyStrings(t *testing.T) {
	attraction := Attraction{}
	mergedAttractionName, _ := attraction.MergeAttractionNameCityAndState()
	expectedMergedAttractionName := ""

	if mergedAttractionName != expectedMergedAttractionName {
		t.Errorf(
			"Merged attraction name was not expected. Got: %s, expected: %s.",
			mergedAttractionName,
			expectedMergedAttractionName)
	}
}

func TestMergeAttractionNameCityAndState_attractionNameOmitted(t *testing.T) {
	attraction := Attraction{"", "Foobar", "CA", 0.0, 0.0}
	mergedAttractionName, _ := attraction.MergeAttractionNameCityAndState()
	expectedMergedAttractionName := ""

	if mergedAttractionName != expectedMergedAttractionName {
		t.Errorf(
			"Merged attraction name was not expected. Got: %s, expected: %s.",
			mergedAttractionName,
			expectedMergedAttractionName)
	}
}

func TestMergeAttractionNameCityAndState_cityNameOmitted(t *testing.T) {
	attraction := Attraction{"Foobar Bridge", "", "CA", 0.0, 0.0}
	mergedAttractionName, _ := attraction.MergeAttractionNameCityAndState()
	expectedMergedAttractionName := ""

	if mergedAttractionName != expectedMergedAttractionName {
		t.Errorf(
			"Merged attraction name was not expected. Got: %s, expected: %s.",
			mergedAttractionName,
			expectedMergedAttractionName)
	}
}

func TestMergeAttractionNameCityAndState_StateNameOmitted(t *testing.T) {
	attraction := Attraction{"Foobar Bridge", "Foobar City", "", 0.0, 0.0}
	mergedAttractionName, _ := attraction.MergeAttractionNameCityAndState()
	expectedMergedAttractionName := ""

	if mergedAttractionName != expectedMergedAttractionName {
		t.Errorf(
			"Merged attraction name was not expected. Got: %s, expected: %s.",
			mergedAttractionName,
			expectedMergedAttractionName)
	}
}

func TestMergeAttractionNameCityAndState_allAttractionIdentifersPresent(t *testing.T) {
	attraction := Attraction{"Foobar Bridge", "Foobar City", "CA", 0.0, 0.0}
	mergedAttractionName, _ := attraction.MergeAttractionNameCityAndState()
	expectedMergedAttractionName := "Foobar Bridge, Foobar City, CA"

	if mergedAttractionName != expectedMergedAttractionName {
		t.Errorf(
			"Merged attraction name was not expected. Got: %s, expected: %s.",
			mergedAttractionName,
			expectedMergedAttractionName)
	}
}

func TestGeocodeAttraction_locationReturned(t *testing.T) {
	geocoder := StubGeocoder{}
	attraction := Attraction{"Fake Attraction", "Fake City", "CA", 0.0, 0.0}
	location, _ := attraction.GeocodeAttraction(geocoder)

	expectedLatitude := -64.07703
	expectedLongitude := -3.76949

	if location.Lat != expectedLatitude {
		t.Errorf(
			"Location coordinates were unsuccessfully obtained. Got: %.6f, expected: %.6f.",
			location.Lat,
			expectedLatitude)
	}

	if location.Lng != expectedLongitude {
		t.Errorf(
			"Location coordinates were unsuccessfully obtained. Got: %.6f, expected: %.6f.",
			location.Lng,
			expectedLongitude)
	}
}
