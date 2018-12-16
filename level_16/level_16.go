package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const DEBUG = false

func main() {
	filePart1, _ := os.Open("input1")
	defer filePart1.Close()
	filePart2, _ := os.Open("input2")
	defer filePart2.Close()

	registers := make([]int, 4)

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

	var definitelyIncorrect [16][16]bool
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
		if DEBUG {
			fmt.Println(before, "=>", op, "=>", after)
		}

		sampleFits := 0
		for k, f := range funcs {
			copy(registers, before)
			if DEBUG {
				fmt.Println("regs before function", registers)
				fmt.Println("trying", k, "expected", after)
			}
			fParam := make([]reflect.Value, 3)
			fParam[0] = reflect.ValueOf(op[1])
			fParam[1] = reflect.ValueOf(op[2])
			fParam[2] = reflect.ValueOf(op[3])
			reflect.ValueOf(f).Call(fParam)
			if DEBUG {
				fmt.Println("regs after function", registers)
			}
			if registers[0] == after[0] && registers[1] == after[1] && registers[2] == after[2] && registers[3] == after[3] {
				// behaviour match
				sampleFits++
			} else {
				definitelyIncorrect[op[0]][k] = true
			}
		}
		if sampleFits >= 3 {
			result1++
		}
	}

	fmt.Println("Result1", result1)

	// part2
	if DEBUG {
		printDi(definitelyIncorrect)
	}
	identified := 0
	remap := make([]int, 16)
	for identified < 16 {
		for k, diRow := range definitelyIncorrect {
			defectsInRow := 0
			dj := 0
			for j := range diRow {
				if !(diRow[j]) {
					dj = j
					defectsInRow++
				}
			}
			if defectsInRow == 1 {
				identified++
				remap[k] = dj
				for si := 0; si < 16; si++ {
					definitelyIncorrect[si][dj] = true
				}
				if DEBUG {
					printDi(definitelyIncorrect)
				}
			}
		}
	}

	registers = []int{0, 0, 0, 0}
	scanner = bufio.NewScanner(filePart2)
	for scanner.Scan() {
		cmdRaw := strings.Split(scanner.Text(), " ")
		if DEBUG {
			fmt.Print("Executing", cmdRaw)
		}
		op := massiveAtoi(cmdRaw)
		f := funcs[remap[op[0]]]
		fParam := make([]reflect.Value, 3)
		fParam[0] = reflect.ValueOf(op[1])
		fParam[1] = reflect.ValueOf(op[2])
		fParam[2] = reflect.ValueOf(op[3])
		reflect.ValueOf(f).Call(fParam)
		if DEBUG {
			fmt.Println(" => ", registers)
		}
	}
	fmt.Println("Result2", registers[0])

}

