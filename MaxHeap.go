package main

import (
	"math"
	"fmt"
	"os"
	"bufio"
	"strconv"
)

var nPorts int
var nMessages int

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

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
	nMessages ++
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

func readFileParameters(filepath string) {
	file, err := os.Open(filepath)
	CheckError(err)

	defer file.Close()

	reader := bufio.NewReader(file)

	// reading number of ports
	line, _, err := reader.ReadLine()
	CheckError(err)
	nPorts, err = strconv.Atoi(string(line))
}

func (h *heap) printHeapInFile() {
	f, err := os.OpenFile("heap.txt", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	msgToPrint := ""

	for i := 0; i < nPorts; i++ {
		id := h.nodes[i].value
		parent := h.nodes[i].parent
		leftChild := h.nodes[i].leftChild
		rightChild := h.nodes[i].rightChild

		msgToPrint = strconv.Itoa(id) + " "

		if parent != nil {
			msgToPrint += strconv.Itoa(parent.value) + " "
		} else {
			msgToPrint += strconv.Itoa(-1) + " "
		}

		if leftChild != nil {
			msgToPrint += strconv.Itoa(leftChild.value) + " "
		} else {
			msgToPrint += strconv.Itoa(-1) + " "
		}

		if rightChild != nil {
			msgToPrint += strconv.Itoa(rightChild.value) + "\n"
		} else {
			msgToPrint += strconv.Itoa(-1) + "\n"
		}
		fmt.Printf("%s", msgToPrint)

		if _, err := f.WriteString(msgToPrint); err != nil {
			fmt.Println(err)
		}
	}
}

func main(){
	heapush := heap{}
	readFileParameters("params.txt")

	for i := 1; i < nPorts+1 ; i++ {
		nMessages = 0
		heapush.insert(i)
	}

	// heapush.printHeap()
	// fmt.Printf("The max value is %d\n", heapush.getMax().value)
	
	heapush.printHeapInFile()

	// print("nMessages = ", 2*(nMessages-1), "\n")
	print("nMessages with brother = ", 2*(nMessages), "\n")
}