package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGame(t *testing.T) {
	testCases := []struct {
		moves          []Move
		expectedWinner int
		expectedError  error
	}{
		// Case when the X-coordinate is outside the board.
		{
			moves: []Move{
				NewMove(-1, 5, 1),
			},
			expectedWinner: 0,
			expectedError:  ErrCoordinatesOutside,
		},
		// Case when the Y-coordinate is outside the board.
		{
			moves: []Move{
				NewMove(5, 15, 1),
			},
			expectedWinner: 0,
			expectedError:  ErrCoordinatesOutside,
		},
		// Case when field is already taken.
		{
			moves: []Move{
				NewMove(0, 0, 1),
				NewMove(0, 0, 2),
			},
			expectedWinner: 0,
			expectedError:  ErrFieldAlreadyTaken,
		},
		// Case when the last move was made by same player.
		{
			moves: []Move{
				NewMove(0, 0, 1),
				NewMove(1, 0, 1),
			},
			expectedWinner: 0,
			expectedError:  ErrInvalidTurn,
		},
		// Case when black should win.
		{
			moves: []Move{
				NewMove(0, 0, 1),
				NewMove(1, 0, 2),
				NewMove(0, 1, 1),
				NewMove(2, 0, 2),
				NewMove(0, 2, 1),
				NewMove(3, 0, 2),
				NewMove(0, 3, 1),
				NewMove(4, 0, 2),
				NewMove(0, 4, 1),
			},
			expectedWinner: 1,
			expectedError:  nil,
		},
		// Case when white should win.
		{
			moves: []Move{
				NewMove(0, 0, 1),
				NewMove(1, 0, 2),
				NewMove(0, 1, 1),
				NewMove(2, 0, 2),
				NewMove(0, 2, 1),
				NewMove(3, 0, 2),
				NewMove(0, 3, 1),
				NewMove(4, 0, 2),
				NewMove(0, 5, 1),
				NewMove(5, 0, 2),
			},
			expectedWinner: 2,
			expectedError:  nil,
		},
	}
	for _, testCase := range testCases {
		game := NewGame()
		for i, move := range testCase.moves {
			actualWinner, actualErr := game.ApplyMove(move)
			if i+1 == len(testCase.moves) {
				require.Equal(t, testCase.expectedWinner, actualWinner)
				require.ErrorIs(t, testCase.expectedError, actualErr)
			} else {
				require.Equal(t, 0, actualWinner)
				require.NoError(t, actualErr)
			}
		}
	}
}
