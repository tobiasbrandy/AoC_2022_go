package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"flag"
	"fmt"
)

func buildTrees(filePath string) [][]int {
	var trees [][]int

	row := 0
	internal.ForEachFileLine(filePath, internal.HandleScanError, func(line string) {
		trees = append(trees, make([]int, len(line)))

		for col, num := range line {
			trees[row][col] = int(num - '0')
		}

		row++
	})

	return trees
}

func part1(filePath string) {
	trees := buildTrees(filePath)

	rows := len(trees)
	cols := len(trees[0])

	visible := make([][]bool, rows)
	for i := range visible {
		visible[i] = make([]bool, cols)
	}

	// Left to right
	for row := 0; row < rows; row++ {
		maxTree := -1
		for col := 0; col < cols; col++ {
			tree := trees[row][col]

			if tree > maxTree {
				visible[row][col] = true
				maxTree = tree
			}
		}
	}

		// Right to left
		for row := 0; row < rows; row++ {
			maxTree := -1
			for col := cols-1; col >= 0; col-- {
				tree := trees[row][col]

				if tree > maxTree {
					visible[row][col] = true
					maxTree = tree
				}
			}
		}

	// Up to down
	for col := 0; col < cols; col++ {
		maxTree := -1
		for row := 0; row < rows; row++ {
			tree := trees[row][col]

			if tree > maxTree {
				visible[row][col] = true
				maxTree = tree
			}
		}
	}

	// Down to up
	for col := 0; col < cols; col++ {
		maxTree := -1
		for row := rows-1; row >= 0; row-- {
			tree := trees[row][col]

			if tree > maxTree {
				visible[row][col] = true
				maxTree = tree
			}
		}
	}

	totalVisible := 0
	for row := range visible {
		for col := range visible[0] {
			if visible[row][col] {
				totalVisible++
			}
		}
	}

	fmt.Println(totalVisible)
}

func part2(filePath string) {
	trees := buildTrees(filePath)

	rows := len(trees)
	cols := len(trees[0])

	maxVisibility := 0

	for row := range trees {
		for col := range trees[0] {
			height := trees[row][col]
			
			right := 0
			for colRight := col+1; colRight < cols; colRight++ {
				right++
				if height <= trees[row][colRight] {
					break
				}
			}

			left := 0
			for colLeft := col-1; colLeft >= 0; colLeft-- {
				left++
				if height <= trees[row][colLeft] {
					break
				}
			}

			down := 0
			for rowDown := row+1; rowDown < rows; rowDown++ {
				down++
				if height <= trees[rowDown][col] {
					break
				}
			}

			up := 0
			for rowUp := row-1; rowUp >= 0; rowUp-- {
				up++
				if height <= trees[rowUp][col] {
					break
				}
			}

			visibility := right * left * down * up
			if visibility > maxVisibility {
				maxVisibility = visibility
			}
		}
	}

	fmt.Println(maxVisibility)
}

func main() {
	inputPath := flag.String("input", "input1.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	switch *part {
	case 1:
		part1(*inputPath)
	case 2:
		part2(*inputPath)
	default:
		internal.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}
}
