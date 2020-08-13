package main

//goland:noinspection SpellCheckingInspection
import (
	. "../utils"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map"
	"sort"
)

type Point2d struct {
	x int
	y int
}

type Node struct {
	point    Point2d
	distance int
	isRoot   bool
}

type Unit struct {
	id, hp   int
	location Point2d
	race     byte
}

const attackPower = 3
const initialHp = 200

var levelMap [][]bool
var units []Unit

func (u Unit) String() string {
	return fmt.Sprintf("{%v (%v) id_%v}", string(u.race), u.hp, u.id)
}
func (p Point2d) String() string {
	return fmt.Sprintf("%v:%v", p.x, p.y)
}

func FindUnit(slice []Unit, val Unit) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func FindPoint(slice []Point2d, val Point2d) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func FindUnitByLocation(val Point2d) (int, bool) {
	for i := range units {
		if units[i].location == val && units[i].hp > 0 {
			return i, true
		}
	}
	return -1, false
}

func IsPointInQueue(slice []Node, val Point2d) bool {
	for _, item := range slice {
		if item.point == val {
			return true
		}
	}
	return false
}

func main() {
	DEBUG = false
	readInputFile()
	fmt.Println("Result1", fight())
}

func fight() int {
	turns := 0
	IfDebugPrintln("Initial map:")
	printMap(levelMap, units)

	for ; ; turns++ {
		IfDebugPrintln("Starting turn", turns+1)

		end := playTurn()

		printMap(levelMap, units)
		IfDebugPrintln("Turn complete", turns+1)

		if end {
			var healthSum int
			IfDebugPrintf("Game ended, survivors: ")
			for _, u := range units {
				if u.hp > 0 {
					IfDebugPrintf("%v", u)
					healthSum += u.hp
				}
			}
			IfDebugPrintln("Result1 formula:", healthSum, "*", turns)
			return healthSum * turns
		}
	}
}

func playTurn() bool {

	// Reorder units for current turn
	sort.SliceStable(units, func(i, j int) bool {
		if units[i].location.y != units[j].location.y {
			return units[i].location.y < units[j].location.y
		} else {
			return units[i].location.x < units[j].location.x
		}
	})
	IfDebugPrintln("Sorted units", units)

	for unitId := range units {
		if units[unitId].hp <= 0 {
			continue
		}
		if unitAct(unitId) {
			return true
		}
	}

	return false
}

func unitAct(unitId int) bool {
	unit := units[unitId]
	tryToReach := make(map[Point2d]bool)

	// Filter out allies
	enemies := append([]Unit(nil), units...)
	n := 0
	for id, x := range enemies {
		if x.race != unit.race && units[id].hp > 0 {
			enemies[n] = x
			n++
			for _, neighbor := range getNeighbors(x.location) {
				if isWalkable(unit.location, neighbor) {
					tryToReach[neighbor] = true
				}
			}
		}
	}
	enemies = enemies[:n]
	IfDebugPrintf("Targets for %v are: %v\n", unit, enemies)

	// No living enemies left
	if len(enemies) == 0 {
		return true
	}

	tryToReachKeys := make([]Point2d, len(tryToReach))
	i := 0
	for k := range tryToReach {
		tryToReachKeys[i] = k
		i++
	}
	IfDebugPrintf("In range for %v are: %v\n", unit, tryToReachKeys)

	if _, ok := FindPoint(tryToReachKeys, unit.location); !ok {
		nextLocation, ok := nextStep(unit.location, tryToReachKeys)
		if ok {
			IfDebugPrintln("Moving", unit, "to", nextLocation)
			units[unitId].location = nextLocation
			IfDebugPrintln("Alive units:", units)
		}
	}

	// Attack if possible
	var killable []Unit
	for _, neighbor := range getNeighbors(units[unitId].location) {
		if id, ok := FindUnitByLocation(neighbor); ok && units[id].race != unit.race {
			killable = append(killable, units[id])
		}
	}

	if len(killable) > 0 {
		// Attack weakest first, in reading order if there are multiple weakest
		sort.SliceStable(killable, func(i, j int) bool {
			if killable[i].hp != killable[j].hp {
				return killable[i].hp < killable[j].hp
			} else if killable[i].location.y == killable[j].location.y {
				return killable[i].location.x < killable[j].location.x
			} else {
				return killable[i].location.y < killable[j].location.y
			}
		})
		IfDebugPrintln("Targets for bashing:", killable)

		idToKill, _ := FindUnit(units, killable[0])
		units[idToKill].hp -= attackPower
		if units[idToKill].hp <= 0 {
			IfDebugPrintln(killable[0], "killed, remaining units:", units)
		}
	}
	return false
}

