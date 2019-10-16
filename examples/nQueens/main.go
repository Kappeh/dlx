package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Kappeh/dlx"
)

var (
	n int
)

func init() {
	flag.IntVar(&n, "n", 8, "The number of queens. Must be greater than 0.")
	flag.Parse()

	if n <= 0 {
		log.Fatal("flag n expected positive number")
	}
}

func nQueenDLX(queens int) (*dlx.Matrix, error) {
	primary := 2 * queens
	diagonals := queens*2 - 3
	optional := 2 * diagonals
	size := queens * queens
	matrix, err := dlx.New(primary, optional)
	if err != nil {
		return nil, err
	}

	for i := 0; i < size; i++ {
		row := i / queens
		col := i % queens
		posDiag := row + col - 1
		negDiag := row - col + queens - 2
		indices := make([]int, 2)
		indices[0] = row
		indices[1] = col + queens
		if posDiag >= 0 && posDiag < diagonals {
			indices = append(indices, posDiag+2*queens)
		}
		if negDiag >= 0 && negDiag < diagonals {
			indices = append(indices, negDiag+2*queens+diagonals)
		}
		err := dlx.AddRow(matrix, indices...)
		if err != nil {
			return nil, err
		}
	}

	return matrix, nil
}

func solutionString(queens int, rows []int) string {
	size := queens * queens
	tiles := make([]bool, size)
	for _, v := range rows {
		tiles[v] = true
	}
	result := ""
	for i := 0; i < size; i++ {
		if i > 0 && i%queens == 0 {
			result += "\n"
		}
		if tiles[i] {
			result += "Q "
		} else {
			result += "- "
		}
	}
	return result
}

func main() {
	matrix, err := nQueenDLX(n)
	if err != nil {
		log.Fatal(err)
	}
	solution := dlx.FirstSolution(matrix)
	fmt.Println(solutionString(n, solution))
}
