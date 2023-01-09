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

const (
	rowState = iota
	spaceState
)

var leftRightDir = [2]int{-1, 1}

func coordinatesByInex(indx int) (x, y int) {
	x = indx / BoardSize
	y = indx - x*BoardSize
	return x, y
}

func nextIndex(curIndex, offset int) (int, error) {
	next := curIndex + offset
	var err error
	if (next < 0) || (next > MaxBoardIndex) || (offset*offset == 1 && curIndex/BoardSize != next/BoardSize) {
		err = errors.New("outside of field")
	}
	return next, err
}

type row struct {
	centerLen int
	sideLen   []int
	gapIdx    []int
	endIdx    []int
}

func newRow() row {
	return row{
		sideLen: make([]int, 0),
		gapIdx:  make([]int, 0),
		endIdx:  make([]int, 0),
	}
}

const (
	FOUR = iota
	OPEN_FOUR
	NONE
)

func (g *Game) checkFour(r row, m Move) int {
	if len(r.gapIdx) == 1 { // if 2 sections
		x, y := coordinatesByInex(r.gapIdx[0])
		if g.maxRowAfterMove(NewMove(x, y, m.color)) <= 5 {
			return FOUR
		}
	} else { // one section
		cnt := 0
		for _, v := range r.endIdx {
			x, y := coordinatesByInex(v)
			if g.maxRowAfterMove(NewMove(x, y, m.color)) <= 5 {
				cnt += 1
			}
		}
		if cnt == 2 {
			return OPEN_FOUR
		}
	}
	return NONE
}

func (g *Game) checkForFork(m Move) error {
	fork := []int{}
	startIndex := m.x*BoardSize + m.y

	for _, offset := range allDirectionOffsets {
		r := newRow()
		r.centerLen = 1

		for _, dir := range [2]int{-1, 1} { // i = -1, 1 to go on positive and negative directions
			prevLen := 0
			state := rowState
			finished := false
			rowNum := 0
			curIndex := startIndex

			for {
				var err error
				curIndex, err = nextIndex(curIndex, offset*dir)
				if err != nil {
					switch state {
					case rowState:
						if rowNum == 0 {
							r.centerLen += prevLen
						} else {
							r.sideLen = append(r.sideLen, prevLen)
						}
					}
					break
				}

				switch state {
				case rowState:
					if g.board[curIndex] == m.color {
						prevLen += 1
					} else {
						if rowNum == 0 {
							r.centerLen += prevLen
						} else {
							r.sideLen = append(r.sideLen, prevLen)
						}
						if g.board[curIndex] == Nil {
							state = spaceState
						} else {
							finished = true
						}
					}
				case spaceState:
					if g.board[curIndex] == Nil {
						r.endIdx = append(r.endIdx, curIndex-dir*offset)
						finished = true
					} else {
						if g.board[curIndex] == m.color && rowNum == 0 {
							r.gapIdx = append(r.gapIdx, curIndex-dir*offset)
							rowNum += 1
							prevLen = 1
							state = rowState
						} else {
							r.endIdx = append(r.endIdx, curIndex-dir*offset)
							finished = true
						}
					}
				}

				if finished {
					break
				}
			}
		}

		if len(r.endIdx) < 2 { //blocked on one or both sides
			continue
		}

		if len(r.sideLen) == 2 { // case when we have 3 sections in a row
			if r.centerLen+r.sideLen[0] == 4 && r.centerLen+r.sideLen[1] == 4 {
				cnt := 0
				for v := range r.gapIdx {
					x, y := coordinatesByInex(v)
					if g.maxRowAfterMove(NewMove(x, y, m.color)) <= 5 {
						cnt += 1
					}
				}
				if cnt == 2 {
					fork = append(fork, 4, 4)
				}
			}
			continue
		}

		totalLen := r.centerLen
		for _, v := range r.sideLen {
			totalLen += v
		}

		switch totalLen {
		case 3:
			g.board[startIndex] = m.color
			if len(r.gapIdx) == 1 {
				x, y := coordinatesByInex(r.gapIdx[0])
				r.gapIdx = nil
				tmpMove := NewMove(x, y, m.color)
				err := g.checkForFork(tmpMove)
				if err == nil && g.checkFour(r, tmpMove) == OPEN_FOUR {
					fork = append(fork, 3)
				}
			} else {
				for i, v := range r.endIdx {
					stopCheck := false
					x, y := coordinatesByInex(v)
					tmpMove := NewMove(x, y, m.color)
					err := g.checkForFork(tmpMove)
					if err == nil {
						g.board[v] = m.color
						r.endIdx[i] += []int{-1, 1}[i] * offset
						x, y := coordinatesByInex(r.endIdx[i])
						iMove := NewMove(x, y, m.color)
						if g.checkFour(r, iMove) == OPEN_FOUR {
							fork = append(fork, 3)
							stopCheck = true
						}
						r.endIdx[i] -= []int{-1, 1}[i] * offset
						g.board[v] = Nil
					}
					if stopCheck {
						break
					}
				}
			}
			g.board[startIndex] = Nil
		case 4:
			fourType := g.checkFour(r, m)

			if fourType == FOUR || fourType == OPEN_FOUR {
				fork = append(fork, 4)
			}
		}

	}

	if !forkIsPermittedForColor(fork, m.color) {
		return ErrInvalidForkForBlack
	} else {
		return nil
	}
}

func (g *Game) maxRowAfterMove(m Move) int {
	startIndex := m.x*BoardSize + m.y
	maxLength := 1
	for _, offset := range allDirectionOffsets {
		curLen := 1
		for _, dir := range [2]int{-1, 1} {
			curIndex := startIndex
			for {
				var err error
				curIndex, err = nextIndex(curIndex, offset*dir)
				if err != nil || g.board[curIndex] != m.color {
					break
				}
				curLen += 1
			}
		}
		if curLen > maxLength {
			maxLength = curLen
		}
	}
	return maxLength
}