func readInputFile() {
	inputFile := "input"
	//inputFile = "input_test5"

	inputLines := ReadFile(inputFile)

	for rowNum, unitId, loc := 0, 0, 0; rowNum < len(inputLines); rowNum++ {
		levelMap = append(levelMap, []bool{})
		inputLine := []byte(inputLines[rowNum])
		for x, v := range inputLine {
			if v == 'E' || v == 'G' {
				units = append(units, Unit{unitId, initialHp, Point2d{x, rowNum}, v})
				unitId++
			}
			levelMap[rowNum] = append(levelMap[rowNum], v != '#')
			loc++
		}
	}
}

// BFS with reading priority
func nextStep(root Point2d, destinations []Point2d) (Point2d, bool) {
	IfDebugPrintln("--- BFS ---")
	var queue []Node
	discovered := make(map[Point2d]bool)
	discovered[root] = true

	breadcrumbs := orderedmap.New()

	queue = append(queue, Node{root, 0, false})
	for ok := true; ok; ok = len(queue) > 0 {
		curPos, dist := queue[0].point, queue[0].distance
		queue = queue[1:] // dequeue
		neighbors := getNeighbors(curPos)

		IfDebugPrintln("Current position", curPos, "neighbors:", neighbors)

		for _, neighbor := range neighbors {
			if !isWalkable(curPos, neighbor) {
				continue
			}
			if v, ok := breadcrumbs.Get(neighbor); !ok || v.(Node).distance > (dist+1) {
				breadcrumbs.Set(neighbor, Node{curPos, dist + 1, false})
			}
			if _, ok := discovered[neighbor]; ok {
				continue
			}
			IfDebugPrintln("Searching", neighbor, "in", queue)
			if !IsPointInQueue(queue, neighbor) {
				queue = append(queue, Node{neighbor, dist + 1, false})
				IfDebugPrintln("Not found, added:", queue)
			}

			discovered[curPos] = true
			IfDebugPrintln("Discovered:", discovered)
		}
	}

	minDist := -1
	var minPoint Point2d

	for pair := breadcrumbs.Oldest(); pair != nil; pair = pair.Next() {
		if _, ok := FindPoint(destinations, pair.Key.(Point2d)); !ok {
			// Path doesn't lead to an enemy
			continue
		}
		dist := pair.Value.(Node).distance
		if pair.Value.(Node).isRoot {
			IfDebugPrintln("--- END BFS ---")
			return Point2d{}, false
		}
		if minDist == -1 || // min not found yet
			(dist < minDist || // shorter path
				(dist == minDist && (pair.Key.(Point2d).y < minPoint.y || // same short path but lower in reading order
					(pair.Key.(Point2d).y == minPoint.y && pair.Key.(Point2d).x < minPoint.x)))) {
			minDist = dist
			minPoint = pair.Key.(Point2d)
		}
	}

	if DEBUG {
		for pair := breadcrumbs.Oldest(); pair != nil; pair = pair.Next() {
			fmt.Printf("%v => %v\n", pair.Key, pair.Value.(Node))
		}
	}

	if minDist == -1 {
		return minPoint, false
	}
	IfDebugPrintln("Closest target is at", minDist, "hops:", minPoint)

	// Track back
	for z, _ := breadcrumbs.Get(minPoint); z.(Node).distance > 1; z, _ = breadcrumbs.Get(minPoint) {
		minPoint = z.(Node).point
	}

	z, _ := breadcrumbs.Get(minPoint)
	IfDebugPrintln("--- END BFS ---")
	return minPoint, !z.(Node).isRoot
}

