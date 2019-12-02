package main

import (
	"math"
	"fmt"
)

// implementation of a max heap

type heap struct{
	nodes []*node
	nodeCount int
}

type node struct {
	value int
	parent *node
	leftChild *node
	rightChild *node
}


func (h *heap) insert(num int){
	//children at indices 2i + 1 and 2i + 2
	//its parent at index floor((i − 1) / 2)
	var newNode *node
	if len(h.nodes) == 0 {
		newNode = &node{num, nil, nil, nil}
	} else {
		newNode = &node{num ,h.nodes[int(math.Floor(float64(h.nodeCount - 1.0) / 2.0))], nil, nil}
		if len(h.nodes) % 2 == 0{
			newNode.parent.rightChild = newNode
		}else{
			newNode.parent.leftChild = newNode
		}
	}
	h.nodes = append(h.nodes, newNode)
	h.nodeCount += 1
	h.maxHeapify()
}

func (h *heap) isEmpty() bool {
	if len(h.nodes) == 0{
		return true
	}else{
		return false
	}
}

func (h *heap) heapSize() int{
	return h.nodeCount
}

func siftUp(firstNode *node) {
	switch {
	case firstNode.leftChild != nil && firstNode.rightChild != nil:
		if firstNode.leftChild.value > firstNode.rightChild.value {
			if firstNode.value < firstNode.leftChild.value {
				tempVal := firstNode.value
				firstNode.value = firstNode.leftChild.value
				firstNode.leftChild.value = tempVal
			}
		} else {
			if firstNode.value < firstNode.rightChild.value {
				tempVal := firstNode.value
				firstNode.value = firstNode.rightChild.value
				firstNode.rightChild.value = tempVal
			}
		}

	case firstNode.leftChild != nil:
		if firstNode.value < firstNode.leftChild.value {
			tempVal := firstNode.value
			firstNode.value = firstNode.leftChild.value
			firstNode.leftChild.value = tempVal
		}

	case firstNode.rightChild != nil:
		if firstNode.value < firstNode.value {
			tempVal := firstNode.rightChild.value
			firstNode.value = firstNode.rightChild.value
			firstNode.rightChild.value = tempVal
		}
	}
}

func (h *heap) maxHeapify(){
	i := len(h.nodes)/2
	for i >= 0 {
		cur := h.nodes[i]
		siftUp(cur)
		i--
	}
}

func (h *heap) printHeap(){
	for _, k := range h.nodes{
		fmt.Printf("%+v", *k)
		fmt.Printf("\n")
	}
}

func (h *heap) getMax() *node {
	return h.nodes[0]
}

func main(){
	heapush := heap{}
	
	heapush.insert(2)
	heapush.insert(7)
	heapush.insert(1)
	heapush.insert(36)
	heapush.insert(100)
	heapush.insert(19)
	heapush.insert(17)
	heapush.insert(3)
	heapush.insert(25)
	heapush.printHeap()
	fmt.Printf("The max value is %d\n", heapush.getMax().value)
}