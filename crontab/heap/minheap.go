package heap

import (
	"errors"
	"fmt"
)

var (
	ErrFull = errors.New("Full min heap")
)

// 遍历heap每个结点，对其做WalkFunc操作
type WalkFunc func(e Element) Element

type Element interface {
	// Key get the element's key
	Key() interface{}

	// Value get the element's value
	Value() interface{}

	// Compare compare by element's key
	Compare(key interface{}) int

	// Merge join together val that have the save key
	//Merge(val interface{})
}

// MinHeap is sorted by Element's Key()
type MinHeap struct {
	arr []Element

	maxSize int
}

func NewMinHeap(maxSize int) *MinHeap {
	return &MinHeap{
		arr:     make([]Element, 0),
		maxSize: maxSize,
	}
}

func (h *MinHeap) Push(e Element) error {
	if e == nil {
		return errors.New("Nil Element")
	}

	// append to the tail of heap
	if len(h.arr) >= h.maxSize {
		return ErrFull
	}
	h.arr = append(h.arr, e)

	// fix heap
	var curIdx int = len(h.arr) - 1
	var parentIdx int

	for curIdx > 0 {
		if curIdx^(curIdx-1) == 0 {
			// odd number
			parentIdx = (curIdx - 2) / 2
		} else {
			// even number
			parentIdx = (curIdx - 1) / 2
		}

		if parentIdx < 0 {
			return fmt.Errorf("Node[%d]'s parent node[%d]'s idx < 0", curIdx, parentIdx)
		}

		if h.arr[curIdx].Compare(h.arr[parentIdx].Key()) >= 0 {
			break
		}

		// cur node less then it's parent, exchange them,and then step to it's parent's idx
		h.arr[curIdx], h.arr[parentIdx] = h.arr[parentIdx], h.arr[curIdx]
		curIdx = parentIdx
	}
	return nil
}

// If empty, nil will be returned
func (h *MinHeap) Pop() Element {
	if len(h.arr) <= 0 {
		return nil
	}

	// get the top one element
	ret := h.arr[0]

	// move the last one element to the top
	h.arr[0] = h.arr[len(h.arr)-1]
	h.arr = h.arr[0 : len(h.arr)-1]

	// fix heap
	if len(h.arr) > 0 {
		var (
			curIdx  = 0
			lastIdx = len(h.arr) - 1
		)

		for {
			var (
				minIdx        = curIdx
				leftChildIdx  = 2*curIdx + 1
				rightChildIdx = 2*curIdx + 2
			)

			if leftChildIdx > lastIdx && rightChildIdx > lastIdx {
				// No child node, break
				break
			}
			if leftChildIdx <= lastIdx && h.arr[leftChildIdx].Compare(h.arr[minIdx].Key()) < 0 {
				minIdx = leftChildIdx
			}
			if rightChildIdx <= lastIdx && h.arr[rightChildIdx].Compare(h.arr[minIdx].Key()) < 0 {
				minIdx = rightChildIdx
			}

			if minIdx == curIdx {
				// parent node is the min node, so break
				break
			}

			// exchage curIdx with minIdx, and then loop
			h.arr[curIdx], h.arr[minIdx] = h.arr[minIdx], h.arr[curIdx]
			curIdx = minIdx
		}
	}

	return ret
}

// Top return the smallest element in the heap
func (h *MinHeap) Top() Element {
	if len(h.arr) <= 0 {
		return nil
	}
	return h.arr[0]
}

// Walk apply function f to all elements
func (h *MinHeap) Walk(f WalkFunc) {
	for i, e := range h.arr {
		if newE := f(e); newE != nil {
			h.arr[i] = newE
		}
	}
}

func (h *MinHeap) Debug() {
	for {
		o := h.Pop()
		if o == nil {
			break
		}
		fmt.Printf("key=%v\n", o.Key())
		fmt.Printf("val=%v\n", o.Value())
		fmt.Printf("\n")
	}
}
