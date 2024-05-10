package main

import "fmt"

type Coordinate struct {
	Row    uint8
	Column uint8
}

type Piece interface {
	GetColor() string
	GetPosition() Coordinate
	IsKing() bool
}

type Pawn struct {
	Color    string
	Position Coordinate
}

func (p *Pawn) GetColor() string        { return p.Color }
func (p *Pawn) GetPosition() Coordinate { return p.Position }
func (p *Pawn) IsKing() bool            { return false }

type King struct {
	Color    string
	Position Coordinate
}

func (k *King) GetColor() string        { return k.Color }
func (k *King) GetPosition() Coordinate { return k.Position }
func (k *King) IsKing() bool            { return true }

type Board struct {
	Pieces []Piece
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
	isValidMove(p Piece, toPosition Coordinate) bool    // valid move pos -- check bounds, check piece type and if move is possible
	isValidCapture(p Piece, toPosition Coordinate) bool // valid capture -- check if there is a piece to capture and free space
	promotable(p Piece, toPosition Coordinate) bool     // valid promote -- check move pos
}

type GameEvaluation interface {
	allCaptured(c string) bool               // all captured check (c Color)
	anyLegalMoves(c string) bool             // any legal moves check (c Color)
	anySufficientMaterial(board *Board) bool // sufficient material check (board)
}

func (b *Board) Init() {
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
}

func main() {
	Checkers := &Board{}
	Checkers.Init()

	Checkers.Debugger()
}
