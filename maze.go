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

// Abstract away room access as still exploring best representation.
func (m *Maze) hasWall(x, y, dir int) bool {
	return m.rooms[y][x].wall[dir]
}

func (m *Maze) walledIn(x, y int) bool {
	return m.hasWall(x, y, NORTH) &&
		m.hasWall(x, y, EAST) &&
		m.hasWall(x, y, SOUTH) &&
		m.hasWall(x, y, WEST)
}

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
	m.extrudeWanderers()
}

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

// Weakness: produces "components", pockets of connected rooms that are
// nonetheless unreachable from other parts of the maze.
func (m *Maze) extrudeWanderers() {
	// TODO: optimize this, it's very inefficient currently.
	for {
		x, y, ok := m.findWalledIn()
		if !ok {
			return // We're done.
		}

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
	}
}

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

	m := NewMaze(10, 10)
	m.Build()
	m.Print()
}
