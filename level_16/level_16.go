package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Mapping struct {
	known, unknown int
}

func main() {
	filePart1, _ := os.Open("input1")
	defer filePart1.Close()
	filePart2, _ := os.Open("input2")
	defer filePart2.Close()

	registers := make([]int, 4)
	behaviour := make(map[Mapping]int, 256)

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
	funcs := []interface{}{addi, addr, muli, mulr, bani, banr, bori, borr, seti, setr, gtir, gtri, gtrr, eqir, eqri, eqrr}

	// part1
	scanner := bufio.NewScanner(filePart1)

	result1 := 0
	for scanner.Scan() {
		beforeRaw := strings.Split(strings.Trim(scanner.Text()[8:], "[]"), ", ")
		before := massiveAtoi(beforeRaw)
		scanner.Scan()
		opRaw := strings.Split(scanner.Text(), " ")
		op := massiveAtoi(opRaw)
		scanner.Scan()
		afterRaw := strings.Split(strings.Trim(scanner.Text()[8:], "[]"), ", ")
		after := massiveAtoi(afterRaw)
		scanner.Scan()
		fmt.Println(before, "=>", op, "=>", after)

		sampleFits := 0
		for k, f := range funcs {
			copy(registers, before)
			fmt.Println("regs before function", registers)
			val := reflect.ValueOf(f)
			fmt.Println("trying", k, "expected", after)
			fParam := make([]reflect.Value, 3)
			fParam[0] = reflect.ValueOf(op[1])
			fParam[1] = reflect.ValueOf(op[2])
			fParam[2] = reflect.ValueOf(op[3])
			val.Call(fParam)
			fmt.Println("regs after function", registers)
			rel := Mapping{k, op[0]}
			if behaviour[rel] < 0 {
				continue
			}
			if registers[0] == after[0] && registers[1] == after[1] && registers[2] == after[2] && registers[3] == after[3] {
				fmt.Println("Behaviour match")
				sampleFits++
			}
		}
		if sampleFits >= 3 {
			result1++
		}
		fmt.Println(behaviour)
	}

	fmt.Println("Result1", result1)
}

func massiveAtoi(in []string) []int {
	var result = []int{}
	for _, i := range in {
		v, _ := strconv.Atoi(i)
		result = append(result, v)
	}
	return result
}
