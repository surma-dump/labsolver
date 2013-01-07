package main

import ()

type Vector2 struct {
	X, Y int
}

// Rotate the vector 90 degrees counter-clock wise
// ⎛ 0 -1 ⎞ · ⎛ v.X ⎞
// ⎝ 1  0 ⎠   ⎝ v.Y ⎠
func (v *Vector2) RotateLeft() {
	v.X, v.Y = -v.Y, v.X
}

type LabyrinthSolver struct {
	Labyrinth  Labyrinth
	Start, End Vector2
}

type Type int

const (
	NONE Type = iota
	WALL
	VISITED
)

type Labyrinth interface {
	Walk()
	TurnLeft()
	Look() (left, front, right, bottom Type)
}
