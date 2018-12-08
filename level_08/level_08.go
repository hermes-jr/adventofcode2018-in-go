package main

import (
	"container/list"
	"fmt"
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

	fmt.Println(len(data), data)

	//graph := graphs.NewDigraph()
	stack := list.New()
	readTree(data, 0, stack)
	//graph.Dump()
}

func readTree(data []int, idx int, stack *list.List) int {
	if idx > len(data)-2 {
		fmt.Println("End of array reached, idx", idx)
		return -1 // EOF //todo: do we add this value somewhere. might decrease counter and loop
	}
	nodeId := idx
	fmt.Printf("Depth %v Processing \"%v\", nodeId %v\n", stack.Len(), data[idx], nodeId)
	stack.PushBack(nodeId)
	//dumpIntStack(stack)
	unprocessedChildrenCount := data[idx]
	metaCount := data[idx+1]
	for uc := unprocessedChildrenCount; uc > 0; uc-- {
		// process remaining children
		fmt.Println("Next child is probably at", idx+2)
		possibleNext := readTree(data, idx+2 /* */, stack)
		if possibleNext == -1 {
			idx = nodeId + 2
		} else {
			idx = possibleNext + 2
		}
	}
	fmt.Println("No unprocessed children left")
	nodeMeta := data[idx+2 : idx+2+metaCount]
	fmt.Printf("Node %v meta found, %v items: %v\n", nodeId, metaCount, nodeMeta)
	stack.Remove(stack.Back()) // tree level processed
	//dumpIntStack(stack)
	return idx + metaCount
	// processChildrenRecursively
	// parentNote.metaStartIdx =
}

func dumpIntStack(stack *list.List) {
	fmt.Print("Stack state: <")
	for i := stack.Front(); i != nil; i = i.Next() {
		fmt.Print(i.Value)
		if i.Next() != nil {
			fmt.Print(", ")
		}
	}
	fmt.Println("<")
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
