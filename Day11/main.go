package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"flag"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gammazero/deque"
)

type Monkey struct {
	items       *deque.Deque[int]
	update      func(int) int
	throwTo     func(int) int
	divisor     int
	inspections int
}

var monkeyParser = regexp.MustCompile(
	`Monkey \d+:
  Starting items: (?P<items>(\d+, )*\d+)
  Operation: new = (?P<update>.+)
  Test: divisible by (?P<test>\d+)
    If true: throw to monkey (?P<throwTrue>\d+)
    If false: throw to monkey (?P<throwFalse>\d+)`)

var updateParser = regexp.MustCompile(`(?P<left>old|\d+) (?P<op>[+*/-]) (?P<right>old|\d+)`)

func parseUpdateInput(input string) func(int) int {
	if input == "old" {
		return func(old int) int { return old }
	}

	val := internal.ParseInt(input)
	return func(_ int) int { return val }
}

func parseMonkey(input string) *Monkey {
	// Parse monkey info
	monkeyInfo := internal.NamedCaptureGroups(monkeyParser, input)

	// Parse items
	itemsS := strings.Split(monkeyInfo["items"], ", ")
	items := deque.New[int](len(itemsS))
	for _, itemS := range itemsS {
		items.PushBack(internal.ParseInt(itemS))
	}

	// Parse update
	updateInfo := internal.NamedCaptureGroups(updateParser, monkeyInfo["update"])

	var updateOp func(int, int) int
	switch updateInfo["op"] {
	case "+":
		updateOp = func(l, r int) int { return l + r }
	case "-":
		updateOp = func(l, r int) int { return l - r }
	case "*":
		updateOp = func(l, r int) int { return l * r }
	case "/":
		updateOp = func(l, r int) int { return l / r }
	}

	updateL := parseUpdateInput(updateInfo["left"])
	updateR := parseUpdateInput(updateInfo["right"])

	update := func(old int) int { return updateOp(updateL(old), updateR(old)) }

	// Parse throwTo
	divisor := internal.ParseInt(monkeyInfo["test"])
	throwTrue := internal.ParseInt(monkeyInfo["throwTrue"])
	throwFalse := internal.ParseInt(monkeyInfo["throwFalse"])

	throwTo := func(item int) int {
		if item%divisor == 0 {
			return throwTrue
		}
		return throwFalse
	}

	return &Monkey{
		items:       items,
		update:      update,
		throwTo:     throwTo,
		divisor:     divisor,
		inspections: 0,
	}
}

func solve(filePath string, part int) {
	var rounds int
	if part == 1 {
		rounds = 20
	} else { // part == 2
		rounds = 10_000
	}

	var monkeys []*Monkey

	// Parse monkeys info
	internal.ForEachFileLineSet(filePath, internal.HandleScanError, func(lines []string) {
		monkeys = append(monkeys, parseMonkey(strings.Join(lines, "\n")))
	})

	reductionFactor := 1
	for _, m := range monkeys {
		reductionFactor *= m.divisor
	}

	// Process all rounds
	for i := 0; i < rounds; i++ {
		// printMonkeys(monkeys)

		for _, monkey := range monkeys {
			if monkey.items.Len() == 0 {
				continue
			}

			for monkey.items.Len() > 0 {
				monkey.inspections++

				item := monkey.items.PopFront()
				item = monkey.update(item)
				if part == 1 {
					item = item / 3
				} else { // part == 2
					// If we mod the item value by the least common multiple,
					// 	the divisibility remains the same and the actual number is reduced
					item = item % reductionFactor
				}
				newMonkey := monkeys[monkey.throwTo(item)]
				newMonkey.items.PushBack(item)
			}
		}
	}

	sort.Slice(monkeys, func(i, j int) bool { return monkeys[i].inspections > monkeys[j].inspections })

	business := uint64(monkeys[0].inspections) * uint64(monkeys[1].inspections)
	fmt.Println(business)
}

func main() {
	inputPath := flag.String("input", "input.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	if *part != 1 && *part != 2 {
		internal.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}

	solve(*inputPath, *part)
}
