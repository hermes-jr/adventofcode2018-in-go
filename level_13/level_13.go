package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const DEBUG = false

type Cart struct {
	id, x, y  int
	turnCount int8
	direction byte
}

func (data Cart) String() string {
	return fmt.Sprintf("{%v %v:%v %v}", data.id, data.x, data.y, string(data.direction))
}

type Velocity struct {
	dx, dy int
}
type Repl struct {
	v    Velocity
	repl byte
}

type Carts []Cart

func (c Carts) Len() int      { return len(c) }
func (c Carts) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Carts) Less(i, j int) bool {
	if c[i].y < c[j].y {
		return true
	} else if c[i].y > c[j].y {
		return false
	} else {
		return c[i].x < c[j].x
	}
}

func main() {
	// must not remove trailing spaces in file
	fname := "input"
	//fname = "input_test1"
	//fname = "input_test2"
	//fname = "input_test3"

	file, _ := os.Open(fname)
	defer file.Close()
	var data [][]byte

	scanner := bufio.NewScanner(file)

	carts := Carts{}
	directionToVelocity := map[byte]Repl{
		'>': {Velocity{1, 0}, '-'},
		'<': {Velocity{-1, 0}, '-'},
		'^': {Velocity{0, -1}, '|'},
		'v': {Velocity{0, 1}, '|'},
		'x': {Velocity{0, 0}, 'x'}}

	for rowNum, cartId := 0, 0; scanner.Scan(); rowNum++ {
		data = append(data, []byte{})
		inputLine := scanner.Bytes()
		for colNum, v := range inputLine {
			if v == '>' || v == '<' || v == 'v' || v == '^' {
				// parse cart and replace with underlying track if current cell is a cart
				data[rowNum] = append(data[rowNum], byte(directionToVelocity[v].repl))
				carts = append(carts, Cart{cartId, colNum, rowNum, 0, v})
				cartId++
			} else {
				data[rowNum] = append(data[rowNum], v)
			}
		}

	}
	sort.Sort(carts)

	nCollisions := 0
	alive := len(carts)
outerLoop:
	for tick := 0; alive > 0; tick++ {
		if DEBUG {
			fmt.Println("Starting tick", tick, carts)
			printMap(data, carts)
		}

		for k, cart := range carts {
			if cart.direction == 'x' {
				continue // cart is already broken, next
			}

			cart = moveCart(cart, data, directionToVelocity[cart.direction].v)
			carts[k] = cart

			// detect collisions after moving
			// ineffective here
			for ik, possibleCollision := range carts {
				if possibleCollision.id == cart.id || possibleCollision.direction == 'x' {
					continue // can't collide with self, try next candidate for collision and debris are removed by elves
				}
				if cart.x == possibleCollision.x && cart.y == possibleCollision.y {
					// collision detected
					cart.direction = 'x'
					carts[k] = cart
					possibleCollision.direction = 'x'
					carts[ik] = possibleCollision
					nCollisions++
					alive -= 2
					if nCollisions == 1 {
						fmt.Printf("Result1: %v,%v\n", cart.x, cart.y)
					}
					break
				}
			}
		}
		if alive == 1 {
			for _, survivor := range carts {
				if survivor.direction != 'x' {
					fmt.Printf("Result2: %v,%v\n", survivor.x, survivor.y)
					break outerLoop
				}
			}
		}
		sort.Sort(carts) // resort carts to handle them correctly at next tick
		if DEBUG {
			fmt.Println("Tick complete", tick, carts)
			printMap(data, carts)
		}
	}
}

func moveCart(cart Cart, data [][]byte, v Velocity) Cart {
	if cart.direction == 'x' {
		return cart
	}
	cart.x += v.dx
	cart.y += v.dy
	x := cart.x
	y := cart.y
	railType := data[y][x]
	if railType == '+' {
		// handle intersection
		if cart.turnCount == 0 {
			turnLeft(&cart)
		} else if cart.turnCount == 2 {
			turnRight(&cart)
		}
		// straight otherwise
		cart.turnCount++
		if cart.turnCount == 3 {
			cart.turnCount = 0
		}
	} else if railType == '/' {
		if cart.direction == '>' || cart.direction == '<' {
			turnLeft(&cart)
		} else {
			turnRight(&cart)
		}
	} else if railType == '\\' {
		if cart.direction == '>' || cart.direction == '<' {
			turnRight(&cart)
		} else {
			turnLeft(&cart)
		}
	}
	return cart
}

func turnRight(cart *Cart) {
	switch cart.direction {
	case '>':
		cart.direction = 'v'
	case 'v':
		cart.direction = '<'
	case '<':
		cart.direction = '^'
	case '^':
		cart.direction = '>'
	default:
	}
}

func turnLeft(cart *Cart) {
	switch cart.direction {
	case '>':
		cart.direction = '^'
	case '^':
		cart.direction = '<'
	case '<':
		cart.direction = 'v'
	case 'v':
		cart.direction = '>'
	default:
	}
}

// Enumerate coordinates and print map (no carts yet)
func printMap(data [][]byte, carts Carts) {
	if len(data) == 0 {
		fmt.Println("-EMPTY MAP-")
		return
	}
	fmt.Printf("    ")
	for i := 0; i < len(data[0]); i++ {
		if i/100 != 0 {
			fmt.Printf(" %v", (i/100)%100)
		} else {
			fmt.Print("  ")
		}
	}
	fmt.Printf("\n    ")
	for i := 0; i < len(data[0]); i++ {
		if i/10 != 0 {
			fmt.Printf(" %v", (i/10)%10)
		} else {
			fmt.Print("  ")
		}
	}
	fmt.Printf("\n    ")
	for i := 0; i < len(data[0]); i++ {
		fmt.Printf(" %v", i%10)
	}
	fmt.Println()

	for i := 0; i < len(data); i++ {
		fmt.Printf("%3v: ", i)
		for j := 0; j < len(data[i]); j++ {
			ch := data[i][j]
			for _, ct := range carts {
				if ct.x == j && ct.y == i {
					ch = ct.direction
				}
			}
			fmt.Printf("%v ", string(ch))
		}
		fmt.Println()
	}
}

