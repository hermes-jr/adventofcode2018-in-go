package main

import (
	"container/list"
	"fmt"
	"github.com/thcyron/graphs"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type Tnode struct {
	id       int
	children int
	meta     []int
}

func main() {
	fname := "input"
	//fname = "input_test"

	b, _ := ioutil.ReadFile(fname)
	inTokens := strings.Split(string(b), " ")
	data := make([]int, 0, len(inTokens))
	for _, i := range inTokens {
		n, _ := strconv.Atoi(i)
		data = append(data, n)
	}

	fmt.Println(len(data), data)

	graph := graphs.NewDigraph()
	stack := list.New()
	knownNodes := make(map[int]Tnode)
	readTree(data, 0, stack, graph, knownNodes)
	graph.Dump()

	result1 := 0
	graphs.BFS(graph, 0, func(vertex graphs.Vertex, i *bool) {
		mt := knownNodes[vertex.(int)].meta
		fmt.Println("gonna sum", mt)
		for z := range mt {
			result1 += mt[z]
		}
	})
	fmt.Println("Result1", result1)

	result2 := heavyCount(graph, 0, knownNodes)
	fmt.Println("Result2", result2)
}

func heavyCount(graph *graphs.Graph, vertex int, knownNodes map[int]Tnode) int {
	selfWeight := 0
	thisNode := knownNodes[vertex]
	if thisNode.children == 0 {
		fmt.Printf("Nochild node found %v with meta %v\n", vertex, thisNode.meta)
		for z := range thisNode.meta {
			selfWeight += thisNode.meta[z]
		}
		fmt.Printf("Weight of nochild node %v is %v\n", vertex, selfWeight)
		return selfWeight
	}

	var children []int
	for outEdge := range graph.HalfedgesIter(vertex) {
		children = append(children, outEdge.End.(int))
	}
	sort.Ints(children) // Important
	fmt.Printf("Multichild node %v has children %v and meta %v\n", vertex, children, thisNode.meta)
	for i := range thisNode.meta {
		n := thisNode.meta[i] - 1
		if n >= 0 && n < len(children) {
			childNode := children[n]
			fmt.Println("Gonna calculate weight of", childNode)
			cwt := heavyCount(graph, childNode, knownNodes) // count weight of an Nth child node
			selfWeight += cwt
		}
	}

	fmt.Printf("Weight of multichild node %v is %v\n", vertex, selfWeight)
	return selfWeight
}

// Parses tree encoded in a sequence of integers.
// [Node0_children_count Node0_meta_count [Node1_children...etc] Node0_Meta<>]EOF
//
// Returns last index of a node that is currently being processed
func readTree(data []int, idx int, stack *list.List, graph *graphs.Graph, knownNodes map[int]Tnode) int {
	if idx > len(data)-2 {
		fmt.Println("End of array reached, idx", idx)
		return -1 // EOF
	}
	nodeId := idx
	fmt.Printf("Depth %v Processing \"%v\", nodeId %v\n", stack.Len(), data[idx], nodeId)
	stack.PushBack(nodeId)
	nodeChildrenCount := data[idx]
	metaCount := data[idx+1]
	for uc := nodeChildrenCount; uc > 0; uc-- {
		// process remaining children
		fmt.Println("Next child is probably at", idx+2)
		possibleNext := readTree(data, idx+2, stack, graph, knownNodes)
		if possibleNext != -1 {
			idx = possibleNext
		}
	}
	fmt.Println("No unprocessed children left")
	nodeMeta := data[idx+2 : idx+2+metaCount]
	fmt.Printf("Node %v meta found, %v items: %v\n", nodeId, metaCount, nodeMeta)
	knownNodes[nodeId] = Tnode{nodeId, nodeChildrenCount, nodeMeta}
	if stack.Len() > 1 {
		graph.AddEdge(stack.Back().Prev().Value.(int), stack.Back().Value.(int), float64(stack.Back().Value.(int)))
	}
	stack.Remove(stack.Back()) // tree level processed
	return idx + metaCount
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

--- Part Two ---

The second check is slightly more complicated: you need to find the value of the root node (A in the example above).

The value of a node depends on whether it has child nodes.

If a node has no child nodes, its value is the sum of its metadata entries. So, the value of node B is 10+11+12=33, and the value of node D is 99.

However, if a node does have child nodes, the metadata entries become indexes which refer to those child nodes. A metadata entry of 1 refers to the first child node, 2 to the second, 3 to the third, and so on. The value of this node is the sum of the values of the child nodes referenced by the metadata entries. If a referenced child node does not exist, that reference is skipped. A child node can be referenced multiple time and counts each time it is referenced. A metadata entry of 0 does not refer to any child node.

For example, again using the above nodes:

    Node C has one metadata entry, 2. Because node C has only one child node, 2 references a child node which does not exist, and so the value of node C is 0.
    Node A has three metadata entries: 1, 1, and 2. The 1 references node A's first child node, B, and the 2 references node A's second child node, C. Because node B has a value of 33 and node C has a value of 0, the value of node A is 33+33+0=66.

So, in this example, the value of the root node is 66.

What is the value of the root node?

*/
