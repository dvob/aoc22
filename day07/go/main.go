package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func sum(nums []int) int {
	result := 0
	for _, num := range nums {
		result += num
	}
	return result
}

func readDirs(input io.Reader) (map[string]int, error) {
	scanner := bufio.NewScanner(input)
	dirs := map[string][]int{}
	dir := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "$ cd ") {
			dir = filepath.Join(dir, line[len("$ cd "):])
			if _, ok := dirs[dir]; !ok {
				dirs[dir] = []int{}
			}
			continue
		}
		if strings.HasPrefix(line, "$ ls") {
			continue
		}
		if strings.HasPrefix(line, "dir ") {
			continue
		}

		// file
		strSize, _, found := strings.Cut(line, " ")
		if !found {
			return nil, fmt.Errorf("invalid line '%s'", line)
		}
		size, err := strconv.Atoi(strSize)
		if err != nil {
			return nil, fmt.Errorf("invalid line '%s': %w", line, err)
		}

		dirs[dir] = append(dirs[dir], size)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	dirsTotal := map[string]int{}
	for dir, files := range dirs {
		total := 0
		// find subdirs of dir
		for subDir, subDirFiles := range dirs {
			// is sub dir
			if len(subDir) > len(dir) && strings.HasPrefix(subDir, dir) {
				total += sum(subDirFiles)
			}
		}
		total += sum(files)

		dirsTotal[dir] = total
	}
	return dirsTotal, nil
}

func solve(dirs map[string]int) int {
	result := 0
	for _, size := range dirs {
		if size <= 100_000 {
			result += size
		}
	}
	return result
}

func solve2(dirs map[string]int) int {
	fsSpace := 70_000_000
	updateSize := 30_000_000
	used := dirs["/"]
	free := fsSpace - used
	needed := updateSize - free
	if needed < 0 {
		return 0
	}

	// max int isze
	minDir := int(^uint(0) >> 1)

	for _, size := range dirs {
		if size < needed {
			continue
		}
		if size < minDir {
			minDir = size
		}
	}
	return minDir
}

func run() error {
	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("missing argument: filename")
	}

	filename := flag.Arg(0)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	dirs, err := readDirs(f)
	if err != nil {
		return err
	}

	i := solve(dirs)
	fmt.Println(i)

	i = solve2(dirs)
	fmt.Println(i)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
