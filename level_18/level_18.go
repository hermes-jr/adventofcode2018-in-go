package main

import (
	"bufio"
	"fmt"
	"os"
)

const DEBUG = false

type forestMap [][]byte

func main() {
	fname := "input"
	//fname = "input_test1"

	file, _ := os.Open(fname)
	defer file.Close()
	var data forestMap

	scanner := bufio.NewScanner(file)

	for rowNum := 0; scanner.Scan(); rowNum++ {
		data = append(data, []byte{})
		inputLine := scanner.Bytes()
		for _, v := range inputLine {
			data[rowNum] = append(data[rowNum], v)
		}
	}

	if DEBUG {
		printMap(&data)
	}

	for step := 0; step < 1000000000; step++ {
		nextStepData := make([][]byte, len(data))
		tw := 0
		tl := 0
		for y := range data {
			nextStepData[y] = make([]byte, len(data[y]))
			for x, v := range data[y] {
				_, t, l := calcAdjacent(&data, x, y)
				//    An open acre will become filled with trees if three or more adjacent acres contained trees. Otherwise, nothing happens.
				//    An acre filled with trees will become a lumberyard if three or more adjacent acres were lumberyards. Otherwise, nothing happens.
				//    An acre containing a lumberyard will remain a lumberyard if it was adjacent to at least one other lumberyard and at least one acre containing trees. Otherwise, it becomes open.
				switch v {
				case '.':
					if t >= 3 {
						nextStepData[y][x] = '|'
						tw++
					} else {
						nextStepData[y][x] = '.'
					}
				case '|':
					if l >= 3 {
						nextStepData[y][x] = '#'
						tl++
					} else {
						nextStepData[y][x] = '|'
						tw++
					}
				case '#':
					if l >= 1 && t >= 1 {
						nextStepData[y][x] = '#'
						tl++
					} else {
						nextStepData[y][x] = '.'
					}
				}
			}
		}
		data = nextStepData
		if DEBUG {
			printMap(&data)
		}
		if step == 9 {
			fmt.Println("Result1:", tw*tl)

		}
	}
}

// Calculates number of each type cells among adjacent
func calcAdjacent(data *forestMap, x, y int) (int, int, int) {
	e, t, l := 0, 0, 0
	deref := *data

	for i := y - 1; i < y+2; i++ {
		for j := x - 1; j < x+2; j++ {
			if i < 0 || i >= len(deref) || j < 0 || j >= len(deref[0]) || (i == y && j == x) {
				continue
			}
			v := deref[i][j]
			switch v {
			case '.':
				e++
			case '|':
				t++
			case '#':
				l++
			}
		}
	}
	return e, t, l
}

// Prints forest map
func printMap(data *forestMap) {
	deref := *data
	if len(deref) == 0 {
		fmt.Println("-EMPTY MAP-")
		return
	}

	for i := 0; i < len(deref); i++ {
		for j := 0; j < len(deref[i]); j++ {
			fmt.Print(string(deref[i][j]))
		}
		fmt.Println()
	}
}

// Returns lesser of two int values
func minOf(n1 int, n2 int) int {
	if n1 <= n2 {
		return n1
	}
	return n2
}

// Returns larger of two int values
func maxOf(n1 int, n2 int) int {
	if n1 >= n2 {
		return n1
	}
	return n2
}

/*
--- Day 18: Settlers of The North Pole ---

On the outskirts of the North Pole base construction project, many Elves are collecting lumber.

The lumber collection area is 50 acres by 50 acres; each acre can be either open ground (.), trees (|), or a lumberyard (#). You take a scan of the area (your puzzle input).

Strange magic is at work here: each minute, the landscape looks entirely different. In exactly one minute, an open acre can fill with trees, a wooded acre can be converted to a lumberyard, or a lumberyard can be cleared to open ground (the lumber having been sent to other projects).

The change to each acre is based entirely on the contents of that acre as well as the number of open, wooded, or lumberyard acres adjacent to it at the start of each minute. Here, "adjacent" means any of the eight acres surrounding that acre. (Acres on the edges of the lumber collection area might have fewer than eight adjacent acres; the missing acres aren't counted.)

In particular:

    An open acre will become filled with trees if three or more adjacent acres contained trees. Otherwise, nothing happens.
    An acre filled with trees will become a lumberyard if three or more adjacent acres were lumberyards. Otherwise, nothing happens.
    An acre containing a lumberyard will remain a lumberyard if it was adjacent to at least one other lumberyard and at least one acre containing trees. Otherwise, it becomes open.

These changes happen across all acres simultaneously, each of them using the state of all acres at the beginning of the minute and changing to their new form by the end of that same minute. Changes that happen during the minute don't affect each other.

For example, suppose the lumber collection area is instead only 10 by 10 acres with this initial configuration:

Initial state:
.#.#...|#.
.....#|##|
.|..|...#.
..|#.....#
#.#|||#|#|
...#.||...
.|....|...
||...#|.#|
|.||||..|.
...#.|..|.

After 1 minute:
.......##.
......|###
.|..|...#.
..|#||...#
..##||.|#|
...#||||..
||...|||..
|||||.||.|
||||||||||
....||..|.

After 2 minutes:
.......#..
......|#..
.|.|||....
..##|||..#
..###|||#|
...#|||||.
|||||||||.
||||||||||
||||||||||
.|||||||||

After 3 minutes:
.......#..
....|||#..
.|.||||...
..###|||.#
...##|||#|
.||##|||||
||||||||||
||||||||||
||||||||||
||||||||||

After 4 minutes:
.....|.#..
...||||#..
.|.#||||..
..###||||#
...###||#|
|||##|||||
||||||||||
||||||||||
||||||||||
||||||||||

After 5 minutes:
....|||#..
...||||#..
.|.##||||.
..####|||#
.|.###||#|
|||###||||
||||||||||
||||||||||
||||||||||
||||||||||

After 6 minutes:
...||||#..
...||||#..
.|.###|||.
..#.##|||#
|||#.##|#|
|||###||||
||||#|||||
||||||||||
||||||||||
||||||||||

After 7 minutes:
...||||#..
..||#|##..
.|.####||.
||#..##||#
||##.##|#|
|||####|||
|||###||||
||||||||||
||||||||||
||||||||||

After 8 minutes:
..||||##..
..|#####..
|||#####|.
||#...##|#
||##..###|
||##.###||
|||####|||
||||#|||||
||||||||||
||||||||||

After 9 minutes:
..||###...
.||#####..
||##...##.
||#....###
|##....##|
||##..###|
||######||
|||###||||
||||||||||
||||||||||

After 10 minutes:
.||##.....
||###.....
||##......
|##.....##
|##.....##
|##....##|
||##.####|
||#####|||
||||#|||||
||||||||||

After 10 minutes, there are 37 wooded acres and 31 lumberyards. Multiplying the number of wooded acres by the number of lumberyards gives the total resource value after ten minutes: 37 * 31 = 1147.

What will the total resource value of the lumber collection area be after 10 minutes?

--- Part Two ---

This important natural resource will need to last for at least thousands of years. Are the Elves collecting this lumber sustainably?

What will the total resource value of the lumber collection area be after 1000000000 minutes?

*/
