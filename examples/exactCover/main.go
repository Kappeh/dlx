/*
Information of the exact cover problem can be found at:
https://en.wikipedia.org/wiki/Exact_cover

The following example will show how to use the dlx package
to solve the exact cover problem for the following:

Let S = {A, B, C, D, E, F} be a collection of subsets of a
set X = {1, 2, 3, 4, 5, 6, 7} such that:
A = {1, 4, 7}
B = {1, 4}
C = {4, 5, 7}
D = {3, 5, 6}
E = {2, 3, 6, 7}
F = {2, 7}

This can also be represented as the following matrix:

	    1 2 3 4 5 6 7

	A   1 0 0 1 0 0 1
	B   1 0 0 1 0 0 0
	C   0 0 0 1 1 0 1
	D   0 0 1 0 1 1 0
	E   0 1 1 0 0 1 1
	F   0 1 0 0 0 0 1

The subcollection S* = {B, D, F} is the only exact cover
in this case. Our aim is to use the dlx package to
arrive at this conclusion programatically.
*/

package main

import (
	"log"
	"fmt"
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

	// Row0 <- A = {1, 4, 7}
	handleErr(dlx.AddRow(s, 0, 3, 6))
	// Row1 <- B = {1, 4}
	handleErr(dlx.AddRow(s, 0, 3))
	// Row2 <- C = {4, 5, 7}
	handleErr(dlx.AddRow(s, 3, 4, 6))
	// Row3 <- D = {3, 5, 7}
	handleErr(dlx.AddRow(s, 2, 4, 5))
	// Row4 <- E = {2, 3, 6, 7}
	handleErr(dlx.AddRow(s, 1, 2, 5, 6))
	// Row5 <- F = {2, 7}
	handleErr(dlx.AddRow(s, 1, 6))

	count := 0

	// For each solution
	dlx.ForEachSolution(s, func(s []int) {
		// Covert the row indexes to letters
		letters := make([]string, len(s))
		for i, v := range s {
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
