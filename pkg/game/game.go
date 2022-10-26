package game

import (
	"errors"
)

const BoardSize = 15

// Move structure.
type Move struct {
	x    int // X-coordinate.
	y    int // Y-coordinate.
	user int // User: 1 or 2 (black or white).
}

func NewMove(x, y, user int) Move {
	return Move{
		x:    x,
		y:    y,
		user: user,
	}
}

// Game structure.
type Game struct {
	board    [BoardSize * BoardSize]int // Board 15x15.
	lastMove Move                       // The last move.
}

func NewGame() *Game {
	return &Game{}
}

var (
	ErrCoordinatesOutside = errors.New("coordinates outside the board")
	ErrFieldAlreadyTaken  = errors.New("field is already taken")
	ErrInvalidTurn        = errors.New("invalid turn")
)

func (g *Game) ApplyMove(move Move) (winner int, err error) {
	// Check the coordinates are not outside the board.
	if move.x >= BoardSize || move.x < 0 || move.y >= BoardSize || move.y < 0 {
		return 0, ErrCoordinatesOutside
	}
	// Check the field is not already taken.
	if g.board[move.x*BoardSize+move.y] != 0 {
		return 0, ErrFieldAlreadyTaken
	}
	// Check the last move was made by another player.
	if g.lastMove.user == move.user {
		return 0, ErrInvalidTurn
	}

	// Apply the move and change the board.
	g.board[move.x*BoardSize+move.y] = move.user
	g.lastMove = move

	// After a successful move, we should check if there is a winner.
	if g.hasWinner() {
		return g.lastMove.user, nil
	}

	return 0, nil
}

func (game *Game) hasWinner() bool {
	var xcount, ycount, zcount int
	x, y := game.lastMove.x, game.lastMove.y
	x2, y2 := game.lastMove.x, game.lastMove.y
	if game.lastMove.user == 0 {
		return false
	}
	for i := 0; i < BoardSize; i++ {
		// "-"
		if xcount == 5 {
			return true
		}
		if game.board[x*BoardSize+i] == game.lastMove.user {
			xcount++
		} else {
			xcount = 0
		}
		// "|"
		if ycount == 5 {
			return true
		}
		if game.board[i*BoardSize+y] == game.lastMove.user {
			ycount++
		} else {
			ycount = 0
		}
	}
	// "/"
	for x2 > 0 && y2 < BoardSize {
		x2--
		y2++
	}
	for x2 < BoardSize && y2 > 0 {
		if zcount == 5 {
			return true
		}
		if game.board[x2*BoardSize+y2] == game.lastMove.user {
			zcount++
		} else {
			zcount = 0
		}
		x2++
		y2--
	}
	zcount = 0
	// "\"
	for x > 0 && y > 0 {
		x--
		y--
	}
	for x < BoardSize && y < BoardSize {
		if zcount == 5 {
			return true
		}
		if game.board[x*BoardSize+y] == game.lastMove.user {
			zcount++
		} else {
			zcount = 0
		}
		x++
		y++
	}
	return false
}
