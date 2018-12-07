package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	fname := "input"
	//fname = "input_test"

	file, _ := os.Open(fname)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	unhd := make(map[int][]int)

	for scanner.Scan() {
		inputLine := scanner.Text()
		fmt.Println(inputLine)
		// Step N must be finished before step J can begin.
		lb := int(inputLine[5])
		la := int(inputLine[36])
		fmt.Printf("%v -> needs -> %v\n", string(la), string(lb))
		if _, ok := unhd[lb]; !ok {
			unhd[lb] = []int{}
		}
		if _, ok := unhd[la]; !ok {
			unhd[la] = []int{}
		}
		unhd[la] = append(unhd[la], lb)
	}
	fmt.Println(unhd)

	var result []byte
	//outerloop:
	for {
		if len(unhd) == 0 {
			// Available step candidates depleted
			break
		}
		var availableSteps []int
		fmt.Print("Next available steps unsorted:")
		for k, v := range unhd {
			if len(v) == 0 {
				availableSteps = append(availableSteps, k)
				fmt.Print(string(byte(k)))
			}
		} // 65 - A, 69 - E
		fmt.Println()
		sort.Ints(availableSteps)
		if len(availableSteps) == 0 {
			break
		}
		fmt.Println("Available steps", availableSteps)
		for _, nextStep := range availableSteps {
			result = append(result, byte(nextStep))
			// remove current dependency from others
			fmt.Println("Going to remove", nextStep, "from deps", unhd)
			for zv, survivorStepDeps := range unhd {
				for i, sv := range survivorStepDeps {
					fmt.Println("comparing", sv, nextStep)
					if sv == nextStep {
						fmt.Println(survivorStepDeps)
						unhd[zv] = append(survivorStepDeps[:i], survivorStepDeps[i+1:]...)
					}
				}
				fmt.Println("After cleanup:", unhd)
			}
			//continue outerloop
			delete(unhd, nextStep)
			break
		}
	}
	fmt.Println("Result1", string(result))
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
