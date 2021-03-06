package main

import (
	. "../utils"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Progline struct {
	cmd      string
	operands []int
}

func main() {
	DEBUG = false

	infile := ReadFile("input")

	registers := make([]int, 6)

	addi := func(a, b, c int) {
		registers[c] = registers[a] + b
	}
	addr := func(a, b, c int) {
		addi(a, registers[b], c)
	}
	muli := func(a, b, c int) {
		registers[c] = registers[a] * b
	}
	mulr := func(a, b, c int) {
		muli(a, registers[b], c)
	}
	bani := func(a, b, c int) {
		registers[c] = registers[a] & b
	}
	banr := func(a, b, c int) {
		bani(a, registers[b], c)
	}
	bori := func(a, b, c int) {
		registers[c] = registers[a] | b
	}
	borr := func(a, b, c int) {
		bori(a, registers[b], c)
	}
	seti := func(a, b, c int) {
		registers[c] = a
	}
	setr := func(a, b, c int) {
		seti(registers[a], b, c)
	}

	// gteater-then-test
	gtTest := func(a, b, c int) {
		if a > b {
			registers[c] = 1
		} else {
			registers[c] = 0
		}

	}
	gtir := func(a, b, c int) {
		gtTest(a, registers[b], c)
	}
	gtri := func(a, b, c int) {
		gtTest(registers[a], b, c)
	}
	gtrr := func(a, b, c int) {
		gtTest(registers[a], registers[b], c)
	}

	// equals-test
	eqTest := func(a, b, c int) {
		if a == b {
			registers[c] = 1
		} else {
			registers[c] = 0
		}

	}
	eqir := func(a, b, c int) {
		eqTest(a, registers[b], c)
	}
	eqri := func(a, b, c int) {
		eqTest(registers[a], b, c)
	}
	eqrr := func(a, b, c int) {
		eqTest(registers[a], registers[b], c)
	}
	funcs := map[string]interface{}{"addi": addi, "addr": addr, "muli": muli, "mulr": mulr, "bani": bani,
		"banr": banr, "bori": bori, "borr": borr, "seti": seti, "setr": setr, "gtir": gtir, "gtri": gtri,
		"gtrr": gtrr, "eqir": eqir, "eqri": eqri, "eqrr": eqrr}

	ipBound, _ := strconv.Atoi(strings.Split(infile[0], " ")[1])
	if DEBUG {
		fmt.Println("IP bound to", ipBound)
	}

	var prog []Progline

	for _, line := range infile[1:] {
		opRaw := strings.Split(line, " ")
		operands := massiveAtoi(opRaw[1:])
		if DEBUG {
			fmt.Println(opRaw[0], operands)
		}
		prog = append(prog, Progline{opRaw[0], massiveAtoi(opRaw[1:])})
	}

	lastVal := -1
	var seen = make(map[int]bool)
	for ip := 0; ip < len(prog); ip++ {
		registers[ipBound] = ip
		op := prog[ip].operands
		fParam := make([]reflect.Value, 3)
		fParam[0] = reflect.ValueOf(op[0])
		fParam[1] = reflect.ValueOf(op[1])
		fParam[2] = reflect.ValueOf(op[2])
		reflect.ValueOf(funcs[prog[ip].cmd]).Call(fParam)
		if DEBUG {
			fmt.Println("regs after function", registers)

		}
		if ip == 29 { // eqrr 1 0 5, progline 29
			rv := registers[1]
			if lastVal == -1 {
				fmt.Println("Result1:", rv)
			}
			if seen[rv] {
				fmt.Println("Result2:", lastVal)
				break
			}
			seen[rv] = true
			lastVal = rv
		}
		ip = registers[ipBound]
	}

}

func massiveAtoi(in []string) []int {
	var result []int
	for _, i := range in {
		v, _ := strconv.Atoi(i)
		result = append(result, v)
	}
	return result
}

/*
--- Day 21: Chronal Conversion ---

You should have been watching where you were going, because as you wander the new North Pole base, you trip and fall into a very deep hole!

Just kidding. You're falling through time again.

If you keep up your current pace, you should have resolved all of the temporal anomalies by the next time the device activates. Since you have very little interest in browsing history in 500-year increments for the rest of your life, you need to find a way to get back to your present time.

After a little research, you discover two important facts about the behavior of the device:

First, you discover that the device is hard-wired to always send you back in time in 500-year increments. Changing this is probably not feasible.

Second, you discover the activation system (your puzzle input) for the time travel module. Currently, it appears to run forever without halting.

If you can cause the activation system to halt at a specific moment, maybe you can make the device send you so far back in time that you cause an integer underflow in time itself and wrap around back to your current time!

The device executes the program as specified in manual section one and manual section two.

Your goal is to figure out how the program works and cause it to halt. You can only control register 0; every other register begins at 0 as usual.

Because time travel is a dangerous activity, the activation system begins with a few instructions which verify that bitwise AND (via bani) does a numeric operation and not an operation as if the inputs were interpreted as strings. If the test fails, it enters an infinite loop re-running the test instead of allowing the program to execute normally. If the test passes, the program continues, and assumes that all other bitwise operations (banr, bori, and borr) also interpret their inputs as numbers. (Clearly, the Elves who wrote this system were worried that someone might introduce a bug while trying to emulate this system with a scripting language.)

What is the lowest non-negative integer value for register 0 that causes the program to halt after executing the fewest instructions? (Executing the same instruction multiple times counts as multiple instructions executed.)

--- Part Two ---

In order to determine the timing window for your underflow exploit, you also need an upper bound:

What is the lowest non-negative integer value for register 0 that causes the program to halt after executing the most instructions? (The program must actually halt; running forever does not count as halting.)
*/
