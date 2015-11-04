package me

import(
	"container/heap"
)


// An MinHeap is a min-heap of ints.
type MinHeap []struct{int; string}

var heaps map[string]*MinHeap

func init() {
	heaps = make(map[string]*MinHeap)
}

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].int < h[j].int }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h MinHeap) Get(i int) (int, string) {
	if i < h.Len() {
		return h[i].int, h[i].string
	}
	
	return -1, ""
}
func (h MinHeap) Min() int {
	if ( h.Len() > 0 ) {
		return h[0].int
	}
	
	return -1
}

func (h *MinHeap) PushVal(ts int, id string) {
	heap.Push(h, struct{int; string}{ts, id})
}

func (h *MinHeap) PopVal() (int, string) {
	ev := heap.Pop(h).(struct{int; string})
	return ev.int, ev.string
}

func NewHeap(name string) *MinHeap {
	h := &MinHeap{}
	heap.Init(h)
	heaps[name] = h
	return h
}

func (h *MinHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(struct{int; string}))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func GetHeap(name string) (h *MinHeap, err bool) {
	h, err = heaps[name]	
	return h, err
	
}
