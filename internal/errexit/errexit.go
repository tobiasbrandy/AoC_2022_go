package errexit

import (
	"fmt"
	"os"
)

func HandleError(err error, exitCode int) {
	fmt.Println("Error:", err)
	os.Exit(exitCode)
}

func HandleArgsError(err error) {
	HandleError(err, 1)
}

func HandleScanError(err error) {
	HandleError(err, 2)
}

func HandleMainError(err error) {
	HandleError(err, 3)
}
