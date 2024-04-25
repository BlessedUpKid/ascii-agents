package render

import (
	"viz/pkg/types"

	"github.com/gdamore/tcell/v2"
)

func int2ascii(i int) rune {
	// if i <= 0 {
	// 	return '.'
	// }

	if i > 0 && i < 250 {
		return '.'
	} else if i >= 250 && i < 500 {
		return 'O'
	} else if i >= 500 && i < 750 {
		return '0'
	} else if i >= 750 && i < 950 {
		return '@'
	}

	return '█'
}

func RenderMat(screen tcell.Screen, m *types.IntMatrix) {
	vals := m.Vals()

	for i := range vals {
		for j, n := range vals[i] {
			if n > 0 {
				c := m.GetColor(i, j)
				style := tcell.Style{}.
					Foreground(c)

				screen.SetContent(i, j, int2ascii(n), nil, style)
			}
		}
	}
}
