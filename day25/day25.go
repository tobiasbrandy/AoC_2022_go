package day25

import (
	"fmt"
	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
)

func parseHotDigit(digit byte) int {
	switch digit {
	case '=':
		return -2
	case '-':
		return -1
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	default:
		errexit.HandleMainError(fmt.Errorf("invalid digit: %v", string(digit)))
		return 0
	}
}

func hotDigitString(digit int) byte {
	switch digit {
	case -2:
		return '='
	case -1:
		return '-'
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	default:
		errexit.HandleMainError(fmt.Errorf("invalid digit: %v", digit))
		return 0
	}
}

func parseHotInt(number string) int {
	ret := 0
	digits := []byte(number)
	digitsC := len(digits)
	for i := digitsC - 1; i >= 0; i-- {
		ret += parseHotDigit(digits[i]) * mathext.IntPow(5, uint(digitsC-i-1))
	}
	return ret
}

func hotIntString(n int) string {
	var digits []byte

	for n > 0 {
		d := n % 5
		carry := 0
		if d > 2 {
			carry = 1
			d -= 5
		}
		digits = append(digits, hotDigitString(d))
		n = n/5 + carry
	}

	// Reverse Digits
	digitsC := len(digits)
	for l, r := 0, digitsC-1; l < r; l, r = l+1, r-1 {
		digits[l], digits[r] = digits[r], digits[l]
	}

	return string(digits)
}

func Solve(inputPath string, _ int) any {
	total := 0
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		number := parseHotInt(line)
		total += number
		fmt.Println(number)
	})

	return hotIntString(total)
}
