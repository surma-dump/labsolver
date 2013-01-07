package main

import (
	"image"
	"log"
)

type Vector2 struct {
	image.Point
}

func NewVector2(x, y int) *Vector2 {
	return &Vector2{
		image.Point{
			X: x,
			Y: y,
		},
	}
}

// Usually, a rotation of 90 degrees counter-clock wise would be
// ⎛ 0 -1 ⎞ · ⎛ v.X ⎞
// ⎝ 1  0 ⎠   ⎝ v.Y ⎠
// But our y axis is inverted, so we need an transform our vector
// to the normal coordinate system, do the rotation, and reverse
// the transformation. I.e:
// ⎛ 1  0 ⎞ · ⎛ 0 -1 ⎞ · ⎛ 1  0 ⎞ · ⎛ v.X ⎞
// ⎝ 0 -1 ⎠   ⎝ 1  0 ⎠   ⎝ 0 -1 ⎠   ⎝ v.Y ⎠
func (v *Vector2) RotateLeft() {
	v.X, v.Y = v.Y, -v.X
}

type LabyrinthSolver struct {
	Walker LabyrinthWalker
}

func (s *LabyrinthSolver) Solve() {
	for !s.Walker.Done() {
		for s.Walker.WallAhead() {
			s.Walker.TurnLeft()
		}
		s.Walker.Walk()
	}
}

type LabyrinthWalker interface {
	Walk()
	TurnLeft()
	WallAhead() bool
	Done() bool
}

type DumpWalker struct {
	LabyrinthWalker
	StepCount int
}

func (dw *DumpWalker) Walk() {
	dw.StepCount++
	dw.LabyrinthWalker.Walk()
}

func (dw *DumpWalker) TurnLeft() {
	if dw.StepCount > 0 {
		log.Printf("Walk %d step(s)", dw.StepCount)
	}
	log.Printf("turn left")
	dw.StepCount = 0
	dw.LabyrinthWalker.TurnLeft()
}

func (dw *DumpWalker) Done() bool {
	done := dw.LabyrinthWalker.Done()
	if done {
		log.Printf("Walk %d step(s), done", dw.StepCount)
	}
	return done
}

type ImageWalker struct {
	image    image.Image
	pos, end *Vector2
	dir      *Vector2
}

func NewImageWalker(img image.Image, start, end *Vector2) *ImageWalker {
	return &ImageWalker{
		image: img,
		end:   end,
		pos:   start,
		dir:   NewVector2(1, 0),
	}
}

func (iw *ImageWalker) Walk() {
	if !iw.WallAhead() {
		iw.pos.Point = iw.pos.Add(iw.dir.Point)
	}
}

func (iw *ImageWalker) TurnLeft() {
	iw.dir.RotateLeft()
}

func (iw *ImageWalker) WallAhead() bool {
	// Would the next step move us out of image bounds?
	if !iw.pos.Add(iw.dir.Point).In(iw.image.Bounds()) {
		return true
	}
	// Fix me: Wall detection
	return false
}

func (iw *ImageWalker) Done() bool {
	return iw.pos.Eq(iw.end.Point)
}
