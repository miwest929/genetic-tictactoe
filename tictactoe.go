package main

import (
	"bytes"
	"fmt"
	"github.com/genetic/grid"
	"math/rand"
	//	"strconv"
	"time"
)

type Evolver interface {
	GenerateRandom() []byte
	Fitness([]byte, []byte) float64
}

type TicTacToeSolver struct {
}

// two moves can snuggly fit in one byte
//TODO: Avoid returning illegal moves. For example, player can't make a move in a position that is
//      already occupied.
func getRandomTwoMoves() byte {
	var firstOctet uint8 = uint8(rand.Float64() * 9)
	var secondOctet uint8 = uint8(rand.Float64() * 9)
	var newByte byte = (firstOctet << 4) | secondOctet

	return newByte
}

func GetMove(strategy []byte, stateNum int) byte {
	stateIndex := 0
	if stateNum%2 == 0 {
		stateIndex = stateNum / 2
	} else {
		stateIndex = (stateNum - 1) / 2
	}

	gene := strategy[stateIndex]

	var moveOctet byte
	if stateNum%2 == 0 {
		moveOctet = gene >> 4
	} else {
		moveOctet = gene & 0x0f
	}

	//TODO: Remove debug print statements
	fmt.Println("--------------------------------")
	fmt.Printf("stateNum = %d, ", stateNum)
	fmt.Printf("stateIndex = %d, ", stateIndex)
	fmt.Printf("moveByte = %b, ", gene)
	fmt.Printf("move = %b\n", moveOctet)
	fmt.Println("--------------------------------")

	return moveOctet
}

/*
Return: 0 -> NO WINNER, 1 -> X WON, 2 -> O WON
*/
func (solver *TicTacToeSolver) play(strategy1 []byte, strategy2 []byte) int {
	currentGrid := grid.InitialGrid()
	fmt.Printf("Initial board: %s\n", currentGrid.String())

	winner := grid.NO_WINNER
	for winner == grid.NO_WINNER && !currentGrid.IsFull() {
		stateNum := currentGrid.ToStateNumber()
		firstMove := GetMove(strategy1, stateNum)
		fmt.Printf("For state %d player 'x' makes a move at position %d. ", stateNum, firstMove)
		currentGrid.MakeMove(firstMove, 'x')

		fmt.Printf("New board: %s\n", currentGrid.String())

		winner = currentGrid.CheckWinner()
		if winner != grid.NO_WINNER {
			return winner
		} else if currentGrid.IsFull() {
			return grid.NO_WINNER
		}

		stateNum = currentGrid.ToStateNumber()
		secondMove := GetMove(strategy2, stateNum)
		fmt.Printf("For state %d player 'o' makes a move at position %d. ", stateNum, secondMove)
		currentGrid.MakeMove(secondMove, 'o')
		fmt.Printf("New board: %s\n", currentGrid.String())

		winner = currentGrid.CheckWinner()
	}

	return winner
}

func (solver *TicTacToeSolver) GenerateRandom() []byte {
	// Each byte can represent 2 states
	// Each state needs 4 bits to capture which cell to make the move.
	// Each octet can only have a value between 0000 and 1000
	const NumberStates = 19683

	//NOTE: Why the bytes.MinRead. Read http://openmymind.net/Go-Slices-And-The-Case-Of-The-Missing-Memory/ to understand
	genome := make([]byte, 0, NumberStates/2+bytes.MinRead)

	for i := 0; i < NumberStates/2; i++ {
		genome = append(genome, getRandomTwoMoves())
	}

	return genome
}

type Individual struct {
	genome  []byte
	fitness float64
}

type Population struct {
	individuals []Individual
	totalScore  float64
	size        int
}

type FitnessScoreFn func([]byte, []byte) float64

func computeNextGeneration(population Population, fitnessFn FitnessScoreFn, target []byte) []Individual {
	population.totalScore = 0.0
	for i, ind := range population.individuals {
		population.individuals[i].fitness = fitnessFn(ind.genome, target)
		population.totalScore += population.individuals[i].fitness
	}

	return population.individuals
}

func EvolutionComputation(evolver Evolver, populationSize int, maxGenerations int) {
	// Initialize the population
	population := Population{}
	for i := 0; i < populationSize; i++ {
		newIndividual := Individual{genome: evolver.GenerateRandom()}
		population.individuals = append(population.individuals, newIndividual)
	}

	bestSoFar := evolver.GenerateRandom()

	for i := 0; i < maxGenerations; i++ {
		population.individuals = computeNextGeneration(population, evolver.Fitness, bestSoFar)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	/*	state := [9]rune{' ', 'x', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
		thatGrid := grid.NewGrid(state)
		fmt.Printf("grid = %s, stateIndex = %d", thatGrid.String(), thatGrid.ToStateNumber())
	*/

	solver := TicTacToeSolver{}
	strategy1 := solver.GenerateRandom()
	strategy2 := solver.GenerateRandom()
	result := solver.play(strategy1, strategy2)
	fmt.Printf("Game result = %d", result)

	/*	fmt.Println("Board: " + ticGrid.String() + ", Winner: " + strconv.Itoa(ticGrid.CheckWinner()))
		ticGrid.MakeMove(0, 'o')
		ticGrid.MakeMove(3, 'o')
		ticGrid.MakeMove(6, 'o')
		fmt.Println("Board: " + ticGrid.String() + ", Winner: " + strconv.Itoa(ticGrid.CheckWinner()))*/
}
