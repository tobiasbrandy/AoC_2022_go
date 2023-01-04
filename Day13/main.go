package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"

	"golang.org/x/exp/constraints"
)

type Ord string

const (
	LT Ord = "LT"
	EQ Ord = "EQ"
	GT Ord = "GT"
)

func (ord Ord) Inv() Ord {
	switch ord {
	case LT:
		return GT
	case GT:
		return LT
	default:
		return ord
	}
}

func order[T constraints.Ordered](left, right T) Ord {
	switch {
	case left < right:
		return LT
	case left > right:
		return GT
	default:
		return EQ
	}
}

type Packet []byte

func (packet Packet) ParseInt() (n, len int) {
	// Only for small integers that fit int type.
	neg := false

	if packet[0] == '-' || packet[0] == '+' {
		neg = packet[0] == '-'
		packet = packet[1:]
		len++
	}

	for _, ch := range packet {
		if ch < '0' || ch > '9' {
			break
		}

		n = n*10 + int(ch-'0')
		len++
	}

	if neg {
		n = -n
	}

	return n, len
}

func (packet Packet) IsArr() bool {
	return packet[0] == '['
}

func (packet Packet) IsArrEnd() bool {
	return packet[0] == ']'
}

func packetCompare(left, right Packet) (ord Ord, newLeft, newRight Packet) {
	isLeftEnd, isRightEnd := left.IsArrEnd(), right.IsArrEnd()
	isLeftArr, isRightArr := left.IsArr(), right.IsArr()

	switch {
	case isLeftEnd || isRightEnd: // At least one array has ended
		switch {
		case isLeftEnd && isRightEnd:
			return EQ, left[1:], right[1:]
		case isLeftEnd:
			return LT, left[1:], right
		default: // isRightEnd:
			return GT, left, right[1:]
		}

	case !isLeftArr && !isRightArr: // Both are numbers
		nLeft, endLeft := left.ParseInt()
		nRight, endRight := right.ParseInt()
		return order(nLeft, nRight), left[endLeft:], right[endRight:]

	case isLeftArr && isRightArr: // Both are arrays
		for !left.IsArrEnd() && !right.IsArrEnd() {
			left, right = left[1:], right[1:]
			ord, left, right = packetCompare(left, right)
			if ord != EQ {
				return ord, left, right
			}
		}

		return packetCompare(left, right)

	default: // One array and one number
		var arr Packet
		var num Packet
		if isLeftArr {
			arr, num = left, right
		} else {
			arr, num = right, left
		}

		ord, arr, num := packetCompare(arr[1:], num)
		if ord != EQ {
			if isLeftArr {
				return ord, arr, num
			} else {
				return ord.Inv(), num, arr
			}
		}

		if !arr.IsArrEnd() {
			if isLeftArr {
				return GT, arr, num
			} else {
				return LT, num, arr
			}
		}

		arr = arr[1:]
		if isLeftArr {
			return EQ, arr, num
		} else {
			return EQ, num, arr
		}
	}
}

func (left Packet) Less(right Packet) bool {
	ord, _, _ := packetCompare(left, right)
	return ord == LT
}

func part1(filePath string) {
	total := 0

	index := 1
	fileline.ForEachSet(filePath, errexit.HandleScanError, func(lines []string) {
		if len(lines) != 2 {
			errexit.HandleMainError(fmt.Errorf("should only be 2 inputs between each empty line, not %v", len(lines)))
		}

		left, right := Packet(lines[0]), Packet(lines[1])
		if left.Less(right) {
			total += index
		}

		index++
	})

	fmt.Println(total)
}

func part2(filePath string) {
	div1 := &Packet{'[', '[', '2', ']', ']'}
	div2 := &Packet{'[', '[', '6', ']', ']'}

	packets := []*Packet{div1, div2}

	fileline.ForEach(filePath, errexit.HandleScanError, func(line string) {
		if line == "" {
			return
		}

		packet := Packet(line)
		packets = append(packets, &packet)
	})

	sort.Slice(packets, func(i, j int) bool { return (*packets[i]).Less(*packets[j]) })

	div1Idx, div2Idx := -1, -1
	for i, p := range packets {
		if p == div1 {
			div1Idx = i + 1
		} else if p == div2 {
			div2Idx = i + 1
		}

		if div1Idx != -1 && div2Idx != -1 {
			break
		}
	}

	total := div1Idx * div2Idx
	fmt.Println(total)
}

func main() {
	inputPath := flag.String("input", "input.txt", "Path to the input file")
	part := flag.Int("part", 1, "Part number of the AoC challenge")

	flag.Parse()

	switch *part {
	case 1:
		part1(*inputPath)
	case 2:
		part2(*inputPath)
	default:
		errexit.HandleArgsError(fmt.Errorf("no part %v exists in challenge", *part))
	}
}