func isWalkable(curPos, point Point2d) bool {
	// Walkable: self location, empty points, dead units
	// Non-walkable: walls, live units
	if _, occupied := FindUnitByLocation(point); point == curPos ||
		(levelMap[point.y][point.x] && !occupied) {
		return true
	}
	return false
}

// Get adjacent coordinates (in reading order!)
func getNeighbors(p Point2d) [4]Point2d {
	return [...]Point2d{{p.x, p.y - 1}, {p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y + 1}}
}

// Print map with units
func printMap(data [][]bool, units []Unit) {
	if !DEBUG {
		return
	}
	fmt.Println()
	if len(data) == 0 {
		fmt.Println("-EMPTY MAP-")
		return
	}

	for y := 0; y < len(data); y++ {
		var rowUnits []Unit
		for x := 0; x < len(data[y]); x++ {
			symbol := byte('#')
			if u, found := FindUnitByLocation(Point2d{x, y}); found {
				rowUnits = append(rowUnits, units[u])
				symbol = units[u].race
			} else if data[y][x] {
				symbol = '.'
			}
			fmt.Print(string(symbol))
		}
		if len(rowUnits) > 0 {
			fmt.Print("  ", rowUnits)
		}
		fmt.Println()
	}
	fmt.Println()
}

/*
--- Day 15: Beverage Bandits ---

Having perfected their hot chocolate, the Elves have a new problem: the Goblins that live in these caves will do anything to steal it. Looks like they're here for a fight.

You scan the area, generating a map of the walls (#), open cavern (.), and starting position of every Goblin (G) and Elf (E) (your puzzle input).

Combat proceeds in rounds; in each round, each unit that is still alive takes a turn, resolving all of its actions before the next unit's turn begins. On each unit's turn, it tries to move into range of an enemy (if it isn't already) and then attack (if it is in range).

All units are very disciplined and always follow very strict combat rules. Units never move or attack diagonally, as doing so would be dishonorable. When multiple choices are equally valid, ties are broken in reading order: top-to-bottom, then left-to-right. For instance, the order in which units take their turns within a round is the reading order of their starting positions in that round, regardless of the type of unit or whether other units have moved after the round started. For example:

                 would take their
These units:   turns in this order:
  #######           #######
  #.G.E.#           #.1.2.#
  #E.G.E#           #3.4.5#
  #.G.E.#           #.6.7.#
  #######           #######

Each unit begins its turn by identifying all possible targets (enemy units). If no targets remain, combat ends.

Then, the unit identifies all of the open squares (.) that are in range of each target; these are the squares which are adjacent (immediately up, down, left, or right) to any target and which aren't already occupied by a wall or another unit. Alternatively, the unit might already be in range of a target. If the unit is not already in range of a target, and there are no open squares which are in range of a target, the unit ends its turn.

If the unit is already in range of a target, it does not move, but continues its turn with an attack. Otherwise, since it is not in range of a target, it moves.

To move, the unit first considers the squares that are in range and determines which of those squares it could reach in the fewest steps. A step is a single movement to any adjacent (immediately up, down, left, or right) open (.) square. Units cannot move into walls or other units. The unit does this while considering the current positions of units and does not do any prediction about where units will be later. If the unit cannot reach (find an open path to) any of the squares that are in range, it ends its turn. If multiple squares are in range and tied for being reachable in the fewest steps, the step which is first in reading order is chosen. For example:

Targets:      In range:     Reachable:    Nearest:      Chosen:
#######       #######       #######       #######       #######
#E..G.#       #E.?G?#       #E.@G.#       #E.!G.#       #E.+G.#
#...#.#  -->  #.?.#?#  -->  #.@.#.#  -->  #.!.#.#  -->  #...#.#
#.G.#G#       #?G?#G#       #@G@#G#       #!G.#G#       #.G.#G#
#######       #######       #######       #######       #######

In the above scenario, the Elf has three targets (the three Goblins):

    Each of the Goblins has open, adjacent squares which are in range (marked with a ? on the map).
    Of those squares, four are reachable (marked @); the other two (on the right) would require moving through a wall or unit to reach.
    Three of these reachable squares are nearest, requiring the fewest steps (only 2) to reach (marked !).
    Of those, the square which is first in reading order is chosen (+).

The unit then takes a single step toward the chosen square along the shortest path to that square. If multiple steps would put the unit equally closer to its destination, the unit chooses the step which is first in reading order. (This requires knowing when there is more than one shortest path so that you can consider the first step of each such path.) For example:

In range:     Nearest:      Chosen:       Distance:     Step:
#######       #######       #######       #######       #######
#.E...#       #.E...#       #.E...#       #4E212#       #..E..#
#...?.#  -->  #...!.#  -->  #...+.#  -->  #32101#  -->  #.....#
#..?G?#       #..!G.#       #...G.#       #432G2#       #...G.#
#######       #######       #######       #######       #######

The Elf sees three squares in range of a target (?), two of which are nearest (!), and so the first in reading order is chosen (+). Under "Distance", each open square is marked with its distance from the destination square; the two squares to which the Elf could move on this turn (down and to the right) are both equally good moves and would leave the Elf 2 steps from being in range of the Goblin. Because the step which is first in reading order is chosen, the Elf moves right one square.

Here's a larger example of movement:

Initially:
#########
#G..G..G#
#.......#
#.......#
#G..E..G#
#.......#
#.......#
#G..G..G#
#########

After 1 round:
#########
#.G...G.#
#...G...#
#...E..G#
#.G.....#
#.......#
#G..G..G#
#.......#
#########

After 2 rounds:
#########
#..G.G..#
#...G...#
#.G.E.G.#
#.......#
#G..G..G#
#.......#
#.......#
#########

After 3 rounds:
#########
#.......#
#..GGG..#
#..GEG..#
#G..G...#
#......G#
#.......#
#.......#
#########

Once the Goblins and Elf reach the positions above, they all are either in range of a target or cannot find any square in range of a target, and so none of the units can move until a unit dies.

After moving (or if the unit began its turn in range of a target), the unit attacks.

To attack, the unit first determines all of the targets that are in range of it by being immediately adjacent to it. If there are no such targets, the unit ends its turn. Otherwise, the adjacent target with the fewest hit points is selected; in a tie, the adjacent target with the fewest hit points which is first in reading order is selected.

The unit deals damage equal to its attack power to the selected target, reducing its hit points by that amount. If this reduces its hit points to 0 or fewer, the selected target dies: its square becomes . and it takes no further turns.

Each unit, either Goblin or Elf, has 3 attack power and starts with 200 hit points.

For example, suppose the only Elf is about to attack:

       HP:            HP:
G....  9       G....  9
..G..  4       ..G..  4
..EG.  2  -->  ..E..
..G..  2       ..G..  2
...G.  1       ...G.  1

The "HP" column shows the hit points of the Goblin to the left in the corresponding row. The Elf is in range of three targets: the Goblin above it (with 4 hit points), the Goblin to its right (with 2 hit points), and the Goblin below it (also with 2 hit points). Because three targets are in range, the ones with the lowest hit points are selected: the two Goblins with 2 hit points each (one to the right of the Elf and one below the Elf). Of those, the Goblin first in reading order (the one to the right of the Elf) is selected. The selected Goblin's hit points (2) are reduced by the Elf's attack power (3), reducing its hit points to -1, killing it.

After attacking, the unit's turn ends. Regardless of how the unit's turn ends, the next unit in the round takes its turn. If all units have taken turns in this round, the round ends, and a new round begins.

The Elves look quite outnumbered. You need to determine the outcome of the battle: the number of full rounds that were completed (not counting the round in which combat ends) multiplied by the sum of the hit points of all remaining units at the moment combat ends. (Combat only ends when a unit finds no targets during its turn.)

Below is an entire sample combat. Next to each map, each row's units' hit points are listed from left to right.

Initially:
#######
#.G...#   G(200)
#...EG#   E(200), G(200)
#.#.#G#   G(200)
#..G#E#   G(200), E(200)
#.....#
#######

After 1 round:
#######
#..G..#   G(200)
#...EG#   E(197), G(197)
#.#G#G#   G(200), G(197)
#...#E#   E(197)
#.....#
#######

After 2 rounds:
#######
#...G.#   G(200)
#..GEG#   G(200), E(188), G(194)
#.#.#G#   G(194)
#...#E#   E(194)
#.....#
#######

Combat ensues; eventually, the top Elf dies:

After 23 rounds:
#######
#...G.#   G(200)
#..G.G#   G(200), G(131)
#.#.#G#   G(131)
#...#E#   E(131)
#.....#
#######

After 24 rounds:
#######
#..G..#   G(200)
#...G.#   G(131)
#.#G#G#   G(200), G(128)
#...#E#   E(128)
#.....#
#######

After 25 rounds:
#######
#.G...#   G(200)
#..G..#   G(131)
#.#.#G#   G(125)
#..G#E#   G(200), E(125)
#.....#
#######

After 26 rounds:
#######
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(122)
#...#E#   E(122)
#..G..#   G(200)
#######

After 27 rounds:
#######
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(119)
#...#E#   E(119)
#...G.#   G(200)
#######

After 28 rounds:
#######
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(116)
#...#E#   E(113)
#....G#   G(200)
#######

More combat ensues; eventually, the bottom Elf dies:

After 47 rounds:
#######
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(59)
#...#.#
#....G#   G(200)
#######

Before the 48th round can finish, the top-left Goblin finds that there are no targets remaining, and so combat ends. So, the number of full rounds that were completed is 47, and the sum of the hit points of all remaining units is 200+131+59+200 = 590. From these, the outcome of the battle is 47 * 590 = 27730.

Here are a few example summarized combats:

#######       #######
#G..#E#       #...#E#   E(200)
#E#E.E#       #E#...#   E(197)
#G.##.#  -->  #.E##.#   E(185)
#...#E#       #E..#E#   E(200), E(200)
#...E.#       #.....#
#######       #######

Combat ends after 37 full rounds
Elves win with 982 total hit points left
Outcome: 37 * 982 = 36334

#######       #######
#E..EG#       #.E.E.#   E(164), E(197)
#.#G.E#       #.#E..#   E(200)
#E.##E#  -->  #E.##.#   E(98)
#G..#.#       #.E.#.#   E(200)
#..E#.#       #...#.#
#######       #######

Combat ends after 46 full rounds
Elves win with 859 total hit points left
Outcome: 46 * 859 = 39514

#######       #######
#E.G#.#       #G.G#.#   G(200), G(98)
#.#G..#       #.#G..#   G(200)
#G.#.G#  -->  #..#..#
#G..#.#       #...#G#   G(95)
#...E.#       #...G.#   G(200)
#######       #######

Combat ends after 35 full rounds
Goblins win with 793 total hit points left
Outcome: 35 * 793 = 27755

#######       #######
#.E...#       #.....#
#.#..G#       #.#G..#   G(200)
#.###.#  -->  #.###.#
#E#G#G#       #.#.#.#
#...#G#       #G.G#G#   G(98), G(38), G(200)
#######       #######

Combat ends after 54 full rounds
Goblins win with 536 total hit points left
Outcome: 54 * 536 = 28944

#########       #########
#G......#       #.G.....#   G(137)
#.E.#...#       #G.G#...#   G(200), G(200)
#..##..G#       #.G##...#   G(200)
#...##..#  -->  #...##..#
#...#...#       #.G.#...#   G(200)
#.G...G.#       #.......#


#.....G.#       #.......#
#########       #########

Combat ends after 20 full rounds
Goblins win with 937 total hit points left
Outcome: 20 * 937 = 18740

What is the outcome of the combat described in your puzzle input?

*/
