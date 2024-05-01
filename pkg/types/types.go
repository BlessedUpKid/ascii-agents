package types

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

type IntMatrix struct {
	vals   [][]int
	colors [][]tcell.Color
	n      int
	m      int
}

func NewIntMatrix(n, m int) *IntMatrix {
	initInt := make([][]int, n)
	initColors := make([][]tcell.Color, n)
	for i := 0; i < n; i++ {
		initInt[i] = make([]int, m)
		initColors[i] = make([]tcell.Color, m)
	}

	return &IntMatrix{
		vals:   initInt,
		colors: initColors,
		n:      n,
		m:      m,
	}
}

func checkBounds[T int | tcell.Color](mat [][]T, x, y int) error {
	if x >= len(mat) {
		return errors.New("x out of bounds")
	}

	if y >= len(mat[0]) {
		return errors.New("y out of bounds")
	}

	return nil
}

func (m *IntMatrix) Cols() int {
	return m.n
}

func (m *IntMatrix) Rows() int {
	return m.m
}

func (m *IntMatrix) Set(x, y, val int) error {
	err := checkBounds(m.vals, x, y)
	if err != nil {
		return err
	}

	m.vals[x][y] = val
	return nil
}

func (m *IntMatrix) Get(x, y int) (int, error) {
	err := checkBounds(m.vals, x, y)
	if err != nil {
		return 0, err
	}
	return m.vals[x][y], nil
}

func (m *IntMatrix) Vals() [][]int {
	return m.vals
}

func (m *IntMatrix) Colors() [][]tcell.Color {
	return m.colors
}

func (m *IntMatrix) GetColor(x, y int) (tcell.Color, error) {
	err := checkBounds(m.colors, x, y)
	if err != nil {
		return 0, err
	}
	return m.colors[x][y], nil
}

func (m *IntMatrix) SetColor(x, y int, c tcell.Color) error {
	err := checkBounds(m.colors, x, y)
	if err != nil {
		return err
	}
	m.colors[x][y] = c
	return nil
}

func (m *IntMatrix) Clear() {
	// dumb clear for now
	vals := m.Vals()
	for i := range vals {
		for j := range vals[i] {
			m.Set(i, j, 0)
		}
	}
}

func (m *IntMatrix) Composite(n *IntMatrix, x, y int) {
	nvals := n.Vals()
	ncolors := n.Colors()

	for i := 0; i < len(nvals); i++ {
		for j := 0; j < len(nvals[0]); j++ {
			val, err := n.Get(i, j)
			if err != nil {
				continue
			}
			err = m.Set(x+i, y+j, val)
			if err != nil {
				continue
			}
		}
	}

	for i := 0; i < len(ncolors); i++ {
		for j := 0; j < len(ncolors[0]); j++ {
			tcolor, err := n.GetColor(i, j)
			if err != nil {
				continue
			}
			err = m.SetColor(x+i, y+j, tcolor)
			if err != nil {
				continue
			}
		}
	}
}
