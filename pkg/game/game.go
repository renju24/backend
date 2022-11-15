package game

import (
	"errors"
)

const BoardSize = 15
const MaxBoardIndex = BoardSize*BoardSize - 1

// Move structure.
type Move struct {
	x     int   // X-coordinate.
	y     int   // Y-coordinate.
	color Color // color.
}

type Color int

const (
	Nil   Color = 0
	Black Color = 1
	White Color = 2
)

const (
	diagonalOffsetL  = BoardSize + 1 // "\"
	diagonalOffsetR  = BoardSize - 1 // "/"
	verticalOffset   = BoardSize     // "|"
	horizontalOffset = 1             // "-"
)

var allDirectionOffsets = []int{verticalOffset, horizontalOffset, diagonalOffsetL, diagonalOffsetR}

func NewMove(x, y int, color Color) Move {
	return Move{
		x:     x,
		y:     y,
		color: color,
	}
}

// Game structure.
type Game struct {
	board    [BoardSize * BoardSize]Color // Board 15x15.
	lastMove Move                         // The last move.
}

func NewGame() *Game {
	return &Game{}
}

var (
	ErrFirstMoveShouldBeBlack    = errors.New("first move should be made by black user")
	ErrFirstMoveShouldBeInCenter = errors.New("first move should be in board's center")
	ErrCoordinatesOutside        = errors.New("coordinates outside the board")
	ErrFieldAlreadyTaken         = errors.New("field is already taken")
	ErrInvalidTurn               = errors.New("invalid turn")
	ErrRow6IsBannedForBlack      = errors.New("black player cannot make row of length 6 and greater")
	ErrInvalidForkForBlack       = errors.New("black can make only 3x4 forks")
)

func (g *Game) getColorAt(x, y int) (Color, error) {
	if x >= BoardSize || x < 0 || y >= BoardSize || y < 0 {
		return Nil, ErrCoordinatesOutside
	}
	return g.board[x*BoardSize+y], nil
}

func (g *Game) setColorAt(x, y int, c Color) error {
	if x >= BoardSize || x < 0 || y >= BoardSize || y < 0 {
		return ErrCoordinatesOutside
	}
	g.board[x*BoardSize+y] = c
	return nil
}

func (g *Game) ApplyMove(move Move) (winner Color, err error) {
	e := g.checkMoveIsCorrect(move)
	if e != nil {
		return Nil, e
	}

	lastMoveLength := g.maxRowAfterMove(move)

	if lastMoveLength > 5 && move.color == Black {
		return Nil, ErrRow6IsBannedForBlack
	}

	// Apply the move and change the board.
	g.setColorAt(move.x, move.y, move.color)
	g.lastMove = move

	// After a successful move, we should check if there is a winner.
	if lastMoveLength >= 5 { //if row has length 5 or greater then player wins (if all rules above are passed)
		return g.lastMove.color, nil
	}

	return Nil, nil
}

func (g *Game) checkMoveIsCorrect(move Move) error {
	// If it's the first move, then user should be black and move should be in board's center.
	if g.lastMove.color == Nil {
		if move.color != Black {
			return ErrFirstMoveShouldBeBlack
		}
		if move.x != 7 || move.y != 7 {
			return ErrFirstMoveShouldBeInCenter
		}
	}

	c, err := g.getColorAt(move.x, move.y)

	// Check the coordinates are not outside the board.
	if err != nil {
		return ErrCoordinatesOutside
	}
	// Check the field is not already taken.
	if c != Nil {
		return ErrFieldAlreadyTaken
	}
	// Check the last move was made by another player.
	if g.lastMove.color == move.color {
		return ErrInvalidTurn
	}

	forks := countForksAfterMove(move)
	if len(forks) > 2 || (len(forks) == 2 && forks[0]*forks[1] != 12) {
		return ErrInvalidForkForBlack
	}

	return nil
}

func countForksAfterMove(m Move) []int {
	return []int{} //TODO
}

func (g *Game) maxRowAfterMove(move Move) int {
	startIndex := move.x*BoardSize + move.y
	maxLength := 1

	for _, offset := range allDirectionOffsets {
		curIndex := startIndex - offset
		curLength := 1
		for curIndex >= 0 && g.board[curIndex] == move.color {
			curLength++
			curIndex -= offset
		}
		curIndex = startIndex + offset
		for curIndex <= MaxBoardIndex && g.board[curIndex] == move.color {
			curLength++
			curIndex += offset
		}
		if curLength > maxLength {
			maxLength = curLength
		}
	}
	return maxLength
}
