package api

import (
	"container/heap"
	"testing"
)

func TestGetMinHeapCorrectlySetsRootNode(t *testing.T) {
	frequencyMap := map[string]int{
		"Downtown":   1,
		"South Side": 5,
		"East End":   4,
		"Central":    3,
	}

	h := getMinHeap(frequencyMap)
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

func TestGetMinHeapIsEmptyWhenEmptyMapGiven(t *testing.T) {
	frequencyMap := map[string]int{}

	h := getMinHeap(frequencyMap)

	expectedHeapSize := 0
	if h.Len() != expectedHeapSize {
		t.Errorf("Heap size was incorrect. Got: %d, expected: %d.", h.Len(), expectedHeapSize)
	}
}

func TestMinHeapSwapsCorrectly(t *testing.T) {
	frequencyMap := map[string]int{}

	h := getMinHeap(frequencyMap)
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
	minHeap := getMinHeap(frequencyMap)

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

	minHeap := getMinHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(minHeap)

	expectedNeighborhoodsCount := 0
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}
}

func TestFindNeighborhoodsWithSameFrequency_oneHeapEntryGiven(t *testing.T) {
	expectedNeighborhoodName := "Downtown"
	frequencyMap := map[string]int{expectedNeighborhoodName: 1}

	minHeap := getMinHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(minHeap)

	expectedNeighborhoodsCount := 1
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}

	if neighborhoods[0] != expectedNeighborhoodName {
		t.Errorf("The returned neighborhood name was not correct. Got: %s, expected: %s.", neighborhoods[0], expectedNeighborhoodName)
	}
}
