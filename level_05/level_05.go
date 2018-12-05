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
	cloneOfOriginal := make([]byte, len(data))
	copy(cloneOfOriginal, data)

	result1 := foldPolymer(data)
	fmt.Println("Result1", result1)

	smallestPolymerSeen := len(cloneOfOriginal)
	for letter := byte(65); letter <= byte(90); letter++ {
		withoutSomeAtom := make([]byte, len(cloneOfOriginal))
		copy(withoutSomeAtom, cloneOfOriginal)
		//fmt.Println("Testing", letter, string(byte(letter)))

		changed := false
		lastIdx := len(withoutSomeAtom) - 1
		for i := 0; i <= lastIdx; i++ {
			//fmt.Println("Comparing", i, withoutSomeAtom[i], "to", letter, letter+32)
			if withoutSomeAtom[i] == letter || withoutSomeAtom[i] == letter+32 {
				//fmt.Println("Removing", string(letter), string(letter+32), "at", i, string(withoutSomeAtom))
				withoutSomeAtom = append(withoutSomeAtom[:i], withoutSomeAtom[i+1:]...)
				lastIdx--
				i--
				//fmt.Println("Removed", string(letter), string(letter+32), "at", i, "got", string(withoutSomeAtom))
				changed = true // letter found
			}
		}
		if changed {
			afterFold := foldPolymer(withoutSomeAtom)
			//fmt.Println("With no", string(letter), "size is", afterFold)
			if afterFold < smallestPolymerSeen {
				smallestPolymerSeen = afterFold
			}
		}
	}
	fmt.Println("Result2", smallestPolymerSeen)
}

func foldPolymer(polymer []byte) int {
	for lenChange := false; ; lenChange = false {
		lastIdx := len(polymer) - 1
		for i := 0; i < lastIdx; i++ {
			//fmt.Println("Processing", i, i+1)
			delta := int(polymer[i]) - int(polymer[i+1])
			if delta == 32 || delta == -32 {
				//fmt.Println("Opposite charged sequence found\n", data, "len", len(data), "\nremoving", i, ":", data[i], i+1, ":", data[i+1])
				polymer = append(polymer[:i], polymer[i+2:]...)
				lastIdx -= 2
				lenChange = true
			}
		}
		//fmt.Println("After purge cycle")
		if !lenChange {
			return len(polymer)
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

--- Part Two ---

Time to improve the polymer.

One of the unit types is causing problems; it's preventing the polymer from collapsing as much as it should. Your goal is to figure out which unit type is causing the most problems, remove all instances of it (regardless of polarity), fully react the remaining polymer, and measure its length.

For example, again using the polymer dabAcCaCBAcCcaDA from above:

    Removing all A/a units produces dbcCCBcCcD. Fully reacting this polymer produces dbCBcD, which has length 6.
    Removing all B/b units produces daAcCaCAcCcaDA. Fully reacting this polymer produces daCAcaDA, which has length 8.
    Removing all C/c units produces dabAaBAaDA. Fully reacting this polymer produces daDA, which has length 4.
    Removing all D/d units produces abAcCaCBAcCcaA. Fully reacting this polymer produces abCBAc, which has length 6.

In this example, removing all C/c units was best, producing the answer 4.

What is the length of the shortest polymer you can produce by removing all units of exactly one type and fully reacting the result?

*/
