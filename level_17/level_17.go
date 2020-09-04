package main

import (
	. "../utils"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
)

type Point2D struct {
	x, y int
}

var minX, minY, maxX, maxY = math.MaxInt32, math.MaxInt32, math.MinInt32, math.MinInt32

func main() {
	DEBUG = false
	inputLines := ReadFile("input")
	//inputLines = ReadFile("input_test1")
	//inputLines = ReadFile("input_test2")
	//inputLines = ReadFile("input_test3")
	//inputLines = ReadFile("input_test4")
	//inputLines = ReadFile("input_test5")

	data := parseInput(inputLines)

	IfDebugPrintln(Point2D{minX, minY}, ":", Point2D{maxX, maxY})
	drawMap(&data)

	for retry := true; retry; retry = flow(&data, Point2D{500, 1}) {
	}

	result1 := 0
	result2 := 0
	for k, v := range data {
		if k.y > maxY || k.y < minY {
			continue
		}
		switch v {
		case '~':
			result1++
			result2++
		case '|':
			result1++
		}
	}
	drawMap(&data)
	fmt.Println("Result1", result1)
	fmt.Println("Result2", result2)
}

func flow(data *map[Point2D]rune, entry Point2D) bool {
	IfDebugPrintln("At", entry)
	if entry.x < minX-1 || entry.x > maxX+1 || entry.y > maxY || entry.y < 0 ||
		!canPass(data, &entry) {
		return false // out of map bounds or impenetrable tile
	}
	(*data)[entry] = '|'
	drawMap(data)

	down := Point2D{entry.x, entry.y + 1}

	// try down
	if canPass(data, &down) {
		return flow(data, down)
	} else {
		// check if can settle
		lb, rb := entry.x, entry.x
		for nx := entry.x; nx >= minX-1; nx-- {
			if canPass(data, &Point2D{nx, entry.y + 1}) {
				// water will fall on the left side
				lb = -1
				break
			} else {
				if canPass(data, &Point2D{nx, entry.y}) {
					lb = nx // can flow left
				} else {
					break // can't flow left anymore
				}
			}
		}
		for nx := entry.x; nx <= maxX+1; nx++ {
			if canPass(data, &Point2D{nx, entry.y + 1}) {
				// water will fall on the right side
				rb = -1
				break
			} else {
				if canPass(data, &Point2D{nx, entry.y}) {
					rb = nx // can flow right
				} else {
					break // can't flow right anymore
				}
			}
		}
		right := Point2D{entry.x + 1, entry.y}
		left := Point2D{entry.x - 1, entry.y}
		if lb != -1 && rb != -1 {
			// both borders present, fill layer with settled water
			IfDebugPrintf("Can settle layer; y:%d, x:%d..%d\n", entry.y, lb, rb)
			drawMap(data)

			for nx := lb; nx <= rb; nx++ {
				(*data)[Point2D{nx, entry.y}] = '~' // settled
			}

			// Couldn't properly backtrack, so there's a dirty workaround:
			// Dry everything, retry from the beginning
			for k, v := range *data {
				if v == '|' {
					delete(*data, k)
				}
			}
			return true
		} else {
			// one or both borders missing
			retries := false
			if getTileType(data, &left) != '|' {
				retries = flow(data, left) || retries
			}
			if getTileType(data, &right) != '|' {
				retries = flow(data, right) || retries
			}
			return retries
		}
	}
}

func canPass(data *map[Point2D]rune, t *Point2D) bool {
	tileType := getTileType(data, t)
	if tileType == '.' || tileType == '|' {
		return true
	}
	return false
}

func getTileType(data *map[Point2D]rune, location *Point2D) rune {
	if tileType, present := (*data)[*location]; !present {
		return '.'
	} else {
		return tileType
	}
}

func parseInput(inputLines []string) map[Point2D]rune {
	data := make(map[Point2D]rune)

	re := regexp.MustCompile(`^[x,y]=(\d+), [x,y]=(\d+)..(\d+)$`)

	for _, line := range inputLines {
		direction := line[0] == 'y' // true = horizontal, false = vertical

		match := re.FindStringSubmatch(line)
		if match == nil {
			log.Fatal("Couldn't parse: ", line)
		}

		fixed, _ := strconv.Atoi(match[1])
		rl, _ := strconv.Atoi(match[2])
		rr, _ := strconv.Atoi(match[3])

		if direction {
			minY = min(minY, fixed)
			maxY = max(maxY, fixed)
			minX = min(minX, rl)
			maxX = max(maxX, rr)
		} else {
			minX = min(minX, fixed)
			maxX = max(maxX, fixed)
			minY = min(minY, rl)
			maxY = max(maxY, rr)
		}

		for i := rl; i <= rr; i++ {
			if direction {
				data[Point2D{i, fixed}] = '#'
			} else {
				data[Point2D{fixed, i}] = '#'
			}
		}
	}
	return data
}

