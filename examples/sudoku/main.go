package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Kappeh/dlx"
)

type Sudoku struct {
	Unsolved      string
	Solved        string
	FoundSolution bool
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, string(body))
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	sudokus, ok := r.URL.Query()["sudoku"]
	if !ok || len(sudokus[0]) < 81 {
		return
	}
	sudoku := sudokus[0]
	mat, err := sudokuDLX(sudoku)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := Sudoku{
		Unsolved:      sudoku,
		Solved:        "",
		FoundSolution: false,
	}
	solution := dlx.FirstSolution(mat)
	if solution != nil {
		result.Solved = intSliceToString(solution)
		result.FoundSolution = true
	}
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func baseSudokuDLX() (*dlx.Matrix, error) {
	result, err := dlx.New(324, 0)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 729; i++ {
		cellIndex := i / 9
		cellValue := i % 9
		rowIndex := cellIndex / 9
		colIndex := cellIndex % 9
		boxRowIndex := rowIndex / 3
		boxColIndex := colIndex / 3
		boxIndex := 3*boxRowIndex + boxColIndex
		col1 := cellIndex
		col2 := 81 + 9*rowIndex + cellValue
		col3 := 162 + 9*colIndex + cellValue
		col4 := 243 + 9*boxIndex + cellValue
		dlx.AddRow(result, col1, col2, col3, col4)
	}
	return result, nil
}

func sudokuDLX(s string) (*dlx.Matrix, error) {
	mat, err := baseSudokuDLX()
	if err != nil {
		return nil, err
	}
	for i, c := range s {
		if c == '0' {
			continue
		}
		number, err := strconv.Atoi(string(c))
		if err != nil {
			return nil, err
		}
		dlx.AddToSolution(mat, i*9+number-1)
	}
	return mat, nil
}

func intSliceToString(s []int) string {
	result := make([]rune, 81)
	for _, v := range s {
		cellIndex := v / 9
		cellValue := v % 9
		result[cellIndex] = rune("123456789"[cellValue])
	}
	return string(result)
}

func main() {
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/solve", solveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
