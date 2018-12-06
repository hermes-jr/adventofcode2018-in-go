package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	x    int
	y    int
	name byte
}

func (data Point) String() string {
	return fmt.Sprintf("{%v %v:%v}", string(data.name), data.x, data.y)
}

type Points []Point
type Field [][]uint8

func main() {
	fname := "input"
	fname = "input_test"

	file, _ := os.Open(fname)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inData := Points{}

	minX, minY, maxX, maxY := -1, -1, 0, 0

	// scan input, find out "populated" area bounds
	for pointNameCode := byte(65); scanner.Scan(); pointNameCode++ {
		inputLine := scanner.Text()
		fmt.Println(inputLine)
		re := regexp.MustCompile("^(\\d+), (\\d+)$")
		match := re.FindStringSubmatch(inputLine)
		y, _ := strconv.Atoi(match[1])
		x, _ := strconv.Atoi(match[2])
		inData = append(inData, Point{x, y, byte(pointNameCode)})
		if maxX < x {
			maxX = x
		}
		if maxY < y {
			maxY = y
		}
		if minX < 0 || x < minX {
			minX = x
		}
		if minY < 0 || y < minY {
			minY = y
		}
		if pointNameCode == 90 {
			pointNameCode = 96 // jump from Z to a
		}
	}
	fmt.Println(inData)

	fh := maxY - minY + 1
	fw := maxX - minX + 1
	fmt.Println("We get an area", fw, "by", fh, "coords", minX, maxX, minY, maxY)

	// initialize region map
	region := make(Field, fw)
	for i := range region {
		region[i] = make([]byte, fh)
	}

	// find out who's on the frontier - for those elements areas are infinite
	var nonInfinites []byte
	for _, v := range inData {
		v.x -= minX
		v.y -= minY
		region[v.x][v.y] = v.name
		if !(v.x == 0 || v.x == fw-1 || v.y == 0 || v.y == fh-1) {
			nonInfinites = append(nonInfinites, v.name)
		}
	}
	printField(region)
	fmt.Println(string(nonInfinites))

	// n^2 is bad. but... make it work first
	/*
		for i := 0; i < len(nonInfinites); i++ {
			for j = i + 1; j < len(nonInfinites); j++ {
				// get
			}
		}
	*/
}

func printField(data Field) {
	for _, xv := range data {
		for _, xy := range xv {
			if xy == 0 {
				fmt.Print(".")
			} else {
				fmt.Print(string(xy))
			}
		}
		fmt.Println()
	}
}

/*
--- Day 6: Chronal Coordinates ---

The device on your wrist beeps several times, and once again you feel like you're falling.

"Situation critical," the device announces. "Destination indeterminate. Chronal interference detected. Please specify new target coordinates."

The device then produces a list of coordinates (your puzzle input). Are they places it thinks are safe or dangerous? It recommends you check manual page 729. The Elves did not give you a manual.

If they're dangerous, maybe you can minimize the danger by finding the coordinate that gives the largest distance from the other points.

Using only the Manhattan distance, determine the area around each coordinate by counting the number of integer X,Y locations that are closest to that coordinate (and aren't tied in distance to any other coordinate).

Your goal is to find the size of the largest area that isn't infinite. For example, consider the following list of coordinates:

1, 1
1, 6
8, 3
3, 4
5, 5
8, 9

If we name these coordinates A through F, we can draw them on a grid, putting 0,0 at the top left:

..........
.A........
..........
........C.
...D......
.....E....
.B........
..........
..........
........F.

This view is partial - the actual grid extends infinitely in all directions. Using the Manhattan distance, each location's closest coordinate can be determined, shown here in lowercase:

aaaaa.cccc
aAaaa.cccc
aaaddecccc
aadddeccCc
..dDdeeccc
bb.deEeecc
bBb.eeee..
bbb.eeefff
bbb.eeffff
bbb.ffffFf

Locations shown as . are equally far from two or more coordinates, and so they don't count as being closest to any.

In this example, the areas of coordinates A, B, C, and F are infinite - while not shown here, their areas extend forever outside the visible grid. However, the areas of coordinates D and E are finite: D is closest to 9 locations, and E is closest to 17 (both including the coordinate's location itself). Therefore, in this example, the size of the largest area is 17.

What is the size of the largest area that isn't infinite?

*/
