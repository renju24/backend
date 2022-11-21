package game

import (
	"fmt"
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

func moveFromStr(str string) Move {
	c := White
	r := rune(str[0])
	if unicode.IsUpper(r) {
		c = Black
		r = unicode.ToLower(r)
	}
	y := int(r - 'a')
	x, _ := strconv.Atoi(str[1:])
	x = 15 - x

	return Move{
		x:     x,
		y:     y,
		color: c,
	}
}

// initializes with specified sequence
func initGame(iniStr string) *Game {
	g := NewGame()
	for _, v := range reg.FindAllString(iniStr, -1) {
		g.lastMove = moveFromStr(v)
		g.board[g.lastMove.x*BoardSize+g.lastMove.y] = g.lastMove.color
	}
	return g
}

func printField(g *Game) {
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			var c rune
			switch g.board[i*BoardSize+j] {
			case Black:
				c = '*'
			case White:
				c = 'o'
			case Nil:
				c = '-'
			}
			fmt.Printf("%c", c)
		}
		fmt.Print("\n")
	}

	fmt.Print("\n")
}

func TestForks(t *testing.T) {
	testCases := []struct {
		iniStr        string
		move          Move
		expectedError error
	}{
		{
			iniStr:        "G9I9j9H8H7k4j3i2",
			move:          moveFromStr("H9"),
			expectedError: nil,
		},
		{
			iniStr:        "c12d11c10d9i14I12I11H8F6I6",
			move:          moveFromStr("I9"),
			expectedError: ErrInvalidForkForBlack,
		},
	}

	for _, testCase := range testCases {
		g := initGame(testCase.iniStr)
		printField(g)
		actualErr := g.checkForFork(testCase.move)
		require.ErrorIs(t, testCase.expectedError, actualErr)
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
				NewMove(5, 15, White),
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
