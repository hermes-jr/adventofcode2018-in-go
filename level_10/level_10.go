package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	xpos, ypos, xvel, yvel int
}

func main() {
	fname := "input"
	//fname = "input_test"

	file, _ := os.Open(fname)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var data []Point

	for scanner.Scan() {
		inputLine := scanner.Text()
		fmt.Println(inputLine)
		// position=<-50212, -20003> velocity=< 5,  2>
		// position=< 3,  6> velocity=<-1, -1>
		re := regexp.MustCompile("^position=<\\s*(-?\\d+),\\s*(-?\\d+)> velocity=<\\s*(-?\\d+),\\s*(-?\\d+)>$")
		match := re.FindStringSubmatch(inputLine)
		if match == nil {
			log.Fatal("Couldn't parse", inputLine)
		}
		xpos, _ := strconv.Atoi(match[1])
		ypos, _ := strconv.Atoi(match[2])
		xvel, _ := strconv.Atoi(match[3])
		yvel, _ := strconv.Atoi(match[4])

		p := Point{xpos: xpos, ypos: ypos, xvel: xvel, yvel: yvel}
		fmt.Println("Point read", p)
		data = append(data, p)
	}

	// Wait until lights spread decreases, as soon as it starts to increase - we've passed the right step
	var lastSurf int
	for step := 0; ; step++ {
		xsize, ysize, _, _ := calcFieldSize(data, step)
		curSurf := xsize * ysize
		if lastSurf == 0 || curSurf < lastSurf {
			lastSurf = curSurf
		} else {
			// overshoot
			step--
			fmt.Println("Result1:")
			takeScreenshot(data, step)
			fmt.Println("Result2:", step)
			break
		}
	}

}

// Get field size at a specific step
func calcFieldSize(data []Point, step int) (int, int, int, int) {
	minx := data[0].xpos + step*data[0].xvel
	maxx := minx
	miny := data[0].ypos + step*data[0].yvel
	maxy := miny

	for _, p := range data {
		xpos := p.xpos + step*p.xvel
		ypos := p.ypos + step*p.yvel
		maxx = maxOf(xpos, maxx)
		maxy = maxOf(ypos, maxy)
		minx = minOf(xpos, minx)
		miny = minOf(ypos, miny)
	}
	return maxx - minx + 1, maxy - miny + 1, minx, miny
}

// Print points clipping area to min-max coordinates
func takeScreenshot(data []Point, step int) {
	xsize, ysize, minx, miny := calcFieldSize(data, step)
	fmt.Printf("field [%v; %v]: minx:%v, miny:%v\n", xsize, ysize, minx, miny)

	region := make([][]bool, ysize)
	for i := range region {
		region[i] = make([]bool, xsize)
	}

	for _, p := range data {
		pointCoordX := p.xpos + p.xvel*step
		pointCoordY := p.ypos + p.yvel*step
		region[pointCoordY-miny][pointCoordX-minx] = true
	}
	for i := 0; i < len(region); i++ {
		for j := 0; j < len(region[i]); j++ {
			if region[i][j] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
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
--- Day 10: The Stars Align ---

It's no use; your navigation system simply isn't capable of providing walking directions in the arctic circle, and certainly not in 1018.

The Elves suggest an alternative. In times like these, North Pole rescue operations will arrange points of light in the sky to guide missing Elves back to base. Unfortunately, the message is easy to miss: the points move slowly enough that it takes hours to align them, but have so much momentum that they only stay aligned for a second. If you blink at the wrong time, it might be hours before another message appears.

You can see these points of light floating in the distance, and record their position in the sky and their velocity, the relative change in position per second (your puzzle input). The coordinates are all given from your perspective; given enough time, those positions and velocities will move the points into a cohesive message!

Rather than wait, you decide to fast-forward the process and calculate what the points will eventually spell.

For example, suppose you note the following points:

position=< 9,  1> velocity=< 0,  2>
position=< 7,  0> velocity=<-1,  0>
position=< 3, -2> velocity=<-1,  1>
position=< 6, 10> velocity=<-2, -1>
position=< 2, -4> velocity=< 2,  2>
position=<-6, 10> velocity=< 2, -2>
position=< 1,  8> velocity=< 1, -1>
position=< 1,  7> velocity=< 1,  0>
position=<-3, 11> velocity=< 1, -2>
position=< 7,  6> velocity=<-1, -1>
position=<-2,  3> velocity=< 1,  0>
position=<-4,  3> velocity=< 2,  0>
position=<10, -3> velocity=<-1,  1>
position=< 5, 11> velocity=< 1, -2>
position=< 4,  7> velocity=< 0, -1>
position=< 8, -2> velocity=< 0,  1>
position=<15,  0> velocity=<-2,  0>
position=< 1,  6> velocity=< 1,  0>
position=< 8,  9> velocity=< 0, -1>
position=< 3,  3> velocity=<-1,  1>
position=< 0,  5> velocity=< 0, -1>
position=<-2,  2> velocity=< 2,  0>
position=< 5, -2> velocity=< 1,  2>
position=< 1,  4> velocity=< 2,  1>
position=<-2,  7> velocity=< 2, -2>
position=< 3,  6> velocity=<-1, -1>
position=< 5,  0> velocity=< 1,  0>
position=<-6,  0> velocity=< 2,  0>
position=< 5,  9> velocity=< 1, -2>
position=<14,  7> velocity=<-2,  0>
position=<-3,  6> velocity=< 2, -1>

Each line represents one point. Positions are given as <X, Y> pairs: X represents how far left (negative) or right (positive) the point appears, while Y represents how far up (negative) or down (positive) the point appears.

At 0 seconds, each point has the position given. Each second, each point's velocity is added to its position. So, a point with velocity <1, -2> is moving to the right, but is moving upward twice as quickly. If this point's initial position were <3, 9>, after 3 seconds, its position would become <6, 3>.

Over time, the points listed above would move like this:

Initially:
........#.............
................#.....
.........#.#..#.......
......................
#..........#.#.......#
...............#......
....#.................
..#.#....#............
.......#..............
......#...............
...#...#.#...#........
....#..#..#.........#.
.......#..............
...........#..#.......
#...........#.........
...#.......#..........

After 1 second:
......................
......................
..........#....#......
........#.....#.......
..#.........#......#..
......................
......#...............
....##.........#......
......#.#.............
.....##.##..#.........
........#.#...........
........#...#.....#...
..#...........#.......
....#.....#.#.........
......................
......................

After 2 seconds:
......................
......................
......................
..............#.......
....#..#...####..#....
......................
........#....#........
......#.#.............
.......#...#..........
.......#..#..#.#......
....#....#.#..........
.....#...#...##.#.....
........#.............
......................
......................
......................

After 3 seconds:
......................
......................
......................
......................
......#...#..###......
......#...#...#.......
......#...#...#.......
......#####...#.......
......#...#...#.......
......#...#...#.......
......#...#...#.......
......#...#..###......
......................
......................
......................
......................

After 4 seconds:
......................
......................
......................
............#.........
........##...#.#......
......#.....#..#......
.....#..##.##.#.......
.......##.#....#......
...........#....#.....
..............#.......
....#......#...#......
.....#.....##.........
...............#......
...............#......
......................
......................

After 3 seconds, the message appeared briefly: HI. Of course, your message will be much longer and will take many more seconds to appear.

What message will eventually appear in the sky?

--- Part Two ---

Good thing you didn't have to wait, because that would have taken a long time - much longer than the 3 seconds in the example above.

Impressed by your sub-hour communication capabilities, the Elves are curious: exactly how many seconds would they have needed to wait for that message to appear?

*/
