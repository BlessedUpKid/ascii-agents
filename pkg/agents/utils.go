package agents

import (
	"math/rand"
)

func wrapCoord(in, upperBound int) int {
	if in < 0 {
		return upperBound - 1
	}

	return in % upperBound
}

func clampHigh(in, upperBound int) int {
	if in == upperBound {
		return upperBound - 1
	}

	return in
}

// Based on an initial location generate the square of distance d from x, y
// The entire neighborhood of points can be listed by calling Ncoords for each
// square of points away from the origin
// I.e. a distance of 3 that includes everything below that would
// be Ncoords(x, y, 1) appended to Ncoords(x, y, 2) and Ncoords(x, y, 3)
func Ncoords(x, y, d int) [][]int {
	var coords [][]int

	ix, iy := (x - d), (y - d)
	l := ((x + d) - (x - d)) + 1

	// Top left corner to top right corner
	for i := 0; i < l; i++ {
		pair := []int{ix + i, iy}
		coords = append(coords, pair)
	}

	ix = coords[len(coords)-1][0]

	// One below the top right to the bottom right
	for i := 1; i < l; i++ {
		pair := []int{ix, iy + i}
		coords = append(coords, pair)
	}

	iy = coords[len(coords)-1][1]

	// bottom right to bottom left
	for i := 1; i < l; i++ {
		pair := []int{ix - i, iy}
		coords = append(coords, pair)
	}

	ix = coords[len(coords)-1][0]

	for i := 1; i < (l - 1); i++ {
		pair := []int{ix, iy - i}
		coords = append(coords, pair)
	}

	return coords
}

type agent struct {
	x          int
	y          int
	val        int
	randomizer *rand.Rand
	mgmt       *Agents
	dirX       int
	dirY       int
	p          float64
	g          int
}

func (a *agent) Lookhead() (int, int) {
	x := wrapCoord(a.x+a.dirX, a.mgmt.x)
	y := wrapCoord(a.y+a.dirY, a.mgmt.y)
	return x, y
}
