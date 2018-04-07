// Problem: place 8 queens on chessboard, such that none is threatened.

package main

import (
	"fmt"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func inverseString(s string) string {
	return "[7m" + s + "[0m"

}

type Square struct {
	File int
	Rank int
}

func (s *Square) IsValid() bool {
	if s.File < 0 || s.File > 7 || s.Rank < 0 || s.Rank > 7 {
		return false
	} else {
		return true
	}
}

func (s Square) String() string {
	return fmt.Sprintf("%c%d", 'A'+s.File, s.Rank+1)
}

type Chessboard [8][8]byte

func (b *Chessboard) Get(s Square) byte {
	if !s.IsValid() {
		return byte(' ')
	}
	if b[s.File][s.Rank] == 0 {
		return byte(' ')
	}
	return b[s.File][s.Rank]
}

func (b *Chessboard) Set(s Square, value byte) {
	if !s.IsValid() {
		return
	}
	b[s.File][s.Rank] = value
}

func squareColor(x, y int) int {
	return (x + y) % 2
}

func (b Chessboard) String() string {
	var s string = "\n"
	for r := 7; r >= 0; r-- {
		s += fmt.Sprintf("%d ", r+1)
		for f := 0; f < 8; f++ {
			piece := b.Get(Square{f, r})
			if squareColor(f, r) == 0 {
				s = s + " " + string(piece) + " "
			} else {
				s = s + inverseString(" "+string(piece)+" ")
			}
		}
		s += "\n"
	}
	s += "   A  B  C  D  E  F  G  H"
	return s
}

func MarkSquares(squares []Square, ch byte) Chessboard {
	var b Chessboard
	for _, s := range squares {
		b.Set(s, ch)
	}
	return b
}

func QueenThreatensSquare(sq, st Square) bool {
	if sq.File == st.File && sq.Rank == st.Rank {
		// No self-threats.
		return false
	}
	if sq.File == st.File || sq.Rank == st.Rank {
		// Rook move
		return true
	}
	file_delta := abs(sq.File - st.File)
	rank_delta := abs(sq.Rank - st.Rank)
	if file_delta == rank_delta {
		// Bishop move
		return true
	}
	return false
}

func SquaresThreatenedByQueenAt(s Square) []Square {
	var l []Square
	for f := 0; f < 8; f++ {
		for r := 0; r < 8; r++ {
			if f == s.File && r == s.Rank {
				// Do not count your own square.
				continue
			}
			s2 := Square{f, r}
			if QueenThreatensSquare(s, s2) {
				l = append(l, s2)
			}
		}
	}
	return l
}

func ShowQueens(queen_placements []Square) {
	var b Chessboard
	for _, s := range queen_placements {
		b.Set(s, 'Q')
	}
	fmt.Println(b)
}

func PlaceQueens() {
	PlaceQueenOnFile(0, []Square{})
}

var solutionCount int = 0

func PlaceQueenOnFile(file int, queens_so_far []Square) {
	if file >= 8 {
		// We went past last file; we are done.
		fmt.Println("\nFound:", queens_so_far)
		ShowQueens(queens_so_far)
		solutionCount++
		return
	}
	for rank := 0; rank < 8; rank++ {
		// Try each position.
		s := Square{file, rank}
		threatened := false
		for _, sq := range queens_so_far {
			if QueenThreatensSquare(sq, s) {
				threatened = true
				break
			}
		}
		if threatened {
			continue
		} else {
			// Viable placement; try next row.
			queens_so_far = append(queens_so_far, Square{file, rank})
			PlaceQueenOnFile(file+1, queens_so_far)
			queens_so_far = queens_so_far[:len(queens_so_far)-1]
		}
	}
}

func main() {
	fmt.Println("Sample square:", Square{1, 1})
	var b Chessboard
	b.Set(Square{0, 0}, 'Q')
	fmt.Printf("Piece at 0,0: %c\n", b.Get(Square{0, 0}))

	l := SquaresThreatenedByQueenAt(Square{0, 0})
	b2 := MarkSquares(l, 'x')
	fmt.Println("Squares threatened by Q at A1:", b2)
	l = SquaresThreatenedByQueenAt(Square{2, 1})
	b2 = MarkSquares(l, 'x')
	fmt.Println("Squares threatened by Q at C2:", b2)

	// Now attempt the 8 queen problem.
	PlaceQueens()
	//fmt.Println("Queens placed:", b3)

	fmt.Println()
	fmt.Printf("Found %d solutions.\n", solutionCount)
}
