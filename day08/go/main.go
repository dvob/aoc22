package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type direction int

const (
	up direction = iota
	right
	down
	left
)

type pos struct {
	x int
	y int
}

type iter struct {
	pos  pos
	dir  direction
	data [][]int
}

func newIter(pos pos, dir direction, data [][]int) *iter {
	return &iter{
		pos:  pos,
		dir:  dir,
		data: data,
	}
}

func (it *iter) next() (pos, int, bool) {
	if it.pos.x < 0 || it.pos.y < 0 || it.pos.y >= len(it.data) || it.pos.x >= len(it.data[it.pos.y]) {
		return pos{}, 0, false
	}
	pos := it.pos
	v := it.data[it.pos.y][it.pos.x]
	switch it.dir {
	case up:
		it.pos.y -= 1
	case right:
		it.pos.x += 1
	case down:
		it.pos.y += 1
	case left:
		it.pos.x -= 1
	}
	return pos, v, true
}

func visible(it *iter, positions map[pos]struct{}) {
	max := -1
	for {
		pos, v, ok := it.next()
		if !ok {
			break
		}
		if v > max {
			max = v
			positions[pos] = struct{}{}
		}
	}
}

func solve1(input [][]int) int {
	visiblePositions := map[pos]struct{}{}

	// rows
	for y := range input[0] {
		// from left to right
		it := newIter(pos{0, y}, right, input)
		visible(it, visiblePositions)

		// from right to left
		it = newIter(pos{len(input[y]) - 1, y}, left, input)
		visible(it, visiblePositions)
	}

	// columns
	for x := range input {
		// top down
		it := newIter(pos{x, 0}, down, input)
		visible(it, visiblePositions)

		// bottom up
		it = newIter(pos{x, len(input) - 1}, up, input)
		visible(it, visiblePositions)
	}

	return len(visiblePositions)
}

func dist(it *iter) int {
	_, start, ok := it.next()
	if !ok {
		panic("invalid position")
	}
	dist := 0
	for {
		_, v, ok := it.next()
		if !ok {
			break
		}
		dist++
		if v >= start {
			break
		}
	}

	return dist
}

func solve2(input [][]int) int {
	cur := 0
	for y := range input {
		for x := range input[0] {
			res := 1
			for _, dir := range []direction{up, right, down, left} {
				d := dist(newIter(pos{x, y}, dir, input))
				res *= d
			}
			if res > cur {
				cur = res
			}
		}
	}
	return cur
}

func parse(input io.Reader) ([][]int, error) {
	scanner := bufio.NewScanner(input)

	data := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, b := range scanner.Text() {
			row[i] = int(b - '0')
		}
		data = append(data, row)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return data, nil
}

func run() error {
	flag.Parse()

	if flag.NArg() < 1 {
		return fmt.Errorf("missing argument: filename")
	}

	filename := flag.Arg(0)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := parse(file)
	if err != nil {
		return err
	}

	result1 := solve1(data)
	fmt.Println(result1)

	result2 := solve2(data)
	fmt.Println(result2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
