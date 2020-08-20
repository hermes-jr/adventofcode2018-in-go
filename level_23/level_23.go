package main

import (
	. "../utils"
	"fmt"
	"regexp"
	"strconv"
)

type Point3D struct {
	x, y, z int64
}

type Drone struct {
	coordinate Point3D
	rng        int64
}

func main() {
	DEBUG = true
	var drones []Drone
	fileLines := ReadFile("input")
	//fileLines = ReadFile("input_test")

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

	var maxD Drone
	for _, drone := range drones {
		if drone.rng >= maxD.rng {
			maxD = drone
		}
	}
	IfDebugPrintln("Strongest:", maxD)
	result1 := countInRangeOfStrongest(&drones, &maxD)
	fmt.Println("Result1", result1)
}

func countInRangeOfStrongest(drones *[]Drone, maxD *Drone) int {
	result := 0
	for _, drone := range *drones {
		if mhd(drone.coordinate, maxD.coordinate) <= maxD.rng {
			result++
		}
	}
	return result
}

func mhd(a, b Point3D) int64 {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}
