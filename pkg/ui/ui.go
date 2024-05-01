package ui

import "viz/pkg/types"

func DrawRect(m *types.IntMatrix, x, y, u, w int) {
	// x -> u
	// y -> w
	// u -> x
	// w -> y

	// --->
	//^   |
	//<-- v

	// (1, 1) ---> (4, 1)
	// -> 3

	lx := (u - x) + 1 // [0, 1, 2, 3]
	ly := (w - y) + 1

	xi := x
	yi := y

	for i := 0; i < lx; i++ {
		m.Set(x+i, y, 2000)
	}

	xi = x + (lx - 1)

	for i := 1; i < ly; i++ {
		m.Set(xi, y+i, 2000)
	}

	yi = y + (ly - 1)

	// i = 1, 2, 3
	// 2, 1, 0,
	for i := 1; i < lx; i++ {
		m.Set(xi-i, yi, 2000)
	}
	xi = xi - (lx - 1)

	for i := 1; i < (ly - 1); i++ {
		m.Set(xi, yi-i, 2000)
	}
}

func DrawWindow() {}
