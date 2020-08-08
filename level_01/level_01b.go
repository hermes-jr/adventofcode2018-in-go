package main

import (
	. "../utils"
	"container/ring"
	"fmt"
	"log"
	"strconv"
)

func main() {
	lines := ReadFile("input1")
	totalLines := len(lines)

	data := ring.New(totalLines)

	for _, line := range lines {
		fmt.Println("Reading value", line)
		z, err := strconv.ParseInt(line, 10, 64)
		fmt.Println("Adding to ring", z)
		data.Value = z
		data = data.Next()
		if err != nil {
			log.Fatal(err)
		}
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


 */
