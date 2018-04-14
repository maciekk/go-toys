package main

import (
	"fmt"
	"math/rand"
	"time"
)

// TODO: understand why/how "expressions can be implicitly repeated" in Go.
const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

type Room struct {
	wall map[int]bool
}

type Maze struct {
	rooms [][]Room
	w, h  int
}

func NewMaze(w, h int) Maze {
	var m Maze
	m.w, m.h = w, h
	m.rooms = make([][]Room, h)
	for i := 0; i < h; i++ {
		m.rooms[i] = make([]Room, w)
		for j := 0; j < w; j++ {
			m.rooms[i][j].wall = make(map[int]bool)
		}
	}
	return m
}

// NOTE: In the following methods we abstract away room access as still
// exploring best representation.

// Check whether room (x,y) has a wall in the given direction.
func (m *Maze) hasWall(x, y, dir int) bool {
	return m.rooms[y][x].wall[dir]
}

// Checks whether the room at (x,y) has all 4 walls.
func (m *Maze) walledIn(x, y int) bool {
	return m.hasWall(x, y, NORTH) &&
		m.hasWall(x, y, EAST) &&
		m.hasWall(x, y, SOUTH) &&
		m.hasWall(x, y, WEST)
}

// Modify a wall in the maze, as given by the (x, y, dir) parameters. 'State'
// controls whether we put in a wall or remove it.
func (m *Maze) setWall(x, y, dir int, state bool) {
	m.rooms[y][x].wall[dir] = state

	// We need to also update the /other/ instance of the wall, as our
	// representation is currently redundant.
	switch dir {
	case NORTH:
		if y < m.h-1 {
			m.rooms[y+1][x].wall[SOUTH] = state
		}
	case EAST:
		if x < m.w-1 {
			m.rooms[y][x+1].wall[WEST] = state
		}
	case SOUTH:
		if y > 0 {
			m.rooms[y-1][x].wall[NORTH] = state
		}
	case WEST:
		if x > 0 {
			m.rooms[y][x-1].wall[EAST] = state
		}
	}
}

// Construct a complete maze.
func (m *Maze) Build() {
	// First, put all walls in.
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			// Don't use m.setWall() here as would result in
			// dulpicate work.
			m.rooms[y][x].wall[NORTH] = true
			m.rooms[y][x].wall[EAST] = true
			m.rooms[y][x].wall[SOUTH] = true
			m.rooms[y][x].wall[WEST] = true
		}
	}
	// Now extrude paths; we have a number of methods for this, as
	// experimenting still for the most pleasing variant.
	//m.extrudeRandCellConnect()
	m.extrudeWalker()
}

// Find a "walled in" room anywhere in the maze.
func (m *Maze) findWalledIn() (x, y int, ok bool) {
	for x := 0; x < m.w; x++ {
		for y := 0; y < m.h; y++ {
			if m.walledIn(x, y) {
				return x, y, true
			}
		}
	}
	return 0, 0, false
}

// Find "explored" room that still has a "walled in" neighbour.
func (m *Maze) findWalledInNeighbour() (x, y int, ok bool) {
	for x := 0; x < m.w; x++ {
		for y := 0; y < m.h; y++ {
			if !m.walledIn(x, y) {
				dirs := m.dirsUnexplored(x, y)
				if len(dirs) > 0 {
					return x, y, true
				}
			}
		}
	}
	return 0, 0, false
}

// Return list of "legal" directions from room (x,y). Legal here means those
// that do not leave the bounds of the maze.
func (m *Maze) dirsLegal(x, y int) []int {
	dirs := make([]int, 0)
	if x > 0 {
		dirs = append(dirs, WEST)
	}
	if y > 0 {
		dirs = append(dirs, SOUTH)
	}
	if x < m.w-1 {
		dirs = append(dirs, EAST)
	}
	if y < m.h-1 {
		dirs = append(dirs, NORTH)
	}
	return dirs
}

