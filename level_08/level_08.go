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
	meta     *[]int
}
type Tnodes []Tnode

func (p Tnodes) Len() int           { return len(p) }
func (p Tnodes) Less(i, j int) bool { return p[i].id < p[j].id }
func (p Tnodes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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
	root, _ := readTree(data, 0, stack, graph)
	graph.Dump()

	result1 := 0
	graphs.BFS(graph, root, func(v graphs.Vertex, i *bool) {
		mt := *(v.(Tnode).meta)
		for z := range mt {
			result1 += mt[z]
		}
	})
	fmt.Println("Result1", result1)

	result2 := heavyCount(graph, root)
	fmt.Println("Result2", result2)
}

// Counts node value by the rules of part 2
func heavyCount(graph *graphs.Graph, vertex graphs.Vertex) int {
	selfValue := 0
	thisNode := vertex.(Tnode)
	if thisNode.children == 0 {
		fmt.Printf("Nochild node found %v with meta %v\n", thisNode, thisNode.meta)
		for z := range *thisNode.meta {
			selfValue += (*thisNode.meta)[z]
		}
		fmt.Printf("Value of nochild node %v is %v\n", thisNode, selfValue)
		return selfValue
	}

	var children Tnodes
	for outEdge := range graph.HalfedgesIter(vertex) {
		children = append(children, outEdge.End.(Tnode))
	}
	sort.Sort(children) // Important
	fmt.Printf("Multichild node %v has children %v and meta %v\n", thisNode, children, thisNode.meta)
	for i := range *thisNode.meta {
		n := (*thisNode.meta)[i] - 1
		if n >= 0 && n < len(children) {
			childNode := children[n]
			fmt.Println("Gonna calculate value of", childNode)
			cwt := heavyCount(graph, childNode) // count value of an Nth child node
			selfValue += cwt
		}
	}

	fmt.Printf("Value of multichild node %v is %v\n", thisNode, selfValue)
	return selfValue
}

// Parses tree encoded in a sequence of integers.
// [Node0_children_count Node0_meta_count [Node1_children...etc] Node0_Meta<>]EOF
//
// Returns tree root (eventually) and last index of a node that is currently being processed
func readTree(data []int, dataIdx int, stack *list.List, graph *graphs.Graph) (graphs.Vertex, int) {
	if dataIdx > len(data)-2 {
		fmt.Println("End of array reached, dataIdx", dataIdx)
		return nil, -1 // EOF
	}
	nodeId := dataIdx
	fmt.Printf("Depth %v Processing \"%v\", nodeId %v\n", stack.Len(), data[dataIdx], nodeId)
	nodeChildrenCount := data[dataIdx]
	currentVertex := Tnode{nodeId, nodeChildrenCount, &[]int{}}
	stack.PushBack(currentVertex)
	metaCount := data[dataIdx+1]
	for uc := nodeChildrenCount; uc > 0; uc-- {
		// process remaining children
		fmt.Println("Next child is probably at", dataIdx+2)
		_, possibleNext := readTree(data, dataIdx+2, stack, graph)
		if possibleNext != -1 {
			dataIdx = possibleNext
		}
	}
	fmt.Println("No unprocessed children left")
	nodeMeta := data[dataIdx+2 : dataIdx+2+metaCount]
	fmt.Printf("Node %v meta found, %v items: %v\n", nodeId, metaCount, nodeMeta)
	*currentVertex.meta = nodeMeta

	// if there's a parent down the stack, create a link
	if stack.Len() > 1 {
		graph.AddEdge(stack.Back().Prev().Value, currentVertex, float64(stack.Back().Value.(Tnode).id))
	}
	stack.Remove(stack.Back()) // tree level processed
	return currentVertex, dataIdx + metaCount
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
