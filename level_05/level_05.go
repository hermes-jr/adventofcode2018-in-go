package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	fname := "input"
	//fname = "input_test"

	// case offset = 32
	//  d   a  b  A  c  C  a  C  B  A  c  C  c  a  D  A
	// [100 97 98 65 99 67 97 67 66 65 99 67 99 97 68 65]

	data, _ := ioutil.ReadFile(fname)
	fmt.Println(data, "len", len(data))

	//outerloop:
	for lenChange := false; ; lenChange = false {
		lastIdx := len(data) - 1
		for i := 0; i < lastIdx; i++ {
			//fmt.Println("Processing", i, i+1)
			delta := int(data[i]) - int(data[i+1])
			if delta == 32 || delta == -32 {
				//fmt.Println("Opposite charged sequence found\n", data, "len", len(data), "\nremoving", i, ":", data[i], i+1, ":", data[i+1])
				data = append(data[:i], data[i+2:]...) // removing two elements in the middle
				lastIdx -= 2
				lenChange = true
			}
		}
		//fmt.Println("After purge cycle")
		fmt.Println(string(data), "len", len(data)) // len() returns actual len + 1 here, wtf?
		if !lenChange {
			fmt.Println("Result1", len(data))
			break
		}
	}

}

/*
--- Day 5: Alchemical Reduction ---

You've managed to sneak in to the prototype suit manufacturing lab. The Elves are making decent progress, but are still struggling with the suit's size reduction capabilities.

While the very latest in 1518 alchemical technology might have solved their problem eventually, you can do better. You scan the chemical composition of the suit's material and discover that it is formed by extremely long polymers (one of which is available as your puzzle input).

The polymer is formed by smaller units which, when triggered, react with each other such that two adjacent units of the same type and opposite polarity are destroyed. Units' types are represented by letters; units' polarity is represented by capitalization. For instance, r and R are units with the same type but opposite polarity, whereas r and s are entirely different types and do not react.

For example:

    In aA, a and A react, leaving nothing behind.
    In abBA, bB destroys itself, leaving aA. As above, this then destroys itself, leaving nothing.
    In abAB, no two adjacent units are of the same type, and so nothing happens.
    In aabAAB, even though aa and AA are of the same type, their polarities match, and so nothing happens.

Now, consider a larger example, dabAcCaCBAcCcaDA:

dabAcCaCBAcCcaDA  The first 'cC' is removed.
dabAaCBAcCcaDA    This creates 'Aa', which is removed.
dabCBAcCcaDA      Either 'cC' or 'Cc' are removed (the result is the same).
dabCBAcaDA        No further actions can be taken.

After all possible reactions, the resulting polymer contains 10 units.

How many units remain after fully reacting the polymer you scanned? (Note: in this puzzle and others, the input is large; if you copy/paste your input, make sure you get the whole thing.)

*/
