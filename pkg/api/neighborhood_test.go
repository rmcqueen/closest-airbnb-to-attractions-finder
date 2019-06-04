package api

import (
	"math"
	"testing"
)

func TestFindNeighborhoodContainingAttraction_noNeighborhoodFound(t *testing.T) {
	attraction := Attraction{"Foobar", "Foobar City", "CA", -32.0, 3.00}

	neighborhood, _ := FindNeighborhoodContainingAttraction(attraction)

	expectedNeighborhoodName := ""
	if neighborhood.Name != expectedNeighborhoodName {
		t.Errorf(
			"Should have returned an empty neighborhood object. Got: %s, expected %s",
			neighborhood.Name,
			expectedNeighborhoodName)
	}
}

func TestFindNeighborhoodContainingAttraction_multipleMatchesExpectedClosestToAttractionReturned(t *testing.T) {
	attraction := Attraction{"Science World", "Vancouver", "BC", 49.2820, -123.1171}

	neighborhood, _ := FindNeighborhoodContainingAttraction(attraction)

	expectedNeighborhoodName := "Downtown"
	if neighborhood.Name != expectedNeighborhoodName {
		t.Errorf(
			"Neighborhood name did not match expected. Got: %s, expected %s",
			neighborhood.Name,
			expectedNeighborhoodName)
	}
}

func TestFindNeighborhoodContainingAttraction_emptyAttractionGiven(t *testing.T) {
	var attraction Attraction

	neighborhood, _ := FindNeighborhoodContainingAttraction(attraction)

	expectedNeighborhoodName := ""
	if neighborhood.Name != expectedNeighborhoodName {
		t.Errorf(
			"Neighborhood name did not match expected. Got: %s, expected %s",
			neighborhood.Name,
			expectedNeighborhoodName)
	}
}

func TestResolveNeighborhoodMultiPolygonsCentroidPoint_neighborhoodNameIsInvalid(t *testing.T) {
	_, err := resolveNeighborhoodMultiPolygonsCentroidPoint("fake", "Vancouver", "BC")

	if err == nil {
		t.Errorf("An exception should have been thrown due to no rows.")
	}
}

func TestResovleNeighborhoodMultiPolygonsCentroidPoint_neighborhoodCentroidResolved(t *testing.T) {
	neighborhoodCoordinates, _ := resolveNeighborhoodMultiPolygonsCentroidPoint("Downtown", "Vancouver", "BC")
	epsilon := 0.0000001
	expectedCoordinates := []float64{-123.116626, 49.280705}
	if math.Abs(neighborhoodCoordinates[0])-math.Abs(expectedCoordinates[0]) > epsilon {
		t.Errorf(
			"Neighborhood latitude is incorrect. Expected: %.6f, got: %.6f",
			expectedCoordinates[0],
			neighborhoodCoordinates[0])
	}
}

func TestGetDistanceBetweenTwoCoordinates_exactSameCoordinatesGiven(t *testing.T) {
	coords1 := []float64{-123.000001, 49.232323}
	res, _ := getDistanceBetweenTwoCoordinates(coords1, coords1)

	expectedCoordinatesDistance := 0.0
	if res != expectedCoordinatesDistance {
		t.Errorf(
			"Exact same coordinates have differing distance. Expected: %.6f, got: %.6f",
			expectedCoordinatesDistance,
			res)
	}
}

func TestFindOptimalNeighborhood_noTies(t *testing.T) {
	g := Graph{}
	nodes := []Neighborhood{
		Neighborhood{"Downtown", "Foobar City", "CA", "USA", -3.1, 0.0},
		Neighborhood{"West Side", "Foobar City", "CA", "USA", -3.2, 0.0},
		Neighborhood{"Central", "Foobar City", "CA", "USA", -3.3, 0.0}}
	g.nodes = nodes
	g.edges = map[string][]Edge{
		nodes[0].Name: {Edge{nodes[0], nodes[1], 3.0}, Edge{nodes[0], nodes[2], 1.0}},
		nodes[1].Name: {Edge{nodes[1], nodes[0], 3.0}, Edge{nodes[1], nodes[2], 5.0}},
		nodes[2].Name: {Edge{nodes[2], nodes[1], 5.0}, Edge{nodes[2], nodes[0], 1.0}}}

	bestNeighborhood, _ := findMinDistanceBetweenNodes(g)

	expectedBestNeighborhood := nodes[0].Name
	if bestNeighborhood.Name != expectedBestNeighborhood {
		t.Errorf(
			"The determined best neighborhood was incorrect. Got: %s, expected: %s.",
			bestNeighborhood.Name,
			expectedBestNeighborhood)
	}
}