/*
--- Day 13: Mine Cart Madness ---

A crop of this size requires significant logistics to transport produce, soil, fertilizer, and so on. The Elves are very busy pushing things around in carts on some kind of rudimentary system of tracks they've come up with.

Seeing as how cart-and-track systems don't appear in recorded history for another 1000 years, the Elves seem to be making this up as they go along. They haven't even figured out how to avoid collisions yet.

You map out the tracks (your puzzle input) and see where you can help.

Tracks consist of straight paths (| and -), curves (/ and \), and intersections (+). Curves connect exactly two perpendicular pieces of track; for example, this is a closed loop:

/----\
|    |
|    |
\----/

Intersections occur when two perpendicular paths cross. At an intersection, a cart is capable of turning left, turning right, or continuing straight. Here are two loops connected by two intersections:

/-----\
|     |
|  /--+--\
|  |  |  |
\--+--/  |
   |     |
   \-----/

Several carts are also on the tracks. Carts always face either up (^), down (v), left (<), or right (>). (On your initial map, the track under each cart is a straight path matching the direction the cart is facing.)

Each time a cart has the option to turn (by arriving at any intersection), it turns left the first time, goes straight the second time, turns right the third time, and then repeats those directions starting again with left the fourth time, straight the fifth time, and so on. This process is independent of the particular intersection at which the cart has arrived - that is, the cart has no per-intersection memory.

Carts all move at the same speed; they take turns moving a single step at a time. They do this based on their current location: carts on the top row move first (acting from left to right), then carts on the second row move (again from left to right), then carts on the third row, and so on. Once each cart has moved one step, the process repeats; each of these loops is called a tick.

For example, suppose there are two carts on a straight track:

|  |  |  |  |
v  |  |  |  |
|  v  v  |  |
|  |  |  v  X
|  |  ^  ^  |
^  ^  |  |  |
|  |  |  |  |

First, the top cart moves. It is facing down (v), so it moves down one square. Second, the bottom cart moves. It is facing up (^), so it moves up one square. Because all carts have moved, the first tick ends. Then, the process repeats, starting with the first cart. The first cart moves down, then the second cart moves up - right into the first cart, colliding with it! (The location of the crash is marked with an X.) This ends the second and last tick.

Here is a longer example:

/->-\
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/

/-->\
|   |  /----\
| /-+--+-\  |
| | |  | |  |
\-+-/  \->--/
  \------/

/---v
|   |  /----\
| /-+--+-\  |
| | |  | |  |
\-+-/  \-+>-/
  \------/

/---\
|   v  /----\
| /-+--+-\  |
| | |  | |  |
\-+-/  \-+->/
  \------/

/---\
|   |  /----\
| /->--+-\  |
| | |  | |  |
\-+-/  \-+--^
  \------/

/---\
|   |  /----\
| /-+>-+-\  |
| | |  | |  ^
\-+-/  \-+--/
  \------/

/---\
|   |  /----\
| /-+->+-\  ^
| | |  | |  |
\-+-/  \-+--/
  \------/

/---\
|   |  /----<
| /-+-->-\  |
| | |  | |  |
\-+-/  \-+--/
  \------/

/---\
|   |  /---<\
| /-+--+>\  |
| | |  | |  |
\-+-/  \-+--/
  \------/

/---\
|   |  /--<-\
| /-+--+-v  |
| | |  | |  |
\-+-/  \-+--/
  \------/

/---\
|   |  /-<--\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/

/---\
|   |  /<---\
| /-+--+-\  |
| | |  | |  |
\-+-/  \-<--/
  \------/

/---\
|   |  v----\
| /-+--+-\  |
| | |  | |  |
\-+-/  \<+--/
  \------/

/---\
|   |  /----\
| /-+--v-\  |
| | |  | |  |
\-+-/  ^-+--/
  \------/

/---\
|   |  /----\
| /-+--+-\  |
| | |  X |  |
\-+-/  \-+--/
  \------/

After following their respective paths for a while, the carts eventually crash. To help prevent crashes, you'd like to know the location of the first crash. Locations are given in X,Y coordinates, where the furthest left column is X=0 and the furthest top row is Y=0:

           111
 0123456789012
0/---\
1|   |  /----\
2| /-+--+-\  |


3| | |  X |  |
4\-+-/  \-+--/
5  \------/

In this example, the location of the first crash is 7,3.

--- Part Two ---

There isn't much you can do to prevent crashes in this ridiculous system. However, by predicting the crashes, the Elves know where to be in advance and instantly remove the two crashing carts the moment any crash occurs.

They can proceed like this for a while, but eventually, they're going to run out of carts. It could be useful to figure out where the last cart that hasn't crashed will end up.

For example:

/>-<\
|   |
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/

/---\
|   |
| v-+-\
| | | |
\-+-/ |
  |   |
  ^---^

/---\
|   |
| /-+-\
| v | |
\-+-/ |
  ^   ^
  \---/

/---\
|   |
| /-+-\
| | | |
\-+-/ ^
  |   |
  \---/

After four very expensive crashes, a tick ends with only one cart remaining; its final location is 6,4.

What is the location of the last cart at the end of the first tick where it is the only cart left?

*/
