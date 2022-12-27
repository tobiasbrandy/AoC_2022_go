package main

import (
	"tobiasbrandy.com/aoc/2022/internal"

	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type FileSystem struct {
	parent *FileSystem
	dirs   map[string]*FileSystem
	files  map[string]int
}

func (fs *FileSystem) String() string {
	var sb strings.Builder
	_printFileSystemRec(&sb, map[string]*FileSystem{"/": fs}, 0)
	return sb.String()
}

func _printFileSystemRec(sb *strings.Builder, dirs map[string]*FileSystem, indent int) {
	for name, dir := range dirs {
		sb.WriteString(strings.Repeat(" ", 4*indent))
		sb.WriteString(fmt.Sprintf("- %v (dir)\n", name))

		for file, size := range dir.files {
			sb.WriteString(strings.Repeat(" ", 4*(indent+1)))
			sb.WriteString(fmt.Sprintf("- %v (file, size=%v)\n", file, size))
		}

		_printFileSystemRec(sb, dir.dirs, indent+1)
	}
}

func (fs *FileSystem) dirSizes(buf *[]int) int {
	total := 0

	for _, dir := range fs.dirs {
		total += dir.dirSizes(buf)
	}

	for _, fileSize := range fs.files {
		total += fileSize
	}

	if buf != nil {
		*buf = append(*buf, total)
	}
	return total
}

func execCommand(scanner *internal.FileLineScanner, root, cwd *FileSystem, cmd string) {
	if !strings.HasPrefix(cmd, "$ ") {
		internal.HandleMainError(fmt.Errorf("command %v doesn't start with '$ '", cmd))
	}
	cmd = cmd[len("$ "):]

	if strings.HasPrefix(cmd, "cd ") {
		execCd(scanner, root, cwd, cmd[len("cd "):])
	} else if strings.HasPrefix(cmd, "ls") {
		execLs(scanner, root, cwd)
	} else {
		internal.HandleMainError(fmt.Errorf("command %v doesn't exist", cmd))
	}
}

func execCd(scanner *internal.FileLineScanner, root, cwd *FileSystem, dir string) {
	cmd, ok := scanner.Read1()
	if !ok {
		// No more commands -> Nothing to do
		return
	}

	switch dir {
	case "/":
		execCommand(scanner, root, root, cmd)
	case "..":
		if cwd.parent == nil {
			// cwd is root
			execCommand(scanner, root, root, cmd)
		} else {
			execCommand(scanner, root, cwd.parent, cmd)
		}
	default:
		newCwd, ok := cwd.dirs[dir]
		if !ok {
			internal.HandleMainError(fmt.Errorf("cd %v: %v doesn't exist", dir, dir))
		}
		execCommand(scanner, root, newCwd, cmd)
	}
}

func execLs(scanner *internal.FileLineScanner, root, cwd *FileSystem) {
	hasNextCmd := false
	var nextCmd string

	scanner.ForEachWhile(func(line string) bool {
		if strings.HasPrefix(line, "$ ") {
			nextCmd = line
			hasNextCmd = true
			return false
		}

		if strings.HasPrefix(line, "dir ") {
			cwd.dirs[line[len("dir "):]] = &FileSystem{
				parent: cwd,
				dirs:   make(map[string]*FileSystem),
				files:  make(map[string]int),
			}
		} else { // file
			endIdx := strings.IndexRune(line, ' ')
			if endIdx == -1 {
				internal.HandleMainError(fmt.Errorf("%s: malformed ls command output", line))
			}

			size, err := strconv.Atoi(line[:endIdx])
			if err != nil {
				internal.HandleMainError(fmt.Errorf("%s: malformed ls command output: %v", line, err))
			}

			cwd.files[line[endIdx+1:]] = size
		}

		return true
	})

	if hasNextCmd {
		execCommand(scanner, root, cwd, nextCmd)
	}
}

func solve(filePath string, part int) {
	root := &FileSystem{
		parent: nil,
		dirs:   make(map[string]*FileSystem),
		files:  make(map[string]int),
	}

	scanner := internal.NewFileLineScanner(filePath, internal.HandleScanError)
	defer scanner.Close()

	cmd, ok := scanner.Read1()
	if !ok {
		internal.HandleMainError(errors.New("input is empty"))
	}
	execCommand(scanner, root, root, cmd)

	var dirSizes []int
	usedSpace := root.dirSizes(&dirSizes)

	if part == 1 {
		const maxDirSize int = 100_000

		total := 0
		for _, size := range dirSizes {
			if size < maxDirSize {
				total += size
			}
		}

		fmt.Println(total)
		return
	}

	// part == 2
	const totalSpace int = 70000000
	const targetSpace int = 30000000

	freeSpace := totalSpace - usedSpace
	missingSpace := targetSpace - freeSpace

	minDelete := usedSpace
	for _, dirSize := range dirSizes {
		if dirSize >= missingSpace && dirSize < minDelete {
			minDelete = dirSize
		}
	}

	fmt.Println(minDelete)
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
