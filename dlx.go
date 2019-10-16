package dlx

import (
	"errors"
)

/*
New constructs a new matrix.
primary 	- The amount of primary columns to be in the matrix.
optional 	- The amount of optional columns to be in the matrix.
*/
func New(primary, optional int) (*Matrix, error) {
	if primary <= 0 {
		return nil, errors.New("primary must be positive")
	}
	if optional < 0 {
		return nil, errors.New("optional must be non-negative")
	}
	cols := primary + optional
	result := Matrix{
		root:     element{},
		headers:  make([]element, cols),
		details:  make([]element, 0),
		rows:     make([]row, 0),
		rowCount: 0,
		solution: newSolution(),
	}
	result.root.up, result.root.down = &result.root, &result.root
	for i := 0; i < primary; i++ {
		result.headers[i] = element{
			size:   0,
			row:    nil,
			header: &result.headers[i],
			up:     &result.headers[i],
			down:   &result.headers[i],
			left:   &result.headers[negativeMod(i-1, primary)],
			right:  &result.headers[(i+1)%primary],
		}
	}
	for i := primary; i < cols; i++ {
		result.headers[i] = element{
			size:   0,
			row:    nil,
			header: &result.headers[i],
			up:     &result.headers[i],
			down:   &result.headers[i],
			left:   &result.headers[i],
			right:  &result.headers[i],
		}
	}
	result.root.left = result.headers[0].left
	result.root.right = &result.headers[0]
	result.root.left.right = &result.root
	result.root.right.left = &result.root
	return &result, nil
}

/*
AddRow adds a new row of elements to a matrix.
dlx 	- The matrix to add the row to.
indexes - The indices of the columns containing 1s.
*/
func AddRow(dlx *Matrix, indexes ...int) error {
	count := len(indexes)
	if count == 0 {
		return nil
	}
	last := -1
	for _, index := range indexes {
		if index < 0 || index >= len(dlx.headers) {
			return errors.New("index out of range")
		}
		if last != -1 && index <= last {
			return errors.New("indexes not in ascending order")
		}
		last = index
	}
	length := len(dlx.details)
	dlx.details = append(dlx.details, make([]element, count)...)
	newRow := dlx.details[length:]
	dlx.rows = append(dlx.rows, row{
		index:        dlx.rowCount,
		firstElement: &newRow[0],
		covered:      false,
	})
	for i, col := range indexes {
		newRow[i] = element{
			size:   0,
			header: &dlx.headers[col],
			up:     dlx.headers[col].up,
			down:   &dlx.headers[col],
			left:   &newRow[negativeMod(i-1, count)],
			right:  &newRow[(i+1)%count],
			row:    &dlx.rows[dlx.rowCount],
		}
		dlx.headers[col].size++
		newRow[i].up.down, newRow[i].down.up = &newRow[i], &newRow[i]
	}
	dlx.rowCount++
	return nil
}

/*
AddToSolution adds a row explicitly to the solution for a matrix.
dlx		- The matrix which contains the row.
index	- The index of the row.
*/
func AddToSolution(dlx *Matrix, index int) error {
	if index < 0 || index >= dlx.rowCount {
		return errors.New("index out of range")
	}
	if dlx.rows[index].covered {
		return errors.New("row is covered, cannot be included in solution")
	}
	dlx.solution.push(index)
	firstElement := dlx.rows[index].firstElement
	coverRow(firstElement)
	coverColumn(firstElement.header)
	for e := firstElement.right; e != firstElement; e = e.right {
		coverColumn(e.header)
	}
	return nil
}

/*
ClearSolution removes all rows from the current solution for a matrix.
This function undoes any calls to AddToSolution.
dlx - The matrix to clear the solution for.
*/
func ClearSolution(dlx *Matrix) {
	for dlx.solution.size() > 0 {
		index, _ := dlx.solution.pop()
		firstElement := dlx.rows[index].firstElement
		for e := firstElement.left; e != firstElement; e = e.left {
			uncoverColumn(e.header)
		}
		uncoverColumn(firstElement.header)
		uncoverRow(firstElement)
	}
}

/*
ForEachSolution calls f with a slice of all row indexes which correspond
to a solution for a matrix.
dlx	- The matrix to find solutions for.
f	- The function to be called when a solution is found.
*/
func ForEachSolution(dlx *Matrix, f func([]int)) {
	if dlx.root.left == &dlx.root {
		f(dlx.solution.values[:dlx.solution.stackptr])
		return
	}
	header, emptyColumn := colPolicy(dlx)
	if emptyColumn {
		return
	}
	for r := header.down; r != header; r = r.down {
		dlx.solution.push(r.row.index)
		coverRow(r)
		coverColumn(r.header)
		for j := r.right; j != r; j = j.right {
			coverColumn(j.header)
		}
		ForEachSolution(dlx, f)
		for j := r.left; j != r; j = j.left {
			uncoverColumn(j.header)
		}
		uncoverColumn(r.header)
		uncoverRow(r)
		dlx.solution.pop()
	}
}

/*
FirstSolution finds a solution for a matrix and returns the row indexes.
dlx	- The matrix to find a solution for.
*/
func FirstSolution(dlx *Matrix) []int {
	if dlx.root.left == &dlx.root {
		return dlx.solution.values[:dlx.solution.stackptr]
	}
	header, emptyColumn := colPolicy(dlx)
	if emptyColumn {
		return nil
	}
	for r := header.down; r != header; r = r.down {
		dlx.solution.push(r.row.index)
		coverRow(r)
		coverColumn(r.header)
		for j := r.right; j != r; j = j.right {
			coverColumn(j.header)
		}
		result := FirstSolution(dlx)
		for j := r.left; j != r; j = j.left {
			uncoverColumn(j.header)
		}
		uncoverColumn(r.header)
		uncoverRow(r)
		dlx.solution.pop()
		if result != nil {
			return result
		}
	}
	return nil
}

func newSolution() solution {
	return solution{
		values:   make([]int, 0),
		stackptr: 0,
	}
}

func colPolicy(dlx *Matrix) (*element, bool) {
	// TODO: Make more efficient column policy
	// algorithm than a linear search
	var best *element
	for h := dlx.root.right; h != &dlx.root; h = h.right {
		if best == nil || h.size < best.size {
			best = h
		}
		if best.size == 0 {
			return nil, true
		}
	}
	return best, false
}

func coverRow(e *element) {
	e.row.covered = true
	e = e.right
	e.up.down, e.down.up = e.down, e.up
	e.header.size--
	for r := e.right; r != e; r = r.right {
		r.up.down, r.down.up = r.down, r.up
		r.header.size--
	}
}

func uncoverRow(e *element) {
	e.row.covered = false
	e = e.left
	e.up.down, e.down.up = e, e
	e.header.size++
	for r := e.left; r != e; r = r.left {
		r.up.down, r.down.up = r, r
		r.header.size++
	}
}

func coverColumn(h *element) {
	h.left.right, h.right.left = h.right, h.left
	for j := h.down; j != h; j = j.down {
		j.left.right, j.right.left = j.right, j.left
		coverRow(j)
	}
}

func uncoverColumn(h *element) {
	for j := h.up; j != h; j = j.up {
		uncoverRow(j)
		j.left.right, j.right.left = j, j
	}
	h.left.right, h.right.left = h, h
}

func negativeMod(a, b int) int {
	result := a % b
	if (result < 0 && b > 0) || (result > 0 && b < 0) {
		return result + b
	}
	return result
}
