package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input1")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	totalLines := 0
	for scanner.Scan() {
		totalLines++
		fmt.Println("Counting lines", totalLines)
	}

	data := ring.New(totalLines)
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		unparsed := scanner.Text()
		fmt.Println("Reading value", unparsed)
		z, err := strconv.ParseInt(unparsed, 10, 64)
		fmt.Println("Adding to ring", z)
		data.Value = z
		data = data.Next()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Printing data")

	var lastfreq int64 = 0
	seenfreqs := map[int64]bool{}
	for {
		fmt.Println("Iteration through ring", data.Value)
		lastfreq += data.Value.(int64)
		if seenfreqs[lastfreq] {
			fmt.Println("Result2: ", lastfreq)
			break
		}
		seenfreqs[lastfreq] = true
		data = data.Next()
	}
}

/*

--- Part Two ---

You notice that the device repeats the same frequency change list over and over. To calibrate the device, you need to find the first frequency it reaches twice.

For example, using the same list of changes above, the device would loop as follows:

    Current frequency  0, change of +1; resulting frequency  1.
    Current frequency  1, change of -2; resulting frequency -1.
    Current frequency -1, change of +3; resulting frequency  2.
    Current frequency  2, change of +1; resulting frequency  3.
    (At this point, the device continues from the start of the list.)
    Current frequency  3, change of +1; resulting frequency  4.
    Current frequency  4, change of -2; resulting frequency  2, which has already been seen.

In this example, the first frequency reached twice is 2. Note that your device might need to repeat its list of frequency changes many times before a duplicate frequency is found, and that duplicates might be found while in the middle of processing the list.

Here are other examples:

    +1, -1 first reaches 0 twice.
    +3, +3, +4, -2, -4 first reaches 10 twice.
    -6, +3, +8, +5, -6 first reaches 5 twice.
    +7, +7, -2, -7, -4 first reaches 14 twice.

What is the first frequency your device reaches twice?

*/
