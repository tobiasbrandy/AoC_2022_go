package day5

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"

	"github.com/gammazero/deque"
)

func parseStackLine(stacks []*deque.Deque[byte], line string) {
	stackCount := len(stacks)

	for i := 0; i < stackCount; i++ {
		crateName := line[(i+1)*4-3]
		if crateName != ' ' {
			stacks[i].PushBack(crateName)
		}
	}
}

var moveRegex = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

func Solve(inputPath string, part int) any {
	scanner := fileline.NewScanner(inputPath, errexit.HandleScanError)
	defer scanner.Close()

	initLine, ok := scanner.Read1()
	if !ok {
		errexit.HandleMainError(errors.New("input is empty"))
	}

	stackCount := (len(initLine) + 1) / 4

	stacks := make([]*deque.Deque[byte], stackCount)
	for i := 0; i < stackCount; i++ {
		stacks[i] = deque.New[byte]()
	}

	parseStackLine(stacks, initLine)
	scanner.ForEachWhile(func(line string) bool {
		if strings.HasPrefix(line, " 1") {
			return false
		}

		parseStackLine(stacks, line)
		return true
	})

	if line, ok := scanner.Read1(); !ok || line != "" {
		errexit.HandleMainError(errors.New("missing empty line dividing initial state from moves"))
	}

	var tmpStack *deque.Deque[byte]
	if part == 2 {
		tmpStack = deque.New[byte]()
	}

	scanner.ForEach(func(line string) {
		move := moveRegex.FindStringSubmatch(line)
		count := parse.Int(move[1])
		init := parse.Int(move[2])
		target := parse.Int(move[3])

		stackInit := stacks[init-1]
		stackTarget := stacks[target-1]

		if part == 1 {
			for i := 0; i < count; i++ {
				if stackInit.Len() == 0 {
					errexit.HandleMainError(fmt.Errorf("tried to move crate from empty stack %v", init))
				}

				stackTarget.PushFront(stackInit.PopFront())
			}
		} else { // part == 2
			for i := 0; i < count; i++ {
				if stackInit.Len() == 0 {
					errexit.HandleMainError(fmt.Errorf("tried to move crate from empty stack %v", init))
				}

				tmpStack.PushFront(stackInit.PopFront())
			}

			for tmpStack.Len() > 0 {
				stackTarget.PushFront(tmpStack.PopFront())
			}
		}
	})

	builder := strings.Builder{}
	for _, stack := range stacks {
		if stack.Len() > 0 {
			builder.WriteByte(stack.Front())
		}
	}

	return builder.String()
}
