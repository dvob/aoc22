package main

import (
	"fmt"
	"testing"
)

func Test_getPos(t *testing.T) {
	data := []byte{
		'a', 'b', 'c',
		'd', 'e', 'f',
		'g', 'h', 'i',
	}
	f := field{
		data:    data,
		lineLen: 3,
	}

	for _, test := range []struct {
		pos         int
		dir         Direction
		expectedPos int
	}{
		{0, LEFT, -1},
		{0, UP, -1},
		{1, UP, -1},
		{1, DOWN, 4},
		{8, DOWN, -1},
		{8, RIGHT, -1},
		{7, DOWN, -1},
		{6, RIGHT, 7},
	} {
		t.Run(fmt.Sprintf("%d_%v", test.pos, test.dir), func(t *testing.T) {
			pos, _ := f.getPos(test.pos, test.dir)
			if pos != test.expectedPos {
				t.Fatalf("got=%d, want=%d", pos, test.expectedPos)
			}
		})
	}
}
