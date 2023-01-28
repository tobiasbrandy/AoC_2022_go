package day20

import (
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/linkedlist"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
)

func Solve(inputPath string, part int) any {
	const (
		bigStep   int = 1000
		stepCount int = 3
	)

	var mixFactor int64
	if part == 1 {
		mixFactor = 1
	} else { // part == 2
		mixFactor = 811_589_153
	}

	var mixCount int
	if part == 1 {
		mixCount = 1
	} else { // part == 2
		mixCount = 10
	}

	var nodes []*linkedlist.CircularDouble[int64]
	var zero *linkedlist.CircularDouble[int64]

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		data := int64(parse.Int(line)) * mixFactor
		l := len(nodes)

		if l == 0 {
			nodes = append(nodes, linkedlist.NewCircularDouble[int64](data))
		} else {
			nodes = append(nodes, nodes[l-1].Append(data))
		}

		if data == 0 {
			zero = nodes[l]
		}
	})

	nodeCount := len(nodes)

	for i := 0; i < mixCount; i++ {
		for _, node := range nodes {
			node.Move(int(node.Data % int64(nodeCount-1)))
		}
	}

	step := bigStep % nodeCount
	var total int64 = 0
	for i, curr := 0, zero; i < stepCount; i++ {
		curr = curr.Get(step)
		total += curr.Data
	}

	return total
}
