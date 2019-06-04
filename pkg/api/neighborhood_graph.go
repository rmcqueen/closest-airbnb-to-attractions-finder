package api

import (
	"log"
	"math"
)

// Edge denotes a connection between two Neighborhood nodes.
type Edge struct {
	sourceNode       Neighborhood
	targetNode       Neighborhood
	distanceInMeters float64
}

// Graph stores all Neighborhoods and their connections between each other.
type Graph struct {
	nodes []Neighborhood
	edges map[string][]Edge
}

func (graph Graph) buildGraphFromNeighborhoods(neighborhoods []Neighborhood) (Graph, error) {
	for _, neighborhood := range neighborhoods {
		graph.nodes = append(graph.nodes, neighborhood)
	}

	return graph, nil
}

func findNeighborhoodWithLeastDistanceToAllOtherNeighborhoods(neighborhoods []Neighborhood) (Neighborhood, error) {
	var graph Graph
	// Ideally, this would be a thread-safe cache to deal with concurrent requests (i.e, Redis).
	distanceCache := make(map[string]float64)

	for _, neighborhood := range neighborhoods {
		sourceNode := neighborhood
		graph.nodes = append(graph.nodes, sourceNode)
		remainingNeighborhoods := composeDifferingNeighborhoodNamesSlice(neighborhood.Name, neighborhoods)
		for _, otherNeighborhood := range remainingNeighborhoods {
			targetNode := neighborhood

			var distanceInMeters float64
			hashedString := generateNeighborhoodCacheKey(neighborhood.Name, otherNeighborhood.Name)
			_, ok := distanceCache[hashedString]

			if ok == false {
				distanceInMeters, _ = getDistanceBetweenTwoCoordinates([]float64{neighborhood.Longitude, neighborhood.Latitude}, []float64{otherNeighborhood.Longitude, otherNeighborhood.Latitude})
				distanceCache[hashedString] = distanceInMeters
			} else {
				distanceInMeters = distanceCache[hashedString]
			}

			edge := Edge{sourceNode, targetNode, distanceInMeters}
			graph.edges[neighborhood.Name] = append(graph.edges[neighborhood.Name], edge)
		}
	}

	optimalNeighborhood, err := findMinDistanceBetweenNodes(graph)
	if err != nil {
		log.Printf("Error after finding optimal neighborhood: %v\n", err)
		return Neighborhood{}, err
	}

	return optimalNeighborhood, nil
}

func composeDifferingNeighborhoodNamesSlice(currentNeighborhoodName string, allNeighborhoodNames []Neighborhood) []Neighborhood {
	var newSlice []Neighborhood
	for _, neighborhood := range allNeighborhoodNames {
		if currentNeighborhoodName != neighborhood.Name {
			newSlice = append(newSlice, neighborhood)
		}
	}

	return newSlice
}

// Searches the constructed graph for the neighborhood with min distance between all other points.
// Time complexity is O(V*E) where V represents the number of vertices to visit, and E represents the
// number of edges to examine.
func findMinDistanceBetweenNodes(graph Graph) (Neighborhood, error) {
	if len(graph.nodes) == 1 {
		return graph.nodes[0], nil
	}

	neighborhoodDistanceSums := make(map[string]float64)
	for sourceNode, edges := range graph.edges {
		_, ok := neighborhoodDistanceSums[sourceNode]
		if ok == true {
			neighborhoodDistanceSums[sourceNode] = 0
		}
		for _, targetNode := range edges {
			neighborhoodDistanceSums[sourceNode] += targetNode.distanceInMeters
		}
	}

	minValue := math.Inf(1)
	var bestNeighborhood Neighborhood
	for _, node := range graph.nodes {
		nodeDistanceSum := neighborhoodDistanceSums[node.Name]
		if nodeDistanceSum < minValue {
			minValue = nodeDistanceSum
			bestNeighborhood = node
		}
	}

	return bestNeighborhood, nil
}
