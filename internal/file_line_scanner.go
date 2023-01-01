package internal

import (
	"bufio"
	"os"
)

func ForEachFileLine(filePath string, errHandler func(error), f func(string)) {
	scanner := NewFileLineScanner(filePath, errHandler)
	defer scanner.Close()

	scanner.ForEach(f)
}

func ForEachFileLineSet(filePath string, errHandler func(error), f func([]string)) {
	scanner := NewFileLineScanner(filePath, errHandler)
	defer scanner.Close()

	scanner.ForEachSet(f)
}

func ForEachFileLineSetN(filePath string, setLen int, errHandler func(error), f func([]string)) {
	scanner := NewFileLineScanner(filePath, errHandler)
	defer scanner.Close()

	scanner.ForEachSetN(setLen, f)
}

func ForEachFileLineSetWhile(filePath string, errHandler func(error), test func(line string) bool, f func([]string)) {
	scanner := NewFileLineScanner(filePath, errHandler)
	defer scanner.Close()

	scanner.ForEachSetWhile(test, f)
}

type FileLineScanner struct {
	file       *os.File
	scanner    *bufio.Scanner
	errHandler func(error)
}

func NewFileLineScanner(filePath string, errHandler func(error)) *FileLineScanner {
	file, err := os.Open(filePath)
	if err != nil {
		if errHandler != nil {
			errHandler(err)
		}
		return nil
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return &FileLineScanner{
		file:       file,
		scanner:    scanner,
		errHandler: errHandler,
	}
}

func (scanner *FileLineScanner) Close() {
	scanner.file.Close()
}

// Line functions

func (scanner *FileLineScanner) ForEach(f func(string)) {
	for scanner.scanner.Scan() {
		f(scanner.scanner.Text())
	}

	err := scanner.scanner.Err()
	if err != nil && scanner.errHandler != nil {
		scanner.errHandler(err)
	}
}

func (scanner *FileLineScanner) ForEachN(n int, f func(string)) int {
	if n <= 0 {
		return 0
	}

	total := n

	for ; n > 0 && scanner.scanner.Scan(); n-- {
		f(scanner.scanner.Text())
	}

	if n > 0 {
		err := scanner.scanner.Err()
		if err != nil && scanner.errHandler != nil {
			scanner.errHandler(err)
		}
	}

	return total - n
}

func (scanner *FileLineScanner) ReadN(n int, buf []string) int {
	i := 0
	return scanner.ForEachN(n, func(s string) {
		buf[i] = s
		i++
	})
}

func (scanner *FileLineScanner) Read1() (string, bool) {
	if scanner.scanner.Scan() {
		return scanner.scanner.Text(), true
	} else {
		err := scanner.scanner.Err()
		if err != nil && scanner.errHandler != nil {
			scanner.errHandler(err)
		}
		
		return "", false
	}
}

// Set Functions

func (scanner *FileLineScanner) ForEachWhile(f func(string) bool) {
	for scanner.scanner.Scan() && f(scanner.scanner.Text()) {
	}

	err := scanner.scanner.Err()
	if err != nil && scanner.errHandler != nil {
		scanner.errHandler(err)
	}
}

func (scanner *FileLineScanner) ForEachSetWhile(test func(line string) bool, f func([]string)) {
	input := make([]string, 5)

	eof := false
	for !eof {
		input = input[0:0]

		for {
			if !scanner.scanner.Scan() {
				eof = true
				break
			}

			s := scanner.scanner.Text()
			if test(s) {
				break
			}

			input = append(input, s)
		}
		
		f(input)
	}

	err := scanner.scanner.Err()
	if err != nil && scanner.errHandler != nil {
		scanner.errHandler(err)
	}
}

func (scanner *FileLineScanner) ForEachSet(f func([]string)) {
	scanner.ForEachSetWhile(func(line string) bool { return line == "" }, f)
}

func (scanner *FileLineScanner) ForEachSetN(setLen int, f func([]string)) {
	input := make([]string, setLen)

	for {
		n := scanner.ReadN(setLen, input)
		if n == 0 {
			break
		}

		f(input[:n])
	}
}
