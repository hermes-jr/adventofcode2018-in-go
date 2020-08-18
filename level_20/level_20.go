package main

import (
	. "../utils"
	"fmt"
	"github.com/thcyron/graphs"
)

type Point2D struct {
	x, y int
}

func add(a, b *Point2D) Point2D {
	return Point2D{
		x: a.x + b.x,
		y: a.y + b.y,
	}
}

var rose = map[rune]Point2D{'N': {0, -1},
	'S': {0, 1}, 'E': {1, 0}, 'W': {-1, 0}}

var g = graphs.NewGraph()

func main() {
	var entryPoint = Point2D{
		x: 0,
		y: 0,
	}
	DEBUG = false
	directions := ReadFile("input")[0]
	//directions = "^WNE$"
	//directions = "^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$" //18
	//directions = "^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$" // 23
	//directions = "^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$" // 31
	//directions = "^NNNNN(EEEEE|NNN)NNNNN$" // 15 (?)

	IfDebugPrintln(directions)
	follow(directions[1:len(directions)-1], entryPoint)
	IfDebugPrintln("g", g)

	result1 := 0
	result2 := 0
	for vertex := range g.VerticesIter() {
		dk := graphs.Dijkstra(g, entryPoint, vertex)
		curDist := dk.Len() - 1
		if curDist >= 1000 {
			result2++
		}
		IfDebugPrintf("Distance to %v: %v\n", vertex, curDist)
		if curDist >= result1 {
			result1 = curDist
		}
	}

	fmt.Println("Result1", result1)
	fmt.Println("Result2", result2)
}

func follow(directions string, curLoc Point2D) Point2D {
	splitPoint := curLoc
	IfDebugPrintln("At", curLoc, "processing", directions)
	for i := 0; i < len(directions); i++ {
		c := rune(directions[i])
		if dif, ok := rose[c]; ok {
			nextLoc := add(&curLoc, &dif)
			g.AddEdge(curLoc, nextLoc, 1)
			IfDebugPrintln("Processing", string(c), "moving from", curLoc, "to", nextLoc)
			curLoc = nextLoc
		} else if c == '(' {
			groupEnd := i + 1
			for depth := 1; depth > 0; groupEnd++ {
				if directions[groupEnd] == ')' {
					depth--
				} else if directions[groupEnd] == '(' {
					depth++
				}
			}
			curLoc = follow(directions[i+1:groupEnd-1], curLoc)
			i += groupEnd - i - 1
			continue
		} else if c == '|' {
			IfDebugPrintln("Processing |, returning from", curLoc, "to", splitPoint)
			curLoc = splitPoint
		} else {
			fmt.Println("UNKNOWN SYMBOL", string(c))
			break
		}
	}
	return curLoc
}