func drawMap(data *map[Point2D]rune) {
	if !DEBUG {
		return
	}
	println()
	for j := -1; j <= maxY-minY; j++ {
		for i := -1; i <= maxX-minX+1; i++ {
			curPoint := Point2D{i + minX, j + minY}
			if i == 500-minX && j == -1 {
				fmt.Print("+") // spring
			} else {
				fmt.Print(string(getTileType(data, &curPoint)))
			}
		}
		fmt.Println()
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

/*
--- Day 17: Reservoir Research ---

You arrive in the year 18. If it weren't for the coat you got in 1018, you would be very cold: the North Pole base hasn't even been constructed.

Rather, it hasn't been constructed yet. The Elves are making a little progress, but there's not a lot of liquid water in this climate, so they're getting very dehydrated. Maybe there's more underground?

You scan a two-dimensional vertical slice of the ground nearby and discover that it is mostly sand with veins of clay. The scan only provides data with a granularity of square meters, but it should be good enough to determine how much water is trapped there. In the scan, x represents the distance to the right, and y represents the distance down. There is also a spring of water near the surface at x=500, y=0. The scan identifies which square meters are clay (your puzzle input).

For example, suppose your scan shows the following veins of clay:

x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504

Rendering clay as #, sand as ., and the water spring as +, and with x increasing to the right and y increasing downward, this becomes:

   44444455555555
   99999900000000
   45678901234567
 0 ......+.......
 1 ............#.
 2 .#..#.......#.
 3 .#..#..#......
 4 .#..#..#......
 5 .#.....#......
 6 .#.....#......
 7 .#######......
 8 ..............
 9 ..............
10 ....#.....#...
11 ....#.....#...
12 ....#.....#...
13 ....#######...

The spring of water will produce water forever. Water can move through sand, but is blocked by clay. Water always moves down when possible, and spreads to the left and right otherwise, filling space that has clay on both sides and falling out otherwise.

For example, if five squares of water are created, they will flow downward until they reach the clay and settle there. Water that has come to rest is shown here as ~, while sand through which water has passed (but which is now dry again) is shown as |:

......+.......
......|.....#.
.#..#.|.....#.
.#..#.|#......
.#..#.|#......
.#....|#......
.#~~~~~#......
.#######......
..............
..............
....#.....#...
....#.....#...
....#.....#...
....#######...

Two squares of water can't occupy the same location. If another five squares of water are created, they will settle on the first five, filling the clay reservoir a little more:

......+.......
......|.....#.
.#..#.|.....#.
.#..#.|#......
.#..#.|#......
.#~~~~~#......
.#~~~~~#......
.#######......
..............
..............
....#.....#...
....#.....#...
....#.....#...
....#######...

Water pressure does not apply in this scenario. If another four squares of water are created, they will stay on the right side of the barrier, and no water will reach the left side:

......+.......
......|.....#.
.#..#.|.....#.
.#..#~~#......
.#..#~~#......
.#~~~~~#......
.#~~~~~#......
.#######......
..............
..............
....#.....#...
....#.....#...
....#.....#...
....#######...

At this point, the top reservoir overflows. While water can reach the tiles above the surface of the water, it cannot settle there, and so the next five squares of water settle like this:

......+.......
......|.....#.
.#..#||||...#.
.#..#~~#|.....
.#..#~~#|.....
.#~~~~~#|.....
.#~~~~~#|.....
.#######|.....
........|.....
........|.....
....#...|.#...
....#...|.#...
....#~~~~~#...
....#######...

Note especially the leftmost |: the new squares of water can reach this tile, but cannot stop there. Instead, eventually, they all fall to the right and settle in the reservoir below.

After 10 more squares of water, the bottom reservoir is also full:

......+.......
......|.....#.
.#..#||||...#.
.#..#~~#|.....
.#..#~~#|.....
.#~~~~~#|.....
.#~~~~~#|.....
.#######|.....
........|.....
........|.....
....#~~~~~#...
....#~~~~~#...
....#~~~~~#...
....#######...

Finally, while there is nowhere left for the water to settle, it can reach a few more tiles before overflowing beyond the bottom of the scanned data:

......+.......    (line not counted: above minimum y value)
......|.....#.
.#..#||||...#.
.#..#~~#|.....
.#..#~~#|.....
.#~~~~~#|.....
.#~~~~~#|.....
.#######|.....
........|.....
...|||||||||..
...|#~~~~~#|..
...|#~~~~~#|..
...|#~~~~~#|..
...|#######|..
...|.......|..    (line not counted: below maximum y value)
...|.......|..    (line not counted: below maximum y value)
...|.......|..    (line not counted: below maximum y value)

How many tiles can be reached by the water? To prevent counting forever, ignore tiles with a y coordinate smaller than the smallest y coordinate in your scan data or larger than the largest one. Any x coordinate is valid. In this example, the lowest y coordinate given is 1, and the highest is 13, causing the water spring (in row 0) and the water falling off the bottom of the render (in rows 14 through infinity) to be ignored.

So, in the example above, counting both water at rest (~) and other sand tiles the water can

hypothetically reach (|), the total number of tiles the water can reach is 57.

How many tiles can the water reach within the range of y values in your scan?

*/
