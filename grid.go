package ump

import (
	"math"
)

type grid struct {
	cellSize float32
	rows     map[int]map[int]*cell
}

func newGrid(cellSize int) *grid {
	return &grid{
		cellSize: float32(cellSize),
		rows:     make(map[int]map[int]*cell),
	}
}

func (g *grid) update(body *Body) {
	for _, c := range body.cells {
		c.leave(body)
	}
	body.cells = []*cell{}
	cl, ct, cw, ch := g.toCellRect(body.x, body.y, body.w, body.h)
	for cy := ct; cy <= ct+ch-1; cy++ {
		for cx := cl; cx <= cl+cw-1; cx++ {
			g.cellAt(float32(cx), float32(cy), true).enter(body)
		}
	}
}

func (g *grid) cellsInRect(l, t, w, h float32) []*cell {
	cl, ct, cw, ch := g.toCellRect(l, t, w, h)
	cells := []*cell{}
	for cy := ct; cy <= ct+ch-1; cy++ {
		row, ok := g.rows[cy]
		if ok {
			for cx := cl; cx <= cl+cw-1; cx++ {
				c, ok := row[cx]
				if ok {
					cells = append(cells, c)
				}
			}
		}
	}
	return cells
}

func (g *grid) toCellRect(x, y, w, h float32) (cx, cy, cw, ch int) {
	cx, cy = g.cellCoordsAt(x, y)
	cr, cb := int(math.Ceil(float64((x+w)/g.cellSize))), int(math.Ceil(float64((y+h)/g.cellSize)))
	return cx, cy, cr - cx, cb - cy
}

func (g *grid) cellCoordsAt(x, y float32) (cx, cy int) {
	return int(math.Floor(float64(x / g.cellSize))), int(math.Floor(float64(y / g.cellSize)))
}

func (g *grid) cellAt(x, y float32, cellCoords bool) *cell {
	var cx, cy int
	if cellCoords == true {
		cx, cy = int(x), int(y)
	} else {
		cx, cy = g.cellCoordsAt(x, y)
	}
	row, ok := g.rows[cy]
	if !ok {
		g.rows[cy] = make(map[int]*cell)
		row = g.rows[cy]
	}
	c, ok := row[cx]
	if !ok {
		row[cx] = &cell{bodies: make(map[uint32]*Body)}
		c = row[cx]
	}
	return c
}

func (g *grid) getCellsTouchedBySegment(x1, y1, x2, y2 float32) []*cell {
	cells := []*cell{}
	visited := map[*cell]bool{}

	g.traceRay(x1, y1, x2, y2, func(cx, cy int) {
		c := g.cellAt(float32(cx), float32(cy), true)
		if _, found := visited[c]; found {
			return
		}
		visited[c] = true
		cells = append(cells, c)
	})

	return cells
}

// traceRay* functions are based on "A Fast Voxel Traversal Algorithm for Ray Tracing",
// by John Amanides and Andrew Woo - http://www.cse.yorku.ca/~amana/research/grid.pdf
// It has been modified to include both cells when the ray "touches a grid corner",
// and with a different exit condition
func (g *grid) rayStep(ct, t1, t2 float32) (int, float32, float32) {
	v := t2 - t1
	delta := g.cellSize / v
	if v > 0 {
		return 1, delta, delta * (1.0 - frac(t1/g.cellSize))
	} else if v < 0 {
		return -1, -delta, -delta * frac(t1/g.cellSize)
	} else {
		return 0, inf, inf
	}
}

func (g *grid) traceRay(x1, y1, x2, y2 float32, f func(cx, cy int)) {
	cx1, cy1 := g.cellCoordsAt(x1, y1)
	cx2, cy2 := g.cellCoordsAt(x2, y2)
	stepX, dx, tx := g.rayStep(float32(cx1), x1, x2)
	stepY, dy, ty := g.rayStep(float32(cy1), y1, y2)
	cx, cy := cx1, cy1

	f(cx, cy)

	// The default implementation had an infinite loop problem when
	// approaching the last cell in some occassions. We finish iterating
	// when we are *next* to the last cell
	for abs(float32(cx-cx2))+abs(float32(cy-cy2)) > 1 {
		if tx < ty {
			tx += dx
			cx += stepX
		} else {
			// Addition: include both cells when going through corners
			if tx == ty {
				f(cx+stepX, cy)
			}
			ty += dy
			cy += stepY
		}
		f(cx, cy)
	}

	// If we have not arrived to the last cell, use it
	if cx != cx2 || cy != cy2 {
		f(cx2, cy2)
	}
}
