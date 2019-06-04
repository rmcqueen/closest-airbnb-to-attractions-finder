package api

import "testing"

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
