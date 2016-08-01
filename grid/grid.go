package grid

import (
	"fmt"
)

type TicTacToeGrid struct {
	grid [9]rune
}

func InitialGrid() *TicTacToeGrid {
	initialGrid := [9]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}

	return &TicTacToeGrid{grid: initialGrid}
}

func NewGrid(state [9]rune) *TicTacToeGrid {
	return &TicTacToeGrid{grid: state}
}

func (grid *TicTacToeGrid) String() string {
	return fmt.Sprintf("[[%c%c%c], [%c%c%c], [%c%c%c]]",
		grid.grid[0], grid.grid[1], grid.grid[2],
		grid.grid[3], grid.grid[4], grid.grid[5],
		grid.grid[6], grid.grid[7], grid.grid[8])
}

func (grid *TicTacToeGrid) MakeMove(position byte, player rune) {
	grid.grid[position] = player
}

// Each genome is a series of bytes that details the exact move to make for each possible state in the state space
// Each value describes the exact cell to place your piece (1-9)
func (grid *TicTacToeGrid) ToStateNumber() int {
	cellType := map[rune]int{
		' ': 0,
		'x': 1,
		'o': 2,
	}

	nums := [9]int{}
	for i := 0; i < 9; i++ {
		nums[i] = cellType[grid.grid[i]]
	}

	stateNum := nums[8]*1 + nums[7]*3 + nums[6]*9 + nums[5]*27 + nums[4]*81 + nums[3]*243 + nums[2]*729 + nums[1]*2187 + nums[0]*6561

	return stateNum
}

/*
  Given a Tic tac toe board return the winner if there is one otherwise return no winner.
  Return: int => 0 - NO WINNER, 1 - X WON, 2 - O WON
*/
func getWinnerStatus(cell rune) int {
	if cell == 'x' {
		return X_WON
	} else if cell == 'o' {
		return O_WON
	} else {
		return NO_WINNER
	}
}

const (
	NO_WINNER = iota
	X_WON
	O_WON
)

func (Grid *TicTacToeGrid) CheckWinner() int {
	grid := Grid.grid
	RowLength := 3
	for i := 0; i < RowLength; i++ {
		// Check rows
		if grid[i*RowLength] == grid[i*RowLength+1] && grid[i*RowLength+1] == grid[i*RowLength+2] {
			status := getWinnerStatus(grid[i*RowLength])
			if status != NO_WINNER {
				return status
			}
		}

		// Check columns
		if grid[i] == grid[i+RowLength] && grid[i+RowLength] == grid[i+2*RowLength] {
			status := getWinnerStatus(grid[i])
			if status != NO_WINNER {
				return status
			}
		}
	}

	// Check diagonals
	if (grid[0] == grid[4] && grid[4] == grid[8]) || (grid[6] == grid[4] && grid[4] == grid[2]) {
		status := getWinnerStatus(grid[4])
		if status != NO_WINNER {
			return status
		}
	}

	return NO_WINNER
}

func (Grid *TicTacToeGrid) IsFull() bool {
	for _, cell := range Grid.grid {
		if cell == ' ' {
			return false
		}
	}

	return true
}
