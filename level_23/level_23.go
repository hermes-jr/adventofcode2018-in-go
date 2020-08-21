package main

import (
	. "../utils"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

type Point3D struct {
	x, y, z int64
}

type SearchMeta struct {
	loc      Point3D
	count    int
	toOrigin int64
}

type Drone struct {
	coordinate Point3D
	rng        int64
}

func main() {
	DEBUG = false
	var drones []Drone
	fileLines := ReadFile("input")
	//fileLines = ReadFile("input_test1")
	//fileLines = ReadFile("input_test2")

	re := regexp.MustCompile("^pos=<(-?\\d+),(-?\\d+),(-?\\d+)>, r=(\\d+)$")
	for _, l := range fileLines {
		match := re.FindStringSubmatch(l)
		x, _ := strconv.ParseInt(match[1], 10, 64)
		y, _ := strconv.ParseInt(match[2], 10, 64)
		z, _ := strconv.ParseInt(match[3], 10, 64)
		r, _ := strconv.ParseInt(match[4], 10, 64)
		drones = append(drones, Drone{Point3D{x, y, z}, r})
	}

	IfDebugPrintln(drones)
	if len(drones) == 0 {
		fmt.Println("ERROR: NO DATA")
		return
	}

	var maxD Drone
	for _, drone := range drones {
		if drone.rng >= maxD.rng {
			maxD = drone
		}
	}
	IfDebugPrintln("Strongest:", maxD)
	result1 := countInRangeOfStrongest(&drones, &maxD)
	fmt.Println("Result1", result1)

	result2 := dissectCube(drones)
	fmt.Println("Result2", result2)
}

func dissectCube(drones []Drone) int64 {
	dNum := len(drones)
	xs := make([]int64, dNum)
	ys := make([]int64, dNum)
	zs := make([]int64, dNum)

	for _, drone := range drones {
		xs = append(xs, drone.coordinate.x)
		ys = append(ys, drone.coordinate.y)
		zs = append(zs, drone.coordinate.z)
	}

	sort.Slice(xs, func(i, j int) bool {
		return xs[i] < xs[j]
	})
	sort.Slice(ys, func(i, j int) bool {
		return ys[i] < ys[j]
	})
	sort.Slice(zs, func(i, j int) bool {
		return zs[i] < zs[j]
	})
	var stepping int64
	for stepping = 1; stepping < (xs[len(xs)-1]-xs[0]) || stepping < ys[len(ys)-1]-ys[0] || stepping < zs[len(zs)-1]-zs[0]; stepping *= 2 {
	}
	IfDebugPrintln("Dist", stepping)

	span := 1
	for ; span < len(drones); span *= 2 {
	}
	forcedCheck := 1
	tried := make(map[int]SearchMeta)

	var bestVal int64
	bestCount := -1

	for {
		if _, ok := tried[forcedCheck]; !ok {
			tried[forcedCheck] = search(&drones, &xs, &ys, &zs, stepping, forcedCheck)
		}
		testVal, testCount := tried[forcedCheck].toOrigin, tried[forcedCheck].count
		IfDebugPrintln("Count:", testCount, "stepping", stepping)

		if testCount == -1 {
			if span > 1 {
				span /= 2
			}
			forcedCheck = max(1, forcedCheck-span)
		} else {
			if bestCount == -1 || bestCount < testCount {
				bestCount = testCount
				bestVal = testVal
			}
			if span == 1 {
				break
			}
			forcedCheck += span
		}
	}

	IfDebugPrintln("Max drones visible:", bestCount)
	return bestVal
}

func search(drones *[]Drone, xs, ys, zs *[]int64, dist int64, forcedCount int) SearchMeta {
	var atTarget []SearchMeta
	IfDebugPrintln(xs, ys, zs)
	for x := (*xs)[0]; x < (*xs)[len(*xs)-1]+1; x += dist {
		for y := (*ys)[0]; y < (*ys)[len(*ys)-1]+1; y += dist {
			for z := (*zs)[0]; z < (*zs)[len(*zs)-1]+1; z += dist {

				curLoc := Point3D{x, y, z}
				IfDebugPrintln("Analyzing", curLoc)

				count := 0
				for _, drone := range *drones {
					calc := mhd(&curLoc, &drone.coordinate)
					if dist == 1 {
						if calc <= drone.rng {
							count++
						}
					} else {
						if calc/dist-3 <= drone.rng/dist {
							count++
						}
					}
				}

				if count >= forcedCount {
					atTarget = append(atTarget, SearchMeta{curLoc, count, distanceToOrigin(&curLoc)})
				}

			}
		}
	}

	for len(atTarget) > 0 {
		var best SearchMeta
		bestIdx := -1

		for k, v := range atTarget {
			if bestIdx == -1 || v.toOrigin < best.toOrigin {
				best = v
				bestIdx = k
			}
		}

		if dist == 1 {
			return best
		} else {
			xs := []int64{best.loc.x, best.loc.x + dist/2}
			ys := []int64{best.loc.y, best.loc.y + dist/2}
			zs := []int64{best.loc.z, best.loc.z + dist/2}
			subResult := search(drones, &xs, &ys, &zs, dist/2, forcedCount)
			if subResult.count == -1 {
				atTarget[bestIdx] = atTarget[len(atTarget)-1] // Copy last element to index i.
				atTarget[len(atTarget)-1] = SearchMeta{}      // Erase last element (write zero value).
				atTarget = atTarget[:len(atTarget)-1]
			} else {
				return subResult
			}
		}

	}

	return SearchMeta{
		toOrigin: 0,
		count:    -1,
	}
}

func distanceToOrigin(point *Point3D) int64 {
	return abs(point.x) + abs(point.y) + abs(point.z)
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func countInRangeOfStrongest(drones *[]Drone, maxD *Drone) int {
	result := 0
	for _, drone := range *drones {
		if mhd(&drone.coordinate, &maxD.coordinate) <= maxD.rng {
			result++
		}
	}
	return result
}

func mhd(a, b *Point3D) int64 {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

/*
--- Day 23: Experimental Emergency Teleportation ---

Using your torch to search the darkness of the rocky cavern, you finally locate the man's friend: a small reindeer.

You're not sure how it got so far in this cave. It looks sick - too sick to walk - and too heavy for you to carry all the way back. Sleighs won't be invented for another 1500 years, of course.

The only option is experimental emergency teleportation.

You hit the "experimental emergency teleportation" button on the device and push I accept the risk on no fewer than 18 different warning messages. Immediately, the device deploys hundreds of tiny nanobots which fly around the cavern, apparently assembling themselves into a very specific formation. The device lists the X,Y,Z position (pos) for each nanobot as well as its signal radius (r) on its tiny screen (your puzzle input).

Each nanobot can transmit signals to any integer coordinate which is a distance away from it less than or equal to its signal radius (as measured by Manhattan distance). Coordinates a distance away of less than or equal to a nanobot's signal radius are said to be in range of that nanobot.

Before you start the teleportation process, you should determine which nanobot is the strongest (that is, which has the largest signal radius) and then, for that nanobot, the total number of nanobots that are in range of it, including itself.

For example, given the following nanobots:

pos=<0,0,0>, r=4
pos=<1,0,0>, r=1
pos=<4,0,0>, r=3
pos=<0,2,0>, r=1
pos=<0,5,0>, r=3
pos=<0,0,3>, r=1
pos=<1,1,1>, r=1
pos=<1,1,2>, r=1
pos=<1,3,1>, r=1

The strongest nanobot is the first one (position 0,0,0) because its signal radius, 4 is the largest. Using that nanobot's location and signal radius, the following nanobots are in or out of range:

    The nanobot at 0,0,0 is distance 0 away, and so it is in range.
    The nanobot at 1,0,0 is distance 1 away, and so it is in range.
    The nanobot at 4,0,0 is distance 4 away, and so it is in range.
    The nanobot at 0,2,0 is distance 2 away, and so it is in range.
    The nanobot at 0,5,0 is distance 5 away, and so it is not in range.
    The nanobot at 0,0,3 is distance 3 away, and so it is in range.
    The nanobot at 1,1,1 is distance 3 away, and so it is in range.
    The nanobot at 1,1,2 is distance 4 away, and so it is in range.
    The nanobot at 1,3,1 is distance 5 away, and so it is not in range.

In this example, in total, 7 nanobots are in range of the nanobot with the largest signal radius.

Find the nanobot with the largest signal radius. How many nanobots are in range of its signals?

Your puzzle answer was 172.

The first half of this puzzle is complete! It provides one gold star: *
--- Part Two ---

Now, you just need to figure out where to position yourself so that you're actually teleported when the nanobots activate.

To increase the probability of success, you need to find the coordinate which puts you in range of the largest number of nanobots. If there are multiple, choose one closest to your position (0,0,0, measured by manhattan distance).

For example, given the following nanobot formation:

pos=<10,12,12>, r=2
pos=<12,14,12>, r=2
pos=<16,12,12>, r=4
pos=<14,14,14>, r=6
pos=<50,50,50>, r=200
pos=<10,10,10>, r=5

Many coordinates are in range of some of the nanobots in this formation. However, only the coordinate 12,12,12 is in range of the most nanobots: it is in range of the first five, but is not in range of the nanobot at 10,10,10. (All other coordinates are in range of fewer than five nanobots.) This coordinate's distance from 0,0,0 is 36.

Find the coordinates that are in range of the largest number of nanobots. What is the shortest manhattan distance between any of those points and 0,0,0?

*/
