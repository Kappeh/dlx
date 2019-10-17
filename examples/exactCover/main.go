package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Kappeh/dlx"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Making new matrix S
	s, err := dlx.New(7, 0)
	handleErr(err)

	handleErr(dlx.AddRow(s, 0, 3, 6))    // Row0 <- A = {1, 4, 7}
	handleErr(dlx.AddRow(s, 0, 3))       // Row1 <- B = {1, 4}
	handleErr(dlx.AddRow(s, 3, 4, 6))    // Row2 <- C = {4, 5, 7}
	handleErr(dlx.AddRow(s, 2, 4, 5))    // Row3 <- D = {3, 5, 7}
	handleErr(dlx.AddRow(s, 1, 2, 5, 6)) // Row4 <- E = {2, 3, 6, 7}
	handleErr(dlx.AddRow(s, 1, 6))       // Row5 <- F = {2, 7}

	count := 0

	// For each solution
	dlx.ForEachSolution(s, func(rows []int) {
		// Covert the row indexes to letters
		letters := make([]string, len(rows))
		for i, v := range rows {
			letters[i] = string("ABCDEF"[v])
		}
		// Combine them in set notation for printing
		setString := "{" + strings.Join(letters, ", ") + "}"
		fmt.Printf("Solution found: %s.\n", setString)
		// Update the solution counter
		count++
	})

	fmt.Printf("Finished. Found %d solution(s).\n", count)
}
