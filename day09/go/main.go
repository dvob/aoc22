package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type direction int

const (
	UP direction = iota
	RIGHT
	DOWN
	LEFT
)

func strToDir(input string) (direction, error) {
	switch input {
	case "U":
		return UP, nil
	case "D":
		return DOWN, nil
	case "L":
		return LEFT, nil
	case "R":
		return RIGHT, nil
	default:
		return -1, fmt.Errorf("invalid direction '%s'", input)
	}
}

type command struct {
	dir  direction
	dist int
}

func strToCmd(input string) (command, error) {
	strDir, strDist, _ := strings.Cut(input, " ")
	dir, err := strToDir(strDir)
	if err != nil {
		return command{}, err
	}
	dist, err := strconv.Atoi(strDist)
	if err != nil {
		return command{}, err
	}
	return command{
		dir:  dir,
		dist: dist,
	}, nil
}

func parse(input io.Reader) ([]command, error) {
	cmds := []command{}

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		cmd, err := strToCmd(line)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return cmds, nil
}

type position struct {
	x int
	y int
}

func movePosition(pos position, dir direction) position {
	cmd := command{
		dir:  dir,
		dist: 1,
	}
	switch cmd.dir {
	case UP:
		pos.y += cmd.dist
	case RIGHT:
		pos.x += cmd.dist
	case DOWN:
		pos.y -= cmd.dist
	case LEFT:
		pos.x -= cmd.dist
	default:
		panic("invalid direction")
	}
	return pos
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func follow(head, tail position) position {
	vec := position{
		x: head.x - tail.x,
		y: head.y - tail.y,
	}
	// adjacent
	if abs(vec.x) <= 1 && abs(vec.y) <= 1 {
		return tail
	}
	//fmt.Printf("%+v", vec)

	if vec.x == 2 {
		vec.x = 1
	}
	if vec.x == -2 {
		vec.x = -1
	}
	if vec.y == 2 {
		vec.y = 1
	}
	if vec.y == -2 {
		vec.y = -1
	}
	tail.x = tail.x + vec.x
	tail.y = tail.y + vec.y
	return tail
}

func moveRope(pos []position, dir direction) {
	if len(pos) < 2 {
		panic("rope to short")
	}
	pos[0] = movePosition(pos[0], dir)
	for i := 1; i < len(pos); i++ {
		pos[i] = follow(pos[i-1], pos[i])
	}
}

func printVisited(h, w int, visited map[position]int) {
	for y := (h - 1); y >= 0; y-- {
		for x := 0; x < w; x++ {
			if _, ok := visited[position{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func solve(cmds []command, length int) int {
	visited := map[position]int{}
	rope := make([]position, length)

	visited[rope[len(rope)-1]]++

	for _, cmd := range cmds {
		for i := 0; i < cmd.dist; i++ {
			moveRope(rope, cmd.dir)
			visited[rope[len(rope)-1]]++
		}
	}
	return len(visited)
}

func run() error {
	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("missign argument: filename")
	}

	filename := flag.Arg(0)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	cmds, err := parse(f)
	if err != nil {
		return err
	}

	result1 := solve(cmds, 2)
	fmt.Println(result1)

	result2 := solve(cmds, 10)
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
