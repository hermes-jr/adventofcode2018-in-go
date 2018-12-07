package main

import (
	"bufio"
	"fmt"
	"github.com/thcyron/graphs"
	"os"
	"sort"
)

//ACBEDSZULXKTIMNFGYWJVPOHRQ

/*
CDASZMULBEXKTINFGYWJVPOHRQ
CDASZULXEIWKNBTFMGYJVPOHRQ
EASCDZULXKBTINFMGYWJVPOHRQ
CDEASZULINXKBTFMGYWJVPOHRQ
ULAXCDKESZINBTFMGWYJVPOHRQ
ULCASDZXKBETINFYWJMGVPOHRQ
CBULDYAXKESZIWTNFJMGVPOHRQ
AULXCDSZEIWKBTNFYJMGVPOHRQ
AULXECDSZIWKBTNFMGYJVPOHRQ
ASCDZMEULXKBTINFGYWJVPOHRQ
ACBEDSZULXKTIMYWNFGJVPOHRQ
ASCDZULXKBETINFMGYWJVPOHRQ
ACBEDSZULXKTIMNFGYWJVPOHRQ
ACBEDSZULXKTINFMGYWJVPOHRQ
ECDASZULIXWNBKTYFJMGVPOHRQ
*/
func main() {
	fname := "input"
	//fname = "input_test"

	file, _ := os.Open(fname)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	graph := graphs.NewDigraph()

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
		graph.AddEdge(lb, la, float64(byte(la[0])))
	}

	fmt.Println("toCcount", toCount)
	fmt.Println("fromCount", fromCount)
	graph.Dump()

	var firstVertex []string
	var lastVertex string
	for k, v := range toCount {
		if v == 0 {
			firstVertex = append(firstVertex, k)
		}
		if fromCount[k] == 0 {
			lastVertex = k
		}
	}
	fmt.Printf("First and last instructions %v -> ... -> %v\n", firstVertex, lastVertex)

	var result1 []string
	result1 = printPart1(graph, lastVertex, fromCount, result1)
	fmt.Print("Result1: ")
	for i := len(result1) - 1; i >= 0; i-- {
		fmt.Print(result1[i])
	}
	fmt.Println()
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

func printPart1(graph *graphs.Graph, node string, tc map[string]int, res []string) []string {
	if tc[node] == 0 {
		res = append(res, node)
		fmt.Println("Cur res", res)
		var children []string
		for v := range graph.HalfedgesIter(node) {
			children = append(children, v.End.(string))
		}
		orderedChildren := sort.StringSlice(children)
		fmt.Printf("For node %v children %v | hist: %v\n", node, orderedChildren, tc)
		for i := len(orderedChildren) - 1; i >= 0; i-- {
			z := orderedChildren[i]
			tc[z]--
			if tc[z] >= 0 {
				res = printPart1(graph, z, tc, res)
			}
		}
	}
	return res
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
