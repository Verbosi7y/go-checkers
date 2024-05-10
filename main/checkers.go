package main

import (
	"fmt"
	"math"
)

type Coordinate struct {
	Row    uint8
	Column uint8
}

type Piece interface {
	GetColor() string
	GetPosition() Coordinate
	IsKing() bool
	GameRules
}

type Pawn struct {
	Color    string
	Position Coordinate
}

func (p *Pawn) GetColor() string { return p.Color }

func (p *Pawn) GetPosition() Coordinate { return p.Position }

func (p *Pawn) IsKing() bool { return false }

func (p Pawn) IsValidMove(to Coordinate, capture bool) bool {
	r_diff := float64(to.Row - p.GetPosition().Row)
	c_diff := float64(to.Column - p.GetPosition().Column)

	// range out of bounds
	if to.Row < 1 || to.Row > 8 || to.Column < 1 || to.Column > 8 {
		return false
	}

	// not diagonal
	if math.Abs(r_diff) != math.Abs(c_diff) {
		return false
	}

	// verify if black is going up
	if p.GetColor() == "black" && r_diff > 0 {
		return false
	}

	// verify if red is going down
	if p.GetColor() == "red" && r_diff < 0 {
		return false
	}

	// diagonal check if capturing
	if (capture == true) && (r_diff != 2 && r_diff != -2) {
		return false
	}

	// diagonal check if not capturing
	if (capture == false) && (r_diff != 1 && r_diff != -1) {
		return false
	}

	return true
}

func (p Pawn) IsValidCapture(b Board, to Coordinate) bool {
	// c >= 1, r >= 1 == right, upper
	// c >= 1, r <= -1 == right, bottom
	// c <= -1, r >= 1 == left, upper
	// c <= -1, r <= -1 == left, bottom
	r_diff := to.Row - p.GetPosition().Row
	c_diff := to.Column - p.GetPosition().Column

	if !p.IsValidMove(to, true) {
		return false
	}

	row := p.GetPosition().Row + (r_diff / 2)
	col := p.GetPosition().Column + (c_diff / 2)

	if (b.GetPiece(Coordinate{Row: row, Column: col}) != nil) {
		return false
	}

	return true
}

func (p Pawn) IsPromotable(toPosition Coordinate) bool {
	if toPosition.Row != 8 && toPosition.Row != 1 {
		return false
	}

	return true
}

type King struct {
	Color    string
	Position Coordinate
}

func (k *King) GetColor() string        { return k.Color }
func (k *King) GetPosition() Coordinate { return k.Position }
func (k *King) IsKing() bool            { return true }

func (k King) IsValidMove(to Coordinate, capture bool) bool {
	r_diff := float64(to.Row - k.GetPosition().Row)
	c_diff := float64(to.Column - k.GetPosition().Column)

	// range out of bounds
	if to.Row < 1 || to.Row > 8 || to.Column < 1 || to.Column > 8 {
		return false
	}

	// not diagonal
	if math.Abs(r_diff) != math.Abs(c_diff) {
		return false
	}

	return true
}

func (k King) IsValidCapture(b Board, to Coordinate) bool {
	// c >= 1, r >= 1 == right, upper
	// c >= 1, r <= -1 == right, bottom
	// c <= -1, r >= 1 == left, upper
	// c <= -1, r <= -1 == left, bottom
	r_diff := int(to.Row - k.GetPosition().Row)
	c_diff := int(to.Column - k.GetPosition().Column)

	row := k.GetPosition().Row
	col := k.GetPosition().Column

	if !k.IsValidMove(to, true) {
		return false
	}

	// verify to pos does not have any pieces
	if b.GetPiece(to) != nil {
		return false
	}

	// Error Case 1: Piece blocked to on upper-right
	// c > 1, r > 1 == right, upper
	if c_diff > 1 && r_diff > 1 {
		for c := 1; c < (c_diff - 1); c++ {
			col++
			row++
			if (b.GetPiece(Coordinate{Row: row, Column: col}) != nil) {
				return false
			}
		}
	}

	// Error Case 2: Piece blocked to on bottom-right
	// c > 1, r < -1 == right, bottom
	if c_diff > 1 && r_diff < -1 {
		for c := 1; c < (c_diff - 1); c++ {
			col++
			row--
			if (b.GetPiece(Coordinate{Row: row, Column: col}) != nil) {
				return false
			}
		}
	}

	// Error Case 3: Piece blocked to on upper-left
	// c <= -1, r >= 1 == left, upper
	if c_diff <= -1 && r_diff >= 1 {
		for c := 1; c < (r_diff - 1); c++ {
			col--
			row++
			if (b.GetPiece(Coordinate{Row: row, Column: col}) != nil) {
				return false
			}
		}
	}

	// Error Case 4: Piece blocked to on bottom-left
	// c <= -1, r <= -1 == left, bottom
	if c_diff <= -1 && r_diff <= -1 {
		for c := 1; c < ((c_diff * -1) - 1); c++ {
			col--
			row--
			if (b.GetPiece(Coordinate{Row: row, Column: col}) != nil) {
				return false
			}
		}
	}

	return true
}

