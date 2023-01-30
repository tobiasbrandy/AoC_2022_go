package day21

import (
	"errors"
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/regext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/stringer"
	"regexp"
)

type MonkeyBinOp byte

func (op MonkeyBinOp) Execute(left, right int) int {
	switch op {
	case '+':
		return left + right
	case '-':
		return left - right
	case '*':
		return left * right
	case '/':
		return left / right
	default:
		errexit.HandleMainError(fmt.Errorf("invalid monkey binary operation: %v", byte(op)))
		return 0
	}
}

func (op MonkeyBinOp) Inv() MonkeyBinOp {
	switch op {
	case '+':
		return '-'
	case '-':
		return '+'
	case '*':
		return '/'
	case '/':
		return '*'
	default:
		errexit.HandleMainError(fmt.Errorf("invalid monkey binary operation: %v", byte(op)))
		return 0
	}
}

func (op MonkeyBinOp) SolveRight(result, left int) int {
	switch op {
	case '+':
		return result - left
	case '-':
		return left - result
	case '*':
		return result / left
	case '/':
		return left / result
	default:
		errexit.HandleMainError(fmt.Errorf("invalid monkey binary operation: %v", byte(op)))
		return 0
	}
}

func (op MonkeyBinOp) SolveLeft(result, left int) int {
	return op.Inv().Execute(result, left)
}

func (op MonkeyBinOp) String() string {
	return string(op)
}

type MonkeyRepo map[string]Monkey

type Monkey interface {
	// Number only works if no unknowns
	Number(repo MonkeyRepo) (int, bool)

	// SolveUnknown returns the number the unknown should have to reduce to target
	SolveUnknown(repo MonkeyRepo, result int) int
}

type LitMonkey int

func (m LitMonkey) Number(_ MonkeyRepo) (int, bool) {
	return int(m), true
}

func (m LitMonkey) SolveUnknown(_ MonkeyRepo, _ int) int {
	errexit.HandleMainError(errors.New("no unknown"))
	return 0
}

func (m LitMonkey) String() string {
	return stringer.String(m)
}

type BinOpMonkey struct {
	left, right string
	op          MonkeyBinOp
}

func (m *BinOpMonkey) Number(repo MonkeyRepo) (int, bool) {
	left, leftOk := repo[m.left].Number(repo)
	if leftOk {
		repo[m.left] = LitMonkey(left)
	}

	right, rightOk := repo[m.right].Number(repo)
	if rightOk {
		repo[m.right] = LitMonkey(right)
	}

	if !leftOk || !rightOk {
		return 0, false
	}

	return m.op.Execute(left, right), true
}

func (m *BinOpMonkey) SolveUnknown(repo MonkeyRepo, result int) int {
	left, leftOk := repo[m.left].Number(repo)
	right, rightOk := repo[m.right].Number(repo)

	if leftOk && rightOk {
		errexit.HandleMainError(errors.New("no unknown"))
		return 0
	}
	if !leftOk && !rightOk {
		errexit.HandleMainError(errors.New("more than one unknown"))
		return 0
	}

	if leftOk {
		fmt.Println(result, "=", left, m.op, "x", "=>", "x", "=", m.op.SolveRight(result, left))
		return repo[m.right].SolveUnknown(repo, m.op.SolveRight(result, left))
	} else { // rightOk
		fmt.Println(result, "=", "x", m.op, right, "=>", "x", "=", m.op.SolveLeft(result, right))
		return repo[m.left].SolveUnknown(repo, m.op.SolveLeft(result, right))
	}
}

func (m *BinOpMonkey) String() string {
	return stringer.String(m)
}

type UnknownMonkey struct{}

func (m UnknownMonkey) Number(_ MonkeyRepo) (int, bool) {
	return 0, false
}

func (m UnknownMonkey) SolveUnknown(_ MonkeyRepo, result int) int {
	return result
}

var monkeyParser = regexp.MustCompile(`(?P<name>[a-zA-Z]+): ((?P<lit>\d+)|(?P<left>[a-zA-Z]+) (?P<op>[+*/-]) (?P<right>[a-zA-Z]+))`)

func parseMonkey(data string) (name string, monkey Monkey) {
	monkeyInfo := regext.NamedCaptureGroups(monkeyParser, data)

	name = monkeyInfo["name"]
	if lit, ok := monkeyInfo["lit"]; ok {
		monkey = LitMonkey(parse.Int(lit))
	} else {
		monkey = &BinOpMonkey{
			left:  monkeyInfo["left"],
			right: monkeyInfo["right"],
			op:    MonkeyBinOp(monkeyInfo["op"][0]),
		}
	}

	return name, monkey
}

func Part1(inputPath string) any {
	const rootName string = "root"

	repo := make(MonkeyRepo)
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		name, monkey := parseMonkey(line)
		repo[name] = monkey
	})

	root := repo[rootName]
	ret, ok := root.Number(repo)
	if !ok {
		errexit.HandleMainError(errors.New("there was an unknown monkey expression in tree"))
	}

	return ret
}

func Part2(inputPath string) any {
	const (
		rootName string = "root"
		myName   string = "humn"
	)

	repo := make(MonkeyRepo)
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		name, monkey := parseMonkey(line)
		repo[name] = monkey
	})

	// We don't know owr value
	repo[myName] = UnknownMonkey{}

	// We know root is a bin op. We assign it the minus operation, so the result must be 0.
	root := repo[rootName].(*BinOpMonkey)
	root.op = '-'

	return root.SolveUnknown(repo, 0)
}
