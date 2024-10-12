package main

func indexTo2DArray(rows, cols int, val string) [][]string {
	matrix := make([][]string, rows)
	for row, _ := range matrix {
		matrix[row] = make([]string, cols)
		for col, _ := range matrix[row] {
			matrix[row][col] = val
		}
	}
	return matrix
}
