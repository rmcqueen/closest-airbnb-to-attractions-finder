package api

import (
	"container/heap"
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

func TestGetMaxHeap_rootNodeCorrectlySet(t *testing.T) {
	frequencyMap := map[string]int{
		"Downtown":   1,
		"South Side": 5,
		"East End":   4,
		"Central":    3,
	}

	h := getMaxHeap(frequencyMap)
	if h.Len() != 4 {
		t.Errorf("Number of heap elements was incorrect. Got: %d, expected: %d.", len(frequencyMap), h.Len())
	}

	rootNode := h.Pop().(neighorboodNameFrequency)
	expectedRootNodeName := "South Side"
	expectedRootNodeCount := 5
	if rootNode.name != expectedRootNodeName {
		t.Errorf("Root node name was incorrect. Got: %s, expected: %s.", rootNode.name, expectedRootNodeName)
	}

	if rootNode.count != expectedRootNodeCount {
		t.Errorf("Root node count was incorrect. Got: %d, expected: %d.", rootNode.count, expectedRootNodeCount)
	}
}

func TestGetMaxHeap_heapIsEmptyWhenEmptyMapGiven(t *testing.T) {
	frequencyMap := map[string]int{}

	h := getMaxHeap(frequencyMap)

	expectedHeapSize := 0
	if h.Len() != expectedHeapSize {
		t.Errorf("Heap size was incorrect. Got: %d, expected: %d.", h.Len(), expectedHeapSize)
	}
}

func TestMaxHeap_elementsSwapCorrectly(t *testing.T) {
	frequencyMap := map[string]int{}

	h := getMaxHeap(frequencyMap)
	heap.Push(h, neighorboodNameFrequency{"foo", 1})
	heap.Push(h, neighorboodNameFrequency{"bar", 2})
	i := 0
	j := 1
	h.Swap(i, j)

	rootNode := h.Pop().(neighorboodNameFrequency)
	expectedRootNodeCount := 1
	if rootNode.count != expectedRootNodeCount {
		t.Errorf("Root node was invalid after swapping. Got: %d, expected: %d.", rootNode.count, expectedRootNodeCount)
	}
}

func TestFindNeighborhoodsWithSameFrequency_onlyOneMaxFrequency(t *testing.T) {
	frequencyMap := map[string]int{
		"Downtown":   1,
		"South Side": 5,
		"East End":   4,
		"Central":    3,
	}

	minHeap := getMaxHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(minHeap)

	expectedNeighborhoodsCount := 1
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}

	expectedNeighborhoodName := "South Side"
	if neighborhoods[0] != expectedNeighborhoodName {
		t.Errorf("The returned neighborhood name was not correct. Got: %s, expected: %s.", neighborhoods[0], expectedNeighborhoodName)
	}
}

func TestFindNeighborhoodsWithSameFrequency_noHeapEntriesGiven(t *testing.T) {
	frequencyMap := map[string]int{}

	minHeap := getMaxHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(minHeap)

	expectedNeighborhoodsCount := 0
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}
}

func TestFindNeighborhoodsWithSameFrequency_oneHeapEntryGiven(t *testing.T) {
	expectedNeighborhoodName := "Downtown"
	frequencyMap := map[string]int{expectedNeighborhoodName: 1}

	minHeap := getMaxHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(minHeap)

	expectedNeighborhoodsCount := 1
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}

	if neighborhoods[0] != expectedNeighborhoodName {
		t.Errorf("The returned neighborhood name was not correct. Got: %s, expected: %s.", neighborhoods[0], expectedNeighborhoodName)
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

// Highly unlikely to ever happen, but still worth testing.
func TestFindOptimalNeighborhood_twoNeighborhoodsTiesForDistance(t *testing.T) {
	g := Graph{}
	nodes := []Neighborhood{
		Neighborhood{"Downtown", "Foobar City", "CA", "USA", -3.1, 0.0},
		Neighborhood{"West Side", "Foobar City", "CA", "USA", -3.2, 0.0},
		Neighborhood{"Central", "Foobar City", "CA", "USA", -3.3, 0.0}}
	g.nodes = nodes
	// Both "A" and "B" are considered to be optimal here.
	g.edges = map[string][]Edge{
		nodes[0].Name: {Edge{nodes[0], nodes[1], 3.0}, Edge{nodes[0], nodes[2], 1.0}},
		nodes[1].Name: {Edge{nodes[1], nodes[0], 3.0}, Edge{nodes[1], nodes[2], 1.0}},
		nodes[2].Name: {Edge{nodes[2], nodes[1], 5.0}, Edge{nodes[2], nodes[0], 1.0}}}

	bestNeighborhood, _ := findMinDistanceBetweenNodes(g)

	expectedOptimalNeighborhoods := map[string]bool{
		nodes[0].Name: true,
		nodes[1].Name: true}
	if expectedOptimalNeighborhoods[bestNeighborhood.Name] == false {
		t.Errorf(
			"The determined optimal neighborhood was incorrect. Got: %s, expected one of: %v.",
			bestNeighborhood.Name,
			expectedOptimalNeighborhoods)
	}
}

func TestFindOptimalNeighborhood_emptyGraphGiven(t *testing.T) {
	g := Graph{}

	bestNeighborhood, _ := findMinDistanceBetweenNodes(g)

	expectedOptimalNeighborhood := ""

	if bestNeighborhood.Name != expectedOptimalNeighborhood {
		t.Errorf(
			"The determined optimal neighborhood was incorrect. Got: %s, expected: %s.",
			bestNeighborhood.Name,
			expectedOptimalNeighborhood)
	}
}
