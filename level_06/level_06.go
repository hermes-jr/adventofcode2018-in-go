package main

import (
	. "../utils"
	"fmt"
	"regexp"
	"sort"
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

type DistancePair struct {
	dist int
	name byte
}

func (data DistancePair) String() string {
	return fmt.Sprintf("{%v %v}", string(data.name), data.dist)
}

type DistancePairs []DistancePair

func (p DistancePairs) Len() int           { return len(p) }
func (p DistancePairs) Less(i, j int) bool { return p[i].dist < p[j].dist }
func (p DistancePairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Points []Point
type Field [][]uint8

const DEBUG = false

func main() {
	fname := "input"
	//fname = "input_test"
	distlimit := 10000
	//distlimit = 32

	lines := ReadFile(fname)

	inData := Points{}

	minX, minY, maxX, maxY := -1, -1, 0, 0

	// scan input, find out "populated" area bounds
	for pointNameCode, i := byte(65), 0; i < len(lines); pointNameCode++ {
		inputLine := lines[i]
		i++
		re := regexp.MustCompile("^(\\d+), (\\d+)$")
		match := re.FindStringSubmatch(inputLine)
		y, _ := strconv.Atoi(match[1])
		x, _ := strconv.Atoi(match[2])
		inData = append(inData, Point{x, y, pointNameCode})
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
	if DEBUG {
		fmt.Println(inData)
	}

	fh := maxY - minY + 1
	fw := maxX - minX + 1
	if DEBUG {
		fmt.Println("We get an area", fw, "by", fh, "coords", minX, maxX, minY, maxY)
	}

	// initialize region map
	region := make(Field, fw)
	for i := range region {
		region[i] = make([]byte, fh)
	}

	// find out who's on the frontier - for those elements areas are infinite
	nonInfinite := make(map[byte]int)
	normalizedData := Points{}
	for _, v := range inData {
		v.x -= minX
		v.y -= minY
		region[v.x][v.y] = v.name
		normalizedData = append(normalizedData, v)
		if !(v.x == 0 || v.x == fw-1 || v.y == 0 || v.y == fh-1) {
			nonInfinite[v.name] = 0
		}
	}
	inData = normalizedData
	printField(region)
	if DEBUG {
		fmt.Println("Non-inf", nonInfinite)
	}

	// n^3 is bad. but... make it work first
	for i := 0; i < len(region); i++ {
		for j := 0; j < len(region[i]); j++ {
			distData := DistancePairs{}
			// find closest point or multiple points
			for _, v := range inData {
				dist := absInt(v.x-i) + absInt(v.y-j)
				distData = append(distData, DistancePair{dist, v.name})
			}
			sort.Sort(distData)
			if DEBUG {
				fmt.Println("At point", i, j, "distances are", distData)
			}
			if distData[0].dist == distData[1].dist {
				region[i][j] = 35
			} else {
				region[i][j] = distData[0].name
				nonInfinite[distData[0].name]++
			}
		}
	}

	if DEBUG {
		printField(region)
		fmt.Println("Non-inf updated", nonInfinite)
	}
	result1 := 0
	for _, v := range nonInfinite {
		if v > result1 {
			result1 = v
		}
	}
	fmt.Println("Result1", result1)

	result2 := 0
	// part 2, different algorithm
	for i := 0; i < len(region); i++ {
	cellLoop:
		for j := 0; j < len(region[i]); j++ {
			var sumForPoint int
			for _, v := range inData {
				dist := absInt(v.x-i) + absInt(v.y-j)
				sumForPoint += dist
				if sumForPoint >= distlimit {
					region[i][j] = 0
					continue cellLoop
				}
			}
			region[i][j] = 35
			result2++
		}
	}
	printField(region)
	fmt.Println("Result2", result2)

}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func printField(data Field) {
	if DEBUG {
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

A.......
........
.......C
..D.....
....E...
B.......
........
........
.......F

This view is partial - the actual grid extends infinitely in all directions. Using the Manhattan distance, each location's closest coordinate can be determined, shown here in lowercase:

Aaaa.ccc
aaddeccc
adddeccC
.dDdeecc
b.deEeec
Bb.eeee.
bb.eeeff
bb.eefff
bb.ffffF

Locations shown as . are equally far from two or more coordinates, and so they don't count as being closest to any.

In this example, the areas of coordinates A, B, C, and F are infinite - while not shown here, their areas extend forever outside the visible grid. However, the areas of coordinates D and E are finite: D is closest to 9 locations, and E is closest to 17 (both including the coordinate's location itself). Therefore, in this example, the size of the largest area is 17.

What is the size of the largest area that isn't infinite?

--- Part Two ---

On the other hand, if the coordinates are safe, maybe the best you can do is try to find a region near as many coordinates as possible.

For example, suppose you want the sum of the Manhattan distance to all of the coordinates to be less than 32. For each location, add up the distances to all of the given coordinates; if the total of those distances is less than 32, that location is within the desired region. Using the same coordinates as above, the resulting region looks like this:

..........
.A........
..........
...###..C.
..#D###...
..###E#...
.B.###....
..........
..........
........F.

In particular, consider the highlighted location 4,3 located at the top middle of the region. Its calculation is as follows, where abs() is the absolute value function:

    Distance to coordinate A: abs(4-1) + abs(3-1) =  5
    Distance to coordinate B: abs(4-1) + abs(3-6) =  6
    Distance to coordinate C: abs(4-8) + abs(3-3) =  4
    Distance to coordinate D: abs(4-3) + abs(3-4) =  2
    Distance to coordinate E: abs(4-5) + abs(3-5) =  3
    Distance to coordinate F: abs(4-8) + abs(3-9) = 10
    Total distance: 5 + 6 + 4 + 2 + 3 + 10 = 30

Because the total distance to all coordinates (30) is less than 32, the location is within the region.

This region, which also includes coordinates D and E, has a total size of 16.

Your actual region will need to be much larger than this example, though, instead including all locations with a total distance of less than 10000.

What is the size of the region containing all locations which have a total distance to all given coordinates of less than 10000?

*/
