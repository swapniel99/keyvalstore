package main

// This file is incomplete. Future work includes using heap for deleting expired values.
import (
//	"container/heap"
)

type node struct {
	expiry int64
	key    string
}

type nodeHeap []node

func (h nodeHeap) Len() int {
	return len(h)
}

func (h nodeHeap) Less(i, j int) bool {
	return (h[i].expiry < h[j].expiry)
}

func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
