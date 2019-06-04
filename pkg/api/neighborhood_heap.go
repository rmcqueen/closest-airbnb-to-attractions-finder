package api

import "container/heap"

type neighorboodNameFrequency struct {
	name  string
	count int
}

type neighborhoodNameFrequencyMaxHeap []neighorboodNameFrequency

func getMaxHeap(m map[string]int) *neighborhoodNameFrequencyMaxHeap {
	h := &neighborhoodNameFrequencyMaxHeap{}
	heap.Init(h)
	for k, v := range m {
		heap.Push(h, neighorboodNameFrequency{k, v})
	}

	return h
}

func (h neighborhoodNameFrequencyMaxHeap) Less(i, j int) bool { return h[i].count > h[j].count }
func (h neighborhoodNameFrequencyMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h neighborhoodNameFrequencyMaxHeap) Len() int           { return len(h) }

func (h *neighborhoodNameFrequencyMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(neighorboodNameFrequency))
}

func (h *neighborhoodNameFrequencyMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
