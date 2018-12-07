package main

import (
	"bufio"
	"fmt"
	"github.com/thcyron/graphs"
	"os"
)

func main() {
	fname := "input"
	fname = "input_test"

	file, _ := os.Open(fname)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	graph := graphs.NewGraph()

	toCount := make(map[string]int)
	fromCount := make(map[string]int)
	for scanner.Scan() {
		inputLine := scanner.Text()
		fmt.Println(inputLine)
		la := inputLine[5:6]
		lb := inputLine[36:37]
		fmt.Printf("%v -> %v\n", la, lb)

		countPaths(toCount, fromCount, lb, la)
		toCount[lb]++
		fromCount[la]++
		graph.AddEdge(lb, la, -1.0*float64(byte(la[0])))
	}

	fmt.Println("toCcount", toCount)
	fmt.Println("fromCount", fromCount)
	graph.Dump()

	var firstVertex, lastVertex string
	for k, v := range toCount {
		if v == 0 {
			firstVertex = k
		}
		if fromCount[k] == 0 {
			lastVertex = k
		}
	}
	fmt.Printf("First and last instructions %v -> ... -> %v\n", firstVertex, lastVertex)

	printPart1(graph, lastVertex)

}

func countPaths(toCount map[string]int, fromCount map[string]int, lb string, la string) {
	if _, ok := toCount[lb]; !ok {
		toCount[lb] = 0
	}
	if _, ok := toCount[la]; !ok {
		toCount[la] = 0
	}
	if _, ok := fromCount[lb]; !ok {
		fromCount[lb] = 0
	}
	if _, ok := fromCount[la]; !ok {
		fromCount[la] = 0
	}
}

func printPart1(graph *graphs.Graph, node string) {
	satisfiedVertices := make(map[string]bool)
	//satisfiedVertices[node] = true

	graphs.BFS(graph, node, func(vertex graphs.Vertex, i *bool) {
		if satisfiedVertices[vertex.(string)] {
			return
		}
		satisfiedVertices[vertex.(string)] = true
		fmt.Println("V: ", vertex)
		/*		if satisfiedVertices[vertex.(string)] {
					result = append(result, vertex.(string))
				}
				// satisfiedVertices[vertex.(string)] = true
				var children []string
				for v := range graph.HalfedgesIter(vertex) {
					children = append(children, v.End.(string))
				}
				orderedChildren := sort.StringSlice(children)
				rr := append([]string, orderedChildren, result)
				result = rr
				fmt.Printf("\nIter\n\tvertex: %v\n\tcurrentlySatisfied: %v\n\tcurrentResult: %v\n\torderedChildren: %v\n", vertex, satisfiedVertices, result, orderedChildren)
		*/
	})
}

/*
--- Day 7: The Sum of Its Parts ---

You find yourself standing on a snow-covered coastline; apparently, you landed a little off course. The region is too hilly to see the North Pole from here, but you do spot some Elves that seem to be trying to unpack something that washed ashore. It's quite cold out, so you decide to risk creating a paradox by asking them for directions.

"Oh, are you the search party?" Somehow, you can understand whatever Elves from the year 1018 speak; you assume it's Ancient Nordic Elvish. Could the device on your wrist also be a translator? "Those clothes don't look very warm; take this." They hand you a heavy coat.

"We do need to find our way back to the North Pole, but we have higher priorities at the moment. You see, believe it or not, this box contains something that will solve all of Santa's transportation problems - at least, that's what it looks like from the pictures in the instructions." It doesn't seem like they can read whatever language it's in, but you can: "Sleigh kit. Some assembly required."

"'Sleigh'? What a wonderful name! You must help us assemble this 'sleigh' at once!" They start excitedly pulling more parts out of the box.

The instructions specify a series of steps and requirements about which steps must be finished before others can begin (your puzzle input). Each step is designated by a single letter. For example, suppose you have the following instructions:

Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.

Visually, these requirements look like this:


  -->A--->B--
 /    \      \
C      -->D----->E
 \           /
  ---->F-----

Your first goal is to determine the order in which the steps should be completed. If more than one step is ready, choose the step which is first alphabetically. In this example, the steps would be completed as follows:

    Only C is available, and so it is done first.
    Next, both A and F are available. A is first alphabetically, so it is done next.
    Then, even though F was available earlier, steps B and D are now also available, and B is the first alphabetically of the three.
    After that, only D and F are available. E is not available because only some of its prerequisites are complete. Therefore, D is completed next.
    F is the only choice, so it is done next.
    Finally, E is completed.

So, in this example, the correct order is CABDFE.

In what order should the steps in your instructions be completed?

*/
