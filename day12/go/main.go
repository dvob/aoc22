package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

type field struct {
	data    []byte
	lineLen int
}

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

func (d Direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	default:
		panic("invalid direction")
	}
}

func (f *field) start() int {
	return bytes.IndexByte(f.data, 'S')
}

func (f *field) end() int {
	return bytes.IndexByte(f.data, 'E')
}

func (f *field) posToStr(i int) string {
	currentLine := i / f.lineLen
	linePos := i % f.lineLen
	return fmt.Sprintf("x=%d, y=%d, data=%c", linePos, currentLine, f.data[i])
}

func (f *field) getPos(pos int, dir Direction) (int, bool) {
	currentLine := pos / f.lineLen
	linePos := pos % f.lineLen
	switch dir {
	case UP:
		currentLine -= 1
	case DOWN:
		currentLine += 1
	case LEFT:
		linePos -= 1
	case RIGHT:
		linePos += 1
	default:
		panic("invalid direction")
	}
	if currentLine < 0 || currentLine >= (len(f.data)/f.lineLen) {
		return -1, false
	}
	if linePos < 0 || linePos >= f.lineLen {
		return -1, false
	}
	pos = currentLine*f.lineLen + linePos
	return pos, true
}

func (f *field) tooSteep(fromPos, toPos int) bool {
	from := f.data[fromPos]
	to := f.data[toPos]
	if to == 'S' {
		to = 'a'
	}
	if from == 'S' {
		from = 'a'
	}
	if to == 'E' {
		to = 'z'
	}
	if from == 'E' {
		from = 'z'
	}

	diff := int(to) - int(from)
	if diff > 1 {
		return true
	}
	return false
}

// get unvisited node with shortest distances
func getMinDistance(visited map[int]struct{}, distances map[int]int) int {
	min := int(^uint(0) >> 1)
	minPos := -1
	for pos, dist := range distances {
		if _, ok := visited[pos]; ok {
			continue
		}
		if dist < min {
			min = dist
			minPos = pos
		}
	}
	return minPos
}

func path(f *field, from, to int) (int, error) {
	visited := map[int]struct{}{}
	distance := map[int]int{}
	prev := map[int]int{}

	distance[from] = 0

	for {
		currentPos := getMinDistance(visited, distance)
		// no more nodes to visit
		if currentPos == -1 {
			break
		}
		currentDistance := distance[currentPos] + 1

		// fmt.Printf("visit %s\n", f.posToStr(currentPos))

		// check neighbors
		for _, direction := range []Direction{UP, DOWN, LEFT, RIGHT} {
			position, ok := f.getPos(currentPos, direction)

			// out of range
			if !ok {
				// fmt.Printf("  %s: out of range\n", direction)
				continue
			}

			if f.tooSteep(currentPos, position) {
				// fmt.Printf("  %s: too steep %s\n", direction, f.posToStr(position))
				continue
			}

			dist, ok := distance[position]
			if !ok || currentDistance < dist {
				distance[position] = currentDistance
				prev[position] = currentPos
				// fmt.Printf("  %s: set distance %s, prev=%s\n", direction, f.posToStr(position), f.posToStr(currentPos))
			} else {
				// fmt.Printf("  %s: has shorter dist %s\n", direction, f.posToStr(position))
			}
		}
		visited[currentPos] = struct{}{}
	}

	// print field
	// for i, c := range f.data {
	// 	if i > 0 && i%f.lineLen == 0 {
	// 		fmt.Println()
	// 	}
	// 	if distance, ok := distance[i]; ok {
	// 		fmt.Printf("%c(%.03d) ", c, distance)
	// 	} else {
	// 		fmt.Printf("%c(n/a) ", c)
	// 	}
	// }
	distanceTo, ok := distance[to]
	if !ok {
		return -1, fmt.Errorf("no path found")
	}
	return distanceTo, nil
}

func readInput(input io.Reader) (*field, error) {
	s := bufio.NewScanner(input)
	data := []byte{}
	lineLen := 0
	for s.Scan() {
		lineLen = len(s.Bytes())
		data = append(data, s.Bytes()...)
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return &field{
		data:    data,
		lineLen: lineLen,
	}, nil
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

	data, err := readInput(f)

	dist, _ := path(data, data.start(), data.end())
	fmt.Println(dist)

	// CPU goes brrrrrrr
	min := int(^uint(0) >> 1)
	for i, v := range data.data {
		if v != 'a' {
			continue
		}

		steps, err := path(data, i, data.end())
		if err != nil {
			continue
		}
		if steps < min {
			min = steps
		}
	}
	fmt.Println(min)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
