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
	// If it's the first move, then user should be black and move should be in board's center.
	if g.lastMove.color == Nil {
		if move.color != Black {
			return Nil, ErrFirstMoveShouldBeBlack
		}
		if move.x != 7 || move.y != 7 {
			return Nil, ErrFirstMoveShouldBeInCenter
		}
	}

	// Check the last move was made by another player.
	if g.lastMove.color == move.color {
		return Nil, ErrInvalidTurn
	}

	c, err := g.getColorAt(move.x, move.y)

	// Check the coordinates are not outside the board.
	if err != nil {
		return Nil, ErrCoordinatesOutside
	}
	// Check the field is not already taken.
	if c != Nil {
		return Nil, ErrFieldAlreadyTaken
	}

	lastMoveLength := g.maxRowAfterMove(move)

	if lastMoveLength > 5 && move.color == Black {
		return Nil, ErrRow6IsBannedForBlack
	}

	err = g.checkForFork(move)

	if err != nil {
		return Nil, err
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

func forkIsPermittedForColor(fork []int, c Color) bool {
	if c == Black {
		if len(fork) > 2 || (len(fork) == 2 && fork[0]*fork[1] != 12) { // if multiplicity > 2 or not 3x4 fork if multiplicity == 2
			return false
		}
	}
	return true
}

func indexIndiseBoard(i int) bool {
	if i >= 0 && i <= MaxBoardIndex {
		return true
	}
	return false
}

func (g *Game) checkForFork(m Move) error {
	g.setColorAt(m.x, m.y, m.color)
	// fork := []int{}

	startIndex := m.x*BoardSize + m.y

	for _, offset := range allDirectionOffsets {
		// freeSpace := 0
		rowLen := 1
		// openBorders := 0

		for i := -1; i <= 1; i += 2 { // i = -1, 1 to go on positive and negative directions
			for curIndex := startIndex + offset*i; indexIndiseBoard(curIndex); curIndex += offset * i {
				if g.board[curIndex] == m.color {
					rowLen += 1
				} else if g.board[curIndex] == Nil {

				} else {
					break // if next cell is opponent
				}
			}
		}
	}

	g.setColorAt(m.x, m.y, Nil)
	return nil
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
