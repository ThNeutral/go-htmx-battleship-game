package main

import "github.com/google/uuid"

type TFields struct {
	First  [][]string
	Second [][]string
}

type TMap struct {
	TFields
	FirstKey  string
	SecondKey string
}

const (
	EMPTY_CELL     = " "
	FILLED_CELL    = "/"
	DESTROYED_CELL = "x"
	MISSED_CELL    = "o"
)

func getNewMap(rows, cols int) *TMap {
	var m TMap
	m.First = indexTo2DArray(rows, cols, EMPTY_CELL)
	m.FirstKey = uuid.NewString()
	m.Second = indexTo2DArray(rows, cols, EMPTY_CELL)
	m.SecondKey = uuid.NewString()
	return &m
}

func (m *TMap) GetFields() *TFields {
	return &TFields{
		First:  m.First,
		Second: m.Second,
	}
}

func (m *TMap) SetFirst(row, col int, val string) {
	m.First[row][col] = val
}

func (m *TMap) SetSecond(row, col int, val string) {
	m.Second[row][col] = val
}
