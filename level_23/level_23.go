package main

import (
	. "../utils"
	"container/heap"
	"fmt"
	"regexp"
	"strconv"
)

type Point3D struct {
	x, y, z int64
}

type SearchMeta struct {
	loc          []Point3D
	maxSeenSoFar int
}

type Drone struct {
	coordinate Point3D
	rng        int64
}

// Priority queue (golang.org)
// An Item is something we manage in a priority queue.
type Item struct {
	value    [2]int64 // The value of the item; arbitrary.
	priority int64    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value [2]int64, priority int64) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
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

	minX, minY, minZ, maxX, maxY, maxZ, minR, maxR :=
		drones[0].coordinate.x, drones[0].coordinate.y, drones[0].coordinate.z,
		drones[0].coordinate.x, drones[0].coordinate.y, drones[0].coordinate.z,
		drones[0].rng, drones[0].rng
	var maxD Drone
	for _, drone := range drones {
		if drone.rng >= maxD.rng {
			maxD = drone
		}
		minMax(&drone, &minX, &maxX, func(fd *Drone) int64 { return fd.coordinate.x })
		minMax(&drone, &minY, &maxY, func(fd *Drone) int64 { return fd.coordinate.y })
		minMax(&drone, &minZ, &maxZ, func(fd *Drone) int64 { return fd.coordinate.z })
		minMax(&drone, &minR, &maxR, func(fd *Drone) int64 { return fd.rng })
	}
	IfDebugPrintln("Strongest:", maxD)
	result1 := countInRangeOfStrongest(&drones, &maxD)
	fmt.Println("Result1", result1)

	result2a := queue(drones)
	fmt.Println("Result2 (queue)", result2a)
	// 22164451 low
	// 92358683 incorrect
	// 125532606 incorrect

	result2b := dissectCube(min(min(minX, minY), minZ), max(max(maxX, maxY), maxZ), minR, drones)
	fmt.Println("Result2 (cube dissect)", result2b)
}

func queue(drones []Drone) int64 {
	pq := make(PriorityQueue, len(drones)*2)
	for i, drone := range drones {
		distance := abs(drone.coordinate.x) + abs(drone.coordinate.y) + abs(drone.coordinate.z)
		pq[i*2] = &Item{
			value:    [2]int64{max(0, distance-drone.rng), 1},
			priority: max(0, distance-drone.rng),
			index:    i * 2,
		}
		pq[i*2+1] = &Item{
			value:    [2]int64{distance + drone.rng + 1, -1},
			priority: distance + drone.rng + 1,
			index:    i*2 + 1,
		}
	}
	heap.Init(&pq)

	var count, maxCount, result int64

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%v ", item.priority, item.value)
		count += item.value[1]
		if count > maxCount {
			result = item.value[0]
			maxCount = count
		}
	}
	fmt.Println()
	return result
}

func dissectCube(minD int64, maxD int64, minR int64, drones []Drone) int64 {
	IfDebugPrintf("World [%v:%v:%v] - [%v:%v:%v]\n", minD, minD, minD, maxD, maxD, maxD)

	// Smallest step to split volume (each coverage cube is guaranteed to be hit):
	IfDebugPrintln("Minimum range", minR)

	hotspots := search(&drones, &Point3D{minD, minD, minD},
		&Point3D{maxD, maxD, maxD}, minR, 0)
	IfDebugPrintln(hotspots)

	result2 := distanceToOrigin(&hotspots.loc[0])
	for _, p := range hotspots.loc {
		distance := abs(p.x) + abs(p.y) + abs(p.z)
		if distance < result2 {
			result2 = distance
		}
	}
	return result2
}

func search(drones *[]Drone, minCorner, maxCorner *Point3D, step int64, maxSeenSoFar int) SearchMeta {
	var result []Point3D
	for x := minCorner.x; x < maxCorner.x; x += step {
		for y := minCorner.y; y < maxCorner.y; y += step {
			for z := minCorner.z; z < maxCorner.z; z += step {
				curLoc := Point3D{x, y, z}
				reachable := pingDrones(drones, &curLoc)
				if reachable > 0 {
					IfDebugPrintln("At", curLoc, "there are", reachable, "drones reachable, stepping", step)
				}
				if reachable < maxSeenSoFar {
					continue // nothing interesting in this quadrant
				} else if reachable == maxSeenSoFar {
					// if stepping is 1, chose closest to origin right here
					if step == 1 && len(result) > 0 {
						if distanceToOrigin(&curLoc) < distanceToOrigin(&result[0]) {
							result = nil
							result = append(result, curLoc)
						}
					} else {
						result = append(result, curLoc)
					}
				} else {
					// new best bet
					result = nil
					maxSeenSoFar = reachable
					result = append(result, curLoc)
				}
			}
		}
	}
	if step == 1 {
		return SearchMeta{result, maxSeenSoFar}
	}

	var lr, rr int64
	for _, nc := range result {
		if step%2 == 0 {
			rr = step / 2
			lr = rr
		} else {
			rr = step/2 + 1
			lr = step / 2
		}
		stepDiv := max(1, step/2)
		IfDebugPrintln("Going to investigate subsample around", nc, "radius", lr, "|", rr, "stepDiv", stepDiv)
		subSample := search(drones, &Point3D{nc.x - lr, nc.y - lr, nc.z - lr},
			&Point3D{nc.x + rr, nc.y + rr, nc.z + rr}, stepDiv, maxSeenSoFar)
		IfDebugPrintln("Subsample investigated:", subSample, "ms: ", subSample.maxSeenSoFar, "stepDiv", stepDiv)
		if subSample.maxSeenSoFar < maxSeenSoFar {
			subSample.loc = nil
			continue
		} else if subSample.maxSeenSoFar == maxSeenSoFar {
			result = append(result, subSample.loc...)
		} else {
			result = subSample.loc
			maxSeenSoFar = subSample.maxSeenSoFar
		}
	}
	return SearchMeta{result, maxSeenSoFar} // shouldn't happen (?)
}

func distanceToOrigin(point *Point3D) int64 {
	return abs(point.x) + abs(point.y) + abs(point.z)
}

func max(a, b int64) int64 {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func pingDrones(drones *[]Drone, curLoc *Point3D) int {
	result := 0
	for _, drone := range *drones {
		if mhd(curLoc, &drone.coordinate) <= drone.rng {
			result++
		}
	}
	return result
}

func minMax(drone *Drone, minParam *int64, maxParam *int64, getParam func(*Drone) int64) {
	if getParam(drone) < *minParam {
		*minParam = getParam(drone)
	} else if getParam(drone) > *maxParam {
		*maxParam = getParam(drone)
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
