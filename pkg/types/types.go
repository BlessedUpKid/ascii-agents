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

func checkBounds(mat [][]int, x, y int) error {
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

func (m *IntMatrix) Get(x, y, val int) error {
	err := checkBounds(m.vals, x, y)
	if err != nil {
		return err
	}

	m.vals[x][y] = val
	return nil
}

func (m *IntMatrix) Vals() [][]int {
	return m.vals
}

func (m *IntMatrix) GetColor(x, y int) tcell.Color {
	return m.colors[x][y]
}

func (m *IntMatrix) SetColor(x, y int, c tcell.Color) {
	m.colors[x][y] = c
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
