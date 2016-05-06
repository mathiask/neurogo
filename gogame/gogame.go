package gogame

const (
	// board size
	n = 3

	// colors
	empty = 0
	black = 1
	white = 2
)

type color byte

// the board
type grid [n*n] color

// transform a board coordinate into a linear array index
func xy(x, y int) int {
	return n*y+x
}

// invert a color black<->white
func invert(c color) color {
	switch c {
	case black:
		return white
	case white:
		return black
	default:
		return empty
	}
}

// find neighbors of a point
func neighbors(xy int) []int {
	rc := make([]int, 0, 4)
	if xy%n>=1 {
		rc = append(rc, xy-1)
	}

	if xy%n+1!=n {
		rc = append(rc, xy+1)
	}

	if xy>=n {
		rc = append(rc, xy-n)
	}

	if xy+n<n*n {
		rc = append(rc, xy+n)
	}

	return rc
}

// find up to max liberties of a chain
func (g *grid) findLiberties(xy int, max int) int {
	libs := 0
	c := g[xy]
	opposite := invert(c)
	g[xy] = opposite // don't look here again

	// look at the neighbors
	for _, nxy := range neighbors(xy) {
		switch g[nxy] {
		case empty:
			// count liberty
			libs += 1
			g[nxy] = opposite // don't look here again

			// count up to max libs
			if libs>=max {
				return libs
			}
		case c:
			// recursively count liberties of neighbor
			libs += g.findLiberties(nxy, max-libs)

			// count up to max libs
			if libs>=max {
				return libs
			}
		}
	}

	// the liberties found so far
	return libs
}

// find up to max liberties of a chain
func (g *grid) liberties(xy int, max int) int {
	t := *g
	return t.findLiberties(xy, max)
}

// remove a chain of stones
func (g *grid) remove(xy int) *grid {
	c := g[xy]
	g[xy] = empty
	for _, nxy := range neighbors(xy) {
		if g[nxy]==c {
			g.remove(nxy)
		}
	}

	return g
}

// play a move
func (g *grid) mkmove(xy int, c color) *grid {
	if g[xy]!=empty {
		// don't play on non-empty points
		return nil
	}

	// play a move
	g[xy] = c

	// check neighbors
	for _, nxy := range neighbors(xy) {
		t := *g
		if t[nxy]==invert(c) && t.findLiberties(nxy, 1)==0 {
			// remove captured stones
			g.remove(nxy)
		}
	}

	// check liberties of the move played
	t := *g
	if t.findLiberties(xy, 1)==0 {
		// illegal move, no liberties
		return nil
	}

	return g
}