// Return list of directions from (x,y) that reach rooms which are still
// "walled in".
func (m *Maze) dirsUnexplored(x, y int) []int {
	dirs := make([]int, 0)
	if x > 0 && m.walledIn(x-1, y) {
		dirs = append(dirs, WEST)
	}
	if y > 0 && m.walledIn(x, y-1) {
		dirs = append(dirs, SOUTH)
	}
	if x < m.w-1 && m.walledIn(x+1, y) {
		dirs = append(dirs, EAST)
	}
	if y < m.h-1 && m.walledIn(x, y+1) {
		dirs = append(dirs, NORTH)
	}
	return dirs
}

// Extrusion method: uses random walks, running each one until it gets stuck.
//
// NOTE: this produces "traditional" looking result, with a single, unique
// path between any pair of rooms. This is great for traditional maze solving
// puzzles, but suboptimal for say cat-and-mouse environment (future plans).
func (m *Maze) extrudeWalker() {
	// TODO: optimize this, it's very inefficient currently.

	// Establish a starting point.
	x, y := 0, 0 // TODO: start in random room?
	var ok bool

	for {
		// Walk around until get "stuck".
		for {
			// Which directions can we go in? Cannot leave bounds of
			// maze, and looking only to enter walled-in rooms.
			dirs := m.dirsUnexplored(x, y)
			if len(dirs) < 1 {
				// We are stuck!

				// Special case: stuck right from the start.
				// That is, we started a walk in a walled-in
				// room that is surrounded by already exlpored
				// space.
				if m.walledIn(x, y) {
					// Connect randomly.
					dirs := m.dirsLegal(x, y)
					dir := dirs[rand.Intn(len(dirs))]
					m.setWall(x, y, dir, false)
				}
				break
			}

			dir := dirs[rand.Intn(len(dirs))]
			m.setWall(x, y, dir, false)
			switch dir {
			case NORTH:
				y += 1
			case EAST:
				x += 1
			case SOUTH:
				y -= 1
			case WEST:
				x -= 1
			}
		}
		// Find a room to re-start the walker in.
		x, y, ok = m.findWalledInNeighbour()
		if !ok {
			// We are done.
			// This should mean that there are no unexplored rooms.
			return
		}

	}
}

// Extrusion method: iterate over all rooms in random order, and for each
// connect it to a random neighbour.
//
// Weakness: not very pleasing result, produces some "empty spaces".
func (m *Maze) extrudeRandCellConnect() {
	for _, cell := range rand.Perm(m.w * m.h) {
		// Unpack cell# to coordinates. Cell# starts at 0 at origin, then
		// counts across row, then up through the rows.
		y := cell / m.w
		x := cell % m.w
		if m.walledIn(x, y) {
			// Construct list of walls we COULD break down.
			dirs := make([]int, 0)
			if x > 0 {
				dirs = append(dirs, WEST)
			}
			if y > 0 {
				dirs = append(dirs, SOUTH)
			}
			if x < m.w-1 {
				dirs = append(dirs, EAST)
			}
			if y < m.h-1 {
				dirs = append(dirs, NORTH)
			}

			// Now break down ONE of these walls, randomly.
			dir := dirs[rand.Intn(len(dirs))]
			m.setWall(x, y, dir, false)
		}
	}
}

// Print the maze to tty.
func (m *Maze) Print() {
	line := "+"
	for x := 0; x < m.w; x++ {
		line += "---+"
	}
	fmt.Println(line)

	// Print such that origin is in lower left of printout.
	for y := m.h - 1; y >= 0; y-- {
		line1 := "|"
		line2 := "+"
		for x := 0; x < m.w; x++ {
			switch m.hasWall(x, y, EAST) {
			case true:
				line1 += "   |"
			default:
				line1 += "    "
			}
			switch m.hasWall(x, y, SOUTH) {
			case true:
				line2 += "---+"
			default:
				line2 += "   +"
			}
		}
		fmt.Println(line1)
		fmt.Println(line2)
	}
}

func main() {
	// Seed random number generator.
	rand.Seed(time.Now().UTC().UnixNano())

	m := NewMaze(20, 20)
	m.Build()
	m.Print()
}
