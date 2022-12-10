package main

import (
	"testing"
)

func Test_moveRope(t *testing.T) {
	for _, test := range []struct {
		name     string
		r        rope
		dir      direction
		expected rope
	}{
		{
			name: "adjacent right",
			r: rope{
				head: position{},
				tail: position{},
			},
			dir: RIGHT,
			expected: rope{
				head: position{x: 1, y: 0},
				tail: position{},
			},
		},
		{
			name: "move tail right",
			r: rope{
				head: position{x: 1, y: 0},
				tail: position{},
			},
			dir: RIGHT,
			expected: rope{
				head: position{x: 2, y: 0},
				tail: position{x: 1, y: 0},
			},
		},
		{
			name: "move tail left",
			r: rope{
				head: position{x: 2, y: 0},
				tail: position{x: 3, y: 0},
			},
			dir: LEFT,
			expected: rope{
				head: position{x: 1, y: 0},
				tail: position{x: 2, y: 0},
			},
		},
		{
			name: "diagonal",
			r: rope{
				head: position{x: 1, y: 1},
				tail: position{},
			},
			dir: RIGHT,
			expected: rope{
				head: position{x: 2, y: 1},
				tail: position{x: 1, y: 1},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := moveRope(test.r, test.dir)
			if result != test.expected {
				t.Fatalf("got=%+v, want=%+v\n", result, test.expected)
			}
		})
	}
}
