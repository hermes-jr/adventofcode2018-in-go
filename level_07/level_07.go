package main

import (
	. "../utils"
	"fmt"
	"sort"
)

func main() {
	fname := "input" // ACBDESULXKYZIMNTFGWJVPOHRQ
	//	fname = "input_test" // CABDFE
	timelag := 60
	//	timelag = 0
	workers := 5
	//	workers = 2

	lines := ReadFile(fname)

	steps := make(map[int][]int)     // [StepID]:[Dependencies]
	stepsCopy := make(map[int][]int) // steps will be mutated, saving current state for part 2

	for _, inputLine := range lines {
		IfDebugPrintln(inputLine)
		// Step N must be finished before step J can begin.
		lb := int(inputLine[5])
		la := int(inputLine[36])

		IfDebugPrintf("%s -> needs -> %s\n", string(rune(la)), string(rune(lb)))
		if _, ok := steps[lb]; !ok {
			steps[lb] = []int{}
		}
		if _, ok := steps[la]; !ok {
			steps[la] = []int{}
		}
		steps[la] = append(steps[la], lb)
	}
	IfDebugPrintln("Steps read complete", steps)

	for k, v := range steps {
		tmp := make([]int, len(v))
		copy(tmp, steps[k])
		stepsCopy[k] = tmp
	}
	IfDebugPrintln("Steps copy for pt.2", stepsCopy)

	var result []byte
	for len(steps) > 0 {
		availableSteps := getNextAvailableStepsSorted(steps)
		IfDebugPrintln("Available steps", availableSteps)
		for _, nextStep := range availableSteps {
			result = append(result, byte(nextStep))
			IfDebugPrintln("Going to remove", nextStep, "from deps", steps)
			removeStepFromDependencies(steps, nextStep)
			delete(steps, nextStep)
			break
		}
	}
	fmt.Println("Result1", string(result))

	steps = stepsCopy

	deadlines := make(map[int]int) // [StepID]:endOfConstructionTime
	freeWorkers := workers

	elapsedTime := 0
	for ; len(steps) > 0; elapsedTime++ {
		for k, v := range deadlines {
			if elapsedTime == v {
				// remove current dependency from others
				IfDebugPrintln("Going to remove", k, "from deps", steps)
				removeStepFromDependencies(steps, k)
				delete(steps, k) // step done
				delete(deadlines, k)
				freeWorkers++
			}
		}

		if freeWorkers == 0 {
			continue // nobody free to take care of things at the moment
		}
		availableSteps := getNextAvailableStepsSortedDl(steps, deadlines)
		availableSteps = availableSteps[0:minOf(freeWorkers, len(availableSteps))] // assign to as many free workers as possible
		freeWorkers -= len(availableSteps)                                         // mark them busy

		IfDebugPrintf("time %v Available steps: %v, deadlines: %v\n", elapsedTime, availableSteps, deadlines)
		// place in execution queue
		for _, nextStep := range availableSteps {
			deadlines[nextStep] = elapsedTime + timelag + (nextStep - 64) // ascii A = 65 but A costs 1 second; Z = 90, costs 26
		}
	}
	fmt.Println("Result2", elapsedTime-1)
}

// Returns steps from pool with no required dependencies
func getNextAvailableStepsSorted(data map[int][]int) []int {
	return getNextAvailableStepsSortedDl(data, make(map[int]int))
}

// Returns steps from pool with no required dependencies
// and which are not currently being built.
func getNextAvailableStepsSortedDl(data map[int][]int, deadlines map[int]int) []int {
	var availableSteps []int
	for k, v := range data {
		if _, ok := deadlines[k]; !ok && len(v) == 0 {
			availableSteps = append(availableSteps, k)
		}
	}
	sort.Ints(availableSteps)
	return availableSteps
}

// Returns lesser of two int values
func minOf(n1 int, n2 int) int {
	if n1 <= n2 {
		return n1
	}
	return n2
}

// Iterates through data pool and removes given step from
// the list of dependencies of other steps
func removeStepFromDependencies(data map[int][]int, step int) {
	for zv, dataStepDependencies := range data {
		for i, dependency := range dataStepDependencies {
			if dependency == step {
				IfDebugPrintln(dataStepDependencies)
				data[zv] = append(dataStepDependencies[:i], dataStepDependencies[i+1:]...)
			}
		}
		IfDebugPrintln("After cleanup:", data)
	}
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

--- Part Two ---

As you're about to begin construction, four of the Elves offer to help. "The sun will set soon; it'll go faster if we work together." Now, you need to account for multiple people working on steps simultaneously. If multiple steps are available, workers should still begin them in alphabetical order.

Each step takes 60 seconds plus an amount corresponding to its letter: A=1, B=2, C=3, and so on. So, step A takes 60+1=61 seconds, while step Z takes 60+26=86 seconds. No time is required between steps.

To simplify things for the example, however, suppose you only have help from one Elf (a total of two workers) and that each step takes 60 fewer seconds (so that step A takes 1 second and step Z takes 26 seconds). Then, using the same instructions as above, this is how each second would be spent:

Second   Worker 1   Worker 2   Done
   0        C          .
   1        C          .
   2        C          .
   3        A          F       C
   4        B          F       CA
   5        B          F       CA
   6        D          F       CAB
   7        D          F       CAB
   8        D          F       CAB
   9        D          .       CABF
  10        E          .       CABFD
  11        E          .       CABFD
  12        E          .       CABFD
  13        E          .       CABFD
  14        E          .       CABFD
  15        .          .       CABFDE

Each row represents one second of time. The Second column identifies how many seconds have passed as of the beginning of that second. Each worker column shows the step that worker is currently doing (or . if they are idle). The Done column shows completed steps.

Note that the order of the steps has changed; this is because steps now take time to finish and multiple workers can begin multiple steps simultaneously.

In this example, it would take 15 seconds for two workers to complete these steps.

With 5 workers and the 60+ second step durations described above, how long will it take to complete all of the steps?

*/
