package parse

import (
	"strconv"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
)

func Int(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		errexit.HandleMainError(err)
	}
	return n
}

