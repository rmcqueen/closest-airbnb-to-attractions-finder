package api

import (
	"container/heap"
	"testing"
)

func TestMaxHeap_popReturnsLargestElement(t *testing.T) {
	frequencyMap := map[string]int{
		"Downtown":   1,
		"South Side": 5,
		"East End":   4,
		"Central":    3,
	}

	h := getMaxHeap(frequencyMap)

	rootNode := heap.Pop(h)

	expectedRootNodeValue := 5
	rootNodeValue := rootNode.(neighorboodNameFrequency).count
	if rootNodeValue != expectedRootNodeValue {
		t.Errorf("The root node's value was incorrect. Expected: %d, got: %d.", expectedRootNodeValue, rootNodeValue)
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

	rootNode := heap.Pop(h).(neighorboodNameFrequency)
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

	rootNode := heap.Pop(h).(neighorboodNameFrequency)
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

	maxHeap := getMaxHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(maxHeap)

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

	maxHeap := getMaxHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(maxHeap)

	expectedNeighborhoodsCount := 0
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}
}

func TestFindNeighborhoodsWithSameFrequency_oneHeapEntryGiven(t *testing.T) {
	expectedNeighborhoodName := "Downtown"
	frequencyMap := map[string]int{expectedNeighborhoodName: 1}

	maxHeap := getMaxHeap(frequencyMap)

	neighborhoods, _ := findNeighborhoodsWithSameFrequency(maxHeap)

	expectedNeighborhoodsCount := 1
	if len(neighborhoods) != expectedNeighborhoodsCount {
		t.Errorf("Number of neighborhoods was invalid. Got: %d, expected: %d.", len(neighborhoods), expectedNeighborhoodsCount)
	}

	if neighborhoods[0] != expectedNeighborhoodName {
		t.Errorf("The returned neighborhood name was not correct. Got: %s, expected: %s.", neighborhoods[0], expectedNeighborhoodName)
	}
}