func (k King) IsPromotable(toPosition Coordinate) bool { return false }

type Board struct {
	Pieces      []Piece
	Multiplayer bool
}

func (b Board) GetPiece(pos Coordinate) Piece {
	for _, p := range b.Pieces {
		if p.GetPosition() == pos {
			return p
		}
	}

	return nil
}

type GameRules interface {
	IsValidMove(to Coordinate, capture bool) bool // valid move pos -- check bounds, check piece type and if move is possible
	IsValidCapture(b Board, to Coordinate) bool   // valid capture -- check if there is a piece to capture and free space
	IsPromotable(toPosition Coordinate) bool      // valid promote -- check move pos
}

type GameEvaluation interface {
	allCaptured(c string) bool               // all captured check (c Color)
	anyLegalMoves(c string) bool             // any legal moves check (c Color)
	anySufficientMaterial(board *Board) bool // sufficient material check (board)
}

func (b *Board) Init(IsMultiplayer bool) {
	b.Pieces = []Piece{
		// BLACK
		&Pawn{Color: "black", Position: Coordinate{Row: 8, Column: 2}},
		&Pawn{Color: "black", Position: Coordinate{Row: 8, Column: 4}},
		&Pawn{Color: "black", Position: Coordinate{Row: 8, Column: 6}},
		&Pawn{Color: "black", Position: Coordinate{Row: 8, Column: 8}},

		&Pawn{Color: "black", Position: Coordinate{Row: 7, Column: 1}},
		&Pawn{Color: "black", Position: Coordinate{Row: 7, Column: 3}},
		&Pawn{Color: "black", Position: Coordinate{Row: 7, Column: 5}},
		&Pawn{Color: "black", Position: Coordinate{Row: 7, Column: 7}},

		&Pawn{Color: "black", Position: Coordinate{Row: 6, Column: 2}},
		&Pawn{Color: "black", Position: Coordinate{Row: 6, Column: 4}},
		&Pawn{Color: "black", Position: Coordinate{Row: 6, Column: 6}},
		&Pawn{Color: "black", Position: Coordinate{Row: 6, Column: 8}},

		// RED
		&Pawn{Color: "red", Position: Coordinate{Row: 3, Column: 1}},
		&Pawn{Color: "red", Position: Coordinate{Row: 3, Column: 3}},
		&Pawn{Color: "red", Position: Coordinate{Row: 3, Column: 5}},
		&Pawn{Color: "red", Position: Coordinate{Row: 3, Column: 7}},

		&Pawn{Color: "red", Position: Coordinate{Row: 2, Column: 2}},
		&Pawn{Color: "red", Position: Coordinate{Row: 2, Column: 4}},
		&Pawn{Color: "red", Position: Coordinate{Row: 2, Column: 6}},
		&Pawn{Color: "red", Position: Coordinate{Row: 2, Column: 8}},

		&Pawn{Color: "red", Position: Coordinate{Row: 1, Column: 1}},
		&Pawn{Color: "red", Position: Coordinate{Row: 1, Column: 3}},
		&Pawn{Color: "red", Position: Coordinate{Row: 1, Column: 5}},
		&Pawn{Color: "red", Position: Coordinate{Row: 1, Column: 7}}}

	b.Multiplayer = IsMultiplayer
}

func (b Board) PrintBoard() {
	var Reset = "\033[0m"
	var Black = "\033[30m"
	var Red = "\033[31m"

	fmt.Println("    Checker Board")
	fmt.Println(" +-----------------+")

	// print board
	for r := uint8(8); r >= 1; r-- {
		fmt.Printf("%d| ", r)

		for c := uint8(1); c <= 8; c++ {
			curr_piece := b.GetPiece(Coordinate{Row: r, Column: c})

			if (r+c)%2 == 1 { // red
				fmt.Print(Red)
			} else { // black
				fmt.Print(Black)
			}

			if curr_piece == nil {
				fmt.Printf("*%s ", Reset)
				continue
			}

			if curr_piece.GetColor() == "red" {
				fmt.Print(Red)
			} else {
				fmt.Print(Black)
			}

			if curr_piece.IsKing() {
				fmt.Print("K")
			} else {
				fmt.Print("P")
			}

			fmt.Printf("%s ", Reset)
		}
		fmt.Println("|")
	}

	fmt.Println(" +-----------------+")
	fmt.Println("   A B C D E F G H")
}

func (b *Board) Debugger() {
	fmt.Println("+=+=+-----Checkers Debugger-----+=+=+")
	b.PrintBoard()

	result := b.GetPiece(Coordinate{3, 1}).IsValidMove(Coordinate{4, 2}, false)

	fmt.Println(result)

}

func main() {
	Checkers := &Board{}
	Checkers.Init(false)

	Checkers.Debugger()
}
