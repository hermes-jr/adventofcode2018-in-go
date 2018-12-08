package main

import (
	"container/list"
	"fmt"
	"github.com/thcyron/graphs"
	"io/ioutil"
	"strconv"
	"strings"
)

type Tnode struct {
	id       int
	children int
	metasize int
	meta     []int
}

func main() {
	fname := "input"
	fname = "input_test"

	b, _ := ioutil.ReadFile(fname)
	inTokens := strings.Split(string(b), " ")
	data := make([]int, 0, len(inTokens))
	for _, i := range inTokens {
		n, _ := strconv.Atoi(i)
		data = append(data, n)
	}

	fmt.Println(data)

	graph := graphs.NewDigraph()
	stack := list.New()
	readTree(graph, data, 0, stack)
	graph.Dump()
}

func readTree(graph *graphs.Graph, data []int, ptr int, stack *list.List) {
	if ptr > len(data)-2 {
		return
	}
	fmt.Printf("Depth %v Processing \"%v\" at index %v\n", stack.Len(), data[ptr], ptr)
	childrenCount := data[ptr]
	metaCount := data[ptr+1]

	accumulatedOffset := 0
	for med := stack.Back(); med != nil; med = med.Prev() {
		accumulatedOffset += 2
		accumulatedOffset += med.Value.(Tnode).metasize
	}
	currentNodeMetaEnd := len(data) - accumulatedOffset
	currentNodeMetaStart := currentNodeMetaEnd - metaCount
	currentNodeMeta := data[currentNodeMetaStart:currentNodeMetaEnd]

	stack.PushBack(Tnode{ptr, childrenCount, metaCount, currentNodeMeta})
	dumpStack(stack)
	if childrenCount == 0 {
		// stub, read meta
		fmt.Println("Current node meta", currentNodeMeta)
		copy(stack.Back().Value.(Tnode).meta, currentNodeMeta)
		fmt.Println(stack.Back().Value.(Tnode))
		ptr += 2 + metaCount
		//fmt.Println(childrenCount, metaCount, currentNodeMeta)
	} else {
		ptr += 2
		//for ci := childrenCount; ci > 0; ci++ {
		//	readTree(graph, data, ptr+1, stack)
		//}
	}
	if stack.Back().Prev() != nil {
		graph.AddEdge(stack.Back().Prev(), stack.Back(), 0)
		stack.Remove(stack.Back())
	}
	readTree(graph, data, ptr, stack)

}

func dumpStack(stack *list.List) {
	fmt.Print("stack [\n")
	for i := stack.Front(); i != nil; i = i.Next() {
		iv := i.Value.(Tnode)
		fmt.Printf("\t{id: %3v children: %3v metadata: %3v:%v}", iv.id, iv.children, iv.metasize, iv.meta)
		if i.Next() != nil {
			fmt.Print(",\n")
		}
	}
	fmt.Println("\n]")
}

/*
--- Day 8: Memory Maneuver ---

The sleigh is much easier to pull than you'd expect for something its weight. Unfortunately, neither you nor the Elves know which way the North Pole is from here.

You check your wrist device for anything that might help. It seems to have some kind of navigation system! Activating the navigation system produces more bad news: "Failed to start navigation system. Could not read software license file."

The navigation system's license file consists of a list of numbers (your puzzle input). The numbers define a data structure which, when processed, produces some kind of tree that can be used to calculate the license number.

The tree is made up of nodes; a single, outermost node forms the tree's root, and it contains all other nodes in the tree (or contains nodes that contain nodes, and so on).

Specifically, a node consists of:

    A header, which is always exactly two numbers:
        The quantity of child nodes.
        The quantity of metadata entries.
    Zero or more child nodes (as specified in the header).
    One or more metadata entries (as specified in the header).

Each child node is itself a node that has its own header, child nodes, and metadata. For example:

2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2
A----------------------------------
    B----------- C-----------
                     D-----

In this example, each node of the tree is also marked with an underline starting with a letter for easier identification. In it, there are four nodes:

    A, which has 2 child nodes (B, C) and 3 metadata entries (1, 1, 2).
    B, which has 0 child nodes and 3 metadata entries (10, 11, 12).
    C, which has 1 child node (D) and 1 metadata entry (2).
    D, which has 0 child nodes and 1 metadata entry (99).

The first check done on the license file is to simply add up all of the metadata entries. In this example, that sum is 1+1+2+10+11+12+2+99=138.

What is the sum of all metadata entries?

*/