func printDi(data [16][16]bool) {
	fmt.Print("   ")
	for i := 0; i < 16; i++ {
		fmt.Printf("%2v", i)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		fmt.Printf("%2v: ", i)
		for j := 0; j < 16; j++ {
			if data[j][i] == false {
				fmt.Print("o ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

func massiveAtoi(in []string) []int {
	var result = []int{}
	for _, i := range in {
		v, _ := strconv.Atoi(i)
		result = append(result, v)
	}
	return result
}

/*
--- Day 16: Chronal Classification ---

As you see the Elves defend their hot chocolate successfully, you go back to falling through time. This is going to become a problem.

If you're ever going to return to your own time, you need to understand how this device on your wrist works. You have a little while before you reach your next destination, and with a bit of trial and error, you manage to pull up a programming manual on the device's tiny screen.

According to the manual, the device has four registers (numbered 0 through 3) that can be manipulated by instructions containing one of 16 opcodes. The registers start with the value 0.

Every instruction consists of four values: an opcode, two inputs (named A and B), and an output (named C), in that order. The opcode specifies the behavior of the instruction and how the inputs are interpreted. The output, C, is always treated as a register.

In the opcode descriptions below, if something says "value A", it means to take the number given as A literally. (This is also called an "immediate" value.) If something says "register A", it means to use the number given as A to read from (or write to) the register with that number. So, if the opcode addi adds register A and value B, storing the result in register C, and the instruction addi 0 7 3 is encountered, it would add 7 to the value contained by register 0 and store the sum in register 3, never modifying registers 0, 1, or 2 in the process.

Many opcodes are similar except for how they interpret their arguments. The opcodes fall into seven general categories:

Addition:

    addr (add register) stores into register C the result of adding register A and register B.
    addi (add immediate) stores into register C the result of adding register A and value B.

Multiplication:

    mulr (multiply register) stores into register C the result of multiplying register A and register B.
    muli (multiply immediate) stores into register C the result of multiplying register A and value B.

Bitwise AND:

    banr (bitwise AND register) stores into register C the result of the bitwise AND of register A and register B.
    bani (bitwise AND immediate) stores into register C the result of the bitwise AND of register A and value B.

Bitwise OR:

    borr (bitwise OR register) stores into register C the result of the bitwise OR of register A and register B.
    bori (bitwise OR immediate) stores into register C the result of the bitwise OR of register A and value B.

Assignment:

    setr (set register) copies the contents of register A into register C. (Input B is ignored.)
    seti (set immediate) stores value A into register C. (Input B is ignored.)

Greater-than testing:

    gtir (greater-than immediate/register) sets register C to 1 if value A is greater than register B. Otherwise, register C is set to 0.
    gtri (greater-than register/immediate) sets register C to 1 if register A is greater than value B. Otherwise, register C is set to 0.
    gtrr (greater-than register/register) sets register C to 1 if register A is greater than register B. Otherwise, register C is set to 0.

Equality testing:

    eqir (equal immediate/register) sets register C to 1 if value A is equal to register B. Otherwise, register C is set to 0.
    eqri (equal register/immediate) sets register C to 1 if register A is equal to value B. Otherwise, register C is set to 0.
    eqrr (equal register/register) sets register C to 1 if register A is equal to register B. Otherwise, register C is set to 0.

Unfortunately, while the manual gives the name of each opcode, it doesn't seem to indicate the number. However, you can monitor the CPU to see the contents of the registers before and after instructions are executed to try to work them out. Each opcode has a number from 0 through 15, but the manual doesn't say which is which. For example, suppose you capture the following sample:

Before: [3, 2, 1, 1]
9 2 1 2
After:  [3, 2, 2, 1]

This sample shows the effect of the instruction 9 2 1 2 on the registers. Before the instruction is executed, register 0 has value 3, register 1 has value 2, and registers 2 and 3 have value 1. After the instruction is executed, register 2's value becomes 2.

The instruction itself, 9 2 1 2, means that opcode 9 was executed with A=2, B=1, and C=2. Opcode 9 could be any of the 16 opcodes listed above, but only three of them behave in a way that would cause the result shown in the sample:

    Opcode 9 could be mulr: register 2 (which has a value of 1) times register 1 (which has a value of 2) produces 2, which matches the value stored in the output register, register 2.
    Opcode 9 could be addi: register 2 (which has a value of 1) plus value 1 produces 2, which matches the value stored in the output register, register 2.
    Opcode 9 could be seti: value 2 matches the value stored in the output register, register 2; the number given for B is irrelevant.

None of the other opcodes produce the result captured in the sample. Because of this, the sample above behaves like three opcodes.

You collect many of these samples (the first section of your puzzle input). The manual also includes a small test program (the second section of your puzzle input) - you can ignore it for now.

Ignoring the opcode numbers, how many samples in your puzzle input behave like three or more opcodes?

--- Part Two ---

Using the samples you collected, work out the number of each opcode and execute the test program (the second section of your puzzle input).

What value is contained in register 0 after executing the test program?

*/
