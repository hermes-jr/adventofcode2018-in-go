package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cart struct {
	id, x, y, velX, velY int
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
	fname = "input_test1"
	//fname = "input_test2"

	file, _ := os.Open(fname)
	defer file.Close()
	var data [][]byte

	scanner := bufio.NewScanner(file)

	carts := Carts{}
	vels := map[byte][]int{'>': {1, 0, '-'}, '<': {-1, 0, '-'}, '^': {0, -1, '|'}, 'v': {0, 1, '|'}}

	for rowNum, cartId := 0, 0; scanner.Scan(); rowNum++ {
		data = append(data, []byte{})
		inputLine := scanner.Bytes()
		for colNum, v := range inputLine {
			if v == '>' || v == '<' || v == 'v' || v == '^' {
				// parse cart and replace with underlying track if current cell is a cart
				data[rowNum] = append(data[rowNum], byte(vels[v][2]))
				carts = append(carts, Cart{cartId, colNum, rowNum, vels[v][0], vels[v][1]})
				cartId++
			} else {
				// parse cart and replace with underlying track if current cell is a cart
				data[rowNum] = append(data[rowNum], v)
			}
		}

	}
	printMap(data)

	fmt.Println(carts)
}

// Enumerate coordinates and print map (no carts yet)
func printMap(data [][]byte) {
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
			fmt.Printf("%v ", string(data[i][j]))
		}
		fmt.Println()
	}
}
