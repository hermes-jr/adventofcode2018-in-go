package main

import (
	"fmt"
)

const DEBUG = false

func main() {
	target := 607331
	//target = 5
	//target = 9
	//target = 18
	//target = 2018
	//target = 51589
	strg := []uint8{0, 1, 2, 4, 5}
	strg = []uint8{9, 2, 5, 1, 0}
	strg = []uint8{6, 0, 7, 3, 3, 1}

	//
	//51589 first appears after 9 recipes.
	//01245 first appears after 5 recipes.
	//92510 first appears after 18 recipes.
	//59414 first appears after 2018 recipes.

	scores := []uint8{3, 7}
	ptr1 := 0
	ptr2 := 1
	p1Progress := -1
	p2Progress := -1
	//strg := split(target)
	tlen := len(strg)
	for step := 0; p1Progress != 2 || p2Progress != 1; step++ {
		scorelen := len(scores)
		r2i := 0
		if DEBUG {
			fmt.Println("Searching ", strg)
			fmt.Println("in ", scores)
		}
		if p2Progress == -1 && scorelen > tlen {
			// calc part2
			fullMatch := true
			for i := 0; i < tlen; i++ {
				if strg[i] != scores[scorelen-tlen-1+i] {
					fullMatch = false
					break
				}
			}
			r2i = scorelen - tlen - 1
			if !fullMatch && scorelen > tlen+1 {
				for i := 0; i < tlen; i++ {
					if strg[i] != scores[scorelen-tlen-2+i] {
						fullMatch = false
						break
					}
				}
				r2i = scorelen - tlen - 2
			}
			if fullMatch {
				fmt.Println("Result2", r2i)
				p2Progress = 1
			}
		}

		if p1Progress == 1 {
			// print part1 result
			fmt.Print("Last ten: ")
			for nt := target; nt < target+10; nt++ {
				fmt.Print(scores[nt])
			}
			fmt.Println()
			p1Progress = 2
		}
		if p1Progress == -1 && scorelen+1 >= target+10 {
			// part1 condition satisfied, print results next time
			p1Progress = 1
		}

		ingredient1 := scores[ptr1]
		ingredient2 := scores[ptr2]
		newRecipe := ingredient1 + ingredient2
		if DEBUG {
			fmt.Printf("%5v: ", step)
			printScoreboard(scores, ptr1, ptr2)
		}

		scores = append(scores, split(int(newRecipe))...)
		scorelen = len(scores)
		ptr1 = nextid(ptr1, scorelen, scores[ptr1])
		ptr2 = nextid(ptr2, scorelen, scores[ptr2])

	}

}

func printScoreboard(scores []uint8, i1 int, i2 int) {
	for i := 0; i < len(scores); i++ {
		if i == i1 {
			fmt.Printf("(%v)", scores[i])
		} else if i == i2 {
			fmt.Printf("[%v]", scores[i])
		} else {
			fmt.Printf(" %v ", scores[i])
		}
	}
	fmt.Println()
}

func nextid(curpos, totalSize int, curval uint8) int {
	return (curpos + 1 + int(curval)) % totalSize
}

func split(in int) []uint8 {
	if in < 10 {
		return []uint8{uint8(in)}
	}
	var result []uint8
	for om := 1; in >= om; om *= 10 {
		result = append([]uint8{uint8((in / om) % 10)}, result...)
	}
	return result
}

/*
--- Day 14: Chocolate Charts ---

You finally have a chance to look at all of the produce moving around. Chocolate, cinnamon, mint, chili peppers, nutmeg, vanilla... the Elves must be growing these plants to make hot chocolate! As you realize this, you hear a conversation in the distance. When you go to investigate, you discover two Elves in what appears to be a makeshift underground kitchen/laboratory.

The Elves are trying to come up with the ultimate hot chocolate recipe; they're even maintaining a scoreboard which tracks the quality score (0-9) of each recipe.

Only two recipes are on the board: the first recipe got a score of 3, the second, 7. Each of the two Elves has a current recipe: the first Elf starts with the first recipe, and the second Elf starts with the second recipe.

To create new recipes, the two Elves combine their current recipes. This creates new recipes from the digits of the sum of the current recipes' scores. With the current recipes' scores of 3 and 7, their sum is 10, and so two new recipes would be created: the first with score 1 and the second with score 0. If the current recipes' scores were 2 and 3, the sum, 5, would only create one recipe (with a score of 5) with its single digit.

The new recipes are added to the end of the scoreboard in the order they are created. So, after the first round, the scoreboard is 3, 7, 1, 0.

After all new recipes are added to the scoreboard, each Elf picks a new current recipe. To do this, the Elf steps forward through the scoreboard a number of recipes equal to 1 plus the score of their current recipe. So, after the first round, the first Elf moves forward 1 + 3 = 4 times, while the second Elf moves forward 1 + 7 = 8 times. If they run out of recipes, they loop back around to the beginning. After the first round, both Elves happen to loop around until they land on the same recipe that they had in the beginning; in general, they will move to different recipes.

Drawing the first Elf as parentheses and the second Elf as square brackets, they continue this process:

(3)[7]
(3)[7] 1  0
 3  7  1 [0](1) 0
 3  7  1  0 [1] 0 (1)
(3) 7  1  0  1  0 [1] 2
 3  7  1  0 (1) 0  1  2 [4]
 3  7  1 [0] 1  0 (1) 2  4  5
 3  7  1  0 [1] 0  1  2 (4) 5  1
 3 (7) 1  0  1  0 [1] 2  4  5  1  5
 3  7  1  0  1  0  1  2 [4](5) 1  5  8
 3 (7) 1  0  1  0  1  2  4  5  1  5  8 [9]
 3  7  1  0  1  0  1 [2] 4 (5) 1  5  8  9  1  6
 3  7  1  0  1  0  1  2  4  5 [1] 5  8  9  1 (6) 7
 3  7  1  0 (1) 0  1  2  4  5  1  5 [8] 9  1  6  7  7
 3  7 [1] 0  1  0 (1) 2  4  5  1  5  8  9  1  6  7  7  9
 3  7  1  0 [1] 0  1  2 (4) 5  1  5  8  9  1  6  7  7  9  2

The Elves think their skill will improve after making a few recipes (your puzzle input). However, that could take ages; you can speed this up considerably by identifying the scores of the ten recipes after that. For example:

    If the Elves think their skill will improve after making 9 recipes, the scores of the ten recipes after the first nine on the scoreboard would be 5158916779 (highlighted in the last line of the diagram).
    After 5 recipes, the scores of the next ten would be 0124515891.
    After 18 recipes, the scores of the next ten would be 9251071085.
    After 2018 recipes, the scores of the next ten would be 5941429882.

What are the scores of the ten recipes immediately after the number of recipes in your puzzle input?

--- Part Two ---

As it turns out, you got the Elves' plan backwards. They actually want to know how many recipes appear on the scoreboard to the left of the first recipes whose scores are the digits from your puzzle input.

    51589 first appears after 9 recipes.
    01245 first appears after 5 recipes.
    92510 first appears after 18 recipes.
    59414 first appears after 2018 recipes.

How many recipes appear on the scoreboard to the left of the score sequence in your puzzle input?

*/
