package game

import (
	"regexp"
	"strconv"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

//move is described as 2 chars -- letter as vertical coordinate and digit as horizontal
//upper-case letter is for black move and lower-case for white move
//i.e. D14 is black move, d14 is white move

var reg = regexp.MustCompile(`\w\d{1,2}`)

// initializes with specified sequence
func initGame(iniStr string) *Game {
	g := NewGame()
	for _, v := range reg.FindAllString(iniStr, -1) {
		c := White
		r := rune(v[0])
		if unicode.IsUpper(r) {
			c = Black
			r = unicode.ToLower(r)
		}
		x := int(r - 'a')
		y, _ := strconv.Atoi(v[1:])
		g.board[x*BoardSize+y] = c
		g.lastMove = Move{
			x:     x,
			y:     y,
			color: c,
		}
	}
	return g
}

func TestGame2(t *testing.T) {
	g := initGame("A0B1A2A3")

	for _, v := range g.board {
		t.Logf("%d", int(v))
	}
}
func TestGame(t *testing.T) {
	testCases := []struct {
		moves          []Move
		expectedWinner Color
		expectedError  error
	}{
		// Case when first move is not black.
		{
			moves: []Move{
				NewMove(1, 5, White),
			},
			expectedWinner: Nil,
			expectedError:  ErrFirstMoveShouldBeBlack,
		},
		// Case when first move is not in center.
		{
			moves: []Move{
				NewMove(1, 2, Black),
			},
			expectedWinner: Nil,
			expectedError:  ErrFirstMoveShouldBeInCenter,
		},
		// Case when the X-coordinate is outside the board.
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(-1, 5, White),
			},
			expectedWinner: Nil,
			expectedError:  ErrCoordinatesOutside,
		},
		// Case when the Y-coordinate is outside the board.
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(5, 15, Black),
			},
			expectedWinner: Nil,
			expectedError:  ErrCoordinatesOutside,
		},
		// Case when field is already taken.
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(7, 7, White),
			},
			expectedWinner: Nil,
			expectedError:  ErrFieldAlreadyTaken,
		},
		// Case when the last move was made by same player.
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(0, 0, White),
				NewMove(1, 0, White),
			},
			expectedWinner: Nil,
			expectedError:  ErrInvalidTurn,
		},
		// Case when black should win.
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(8, 8, White),
				NewMove(0, 0, Black),
				NewMove(1, 0, White),
				NewMove(0, 1, Black),
				NewMove(2, 0, White),
				NewMove(0, 2, Black),
				NewMove(3, 0, White),
				NewMove(0, 3, Black),
				NewMove(4, 0, White),
				NewMove(0, 4, Black),
			},
			expectedWinner: Black,
			expectedError:  nil,
		},
		// Case when white should win.
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(8, 8, White),
				NewMove(0, 0, Black),
				NewMove(1, 0, White),
				NewMove(0, 1, Black),
				NewMove(2, 0, White),
				NewMove(0, 2, Black),
				NewMove(3, 0, White),
				NewMove(0, 3, Black),
				NewMove(4, 0, White),
				NewMove(0, 5, Black),
				NewMove(5, 0, White),
			},
			expectedWinner: White,
			expectedError:  nil,
		},
		{
			moves: []Move{
				NewMove(7, 7, Black),
				NewMove(4, 4, White),
				NewMove(7, 6, Black),
				NewMove(4, 3, White),
				NewMove(7, 5, Black),
				NewMove(1, 1, White),
				NewMove(7, 9, Black),
				NewMove(1, 2, White),
				NewMove(7, 10, Black),
				NewMove(1, 3, White),
				NewMove(7, 8, Black),
			},
			expectedWinner: Nil,
			expectedError:  ErrRow6IsBannedForBlack,
		},
	}
	for testCaseNum, testCase := range testCases {
		game := NewGame()
		for i, move := range testCase.moves {
			actualWinner, actualErr := game.ApplyMove(move)
			if i+1 == len(testCase.moves) {
				require.Equal(t, testCase.expectedWinner, actualWinner)
				require.ErrorIs(t, testCase.expectedError, actualErr)
				t.Logf("testcase %d passed\n", testCaseNum)
			} else {
				require.Equal(t, Nil, actualWinner)
				require.NoError(t, actualErr)
			}
		}
	}
}
