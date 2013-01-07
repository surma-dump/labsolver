package main

import (
	"image"
	"image/draw"
	"log"
	"math"
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

func (v *Vector2) RotateLeft() *Vector2 {
	// Usually, a rotation of 90 degrees counter-clock wise would be
	// ⎛ 0 -1 ⎞ · ⎛ v.X ⎞
	// ⎝ 1  0 ⎠   ⎝ v.Y ⎠
	// But our y axis is inverted, so we need an transform our vector
	// to the normal coordinate system, do the rotation, and reverse
	// the transformation. I.e:
	// ⎛ 1  0 ⎞ · ⎛ 0 -1 ⎞ · ⎛ 1  0 ⎞ · ⎛ v.X ⎞
	// ⎝ 0 -1 ⎠   ⎝ 1  0 ⎠   ⎝ 0 -1 ⎠   ⎝ v.Y ⎠
	return NewVector2(v.Y, -v.X)
}

type LabyrinthSolver struct {
	Walker LabyrinthWalker
}

func (s *LabyrinthSolver) Solve() {
	for !s.Walker.Done() {
		left, front, right := s.Walker.Look()
		if !left {
			s.Walker.TurnLeft()
		} else if front && !right {
			s.Walker.TurnLeft()
			s.Walker.TurnLeft()
			s.Walker.TurnLeft()
		} else if left && front && right {
			s.Walker.TurnLeft()
			s.Walker.TurnLeft()
		}
		s.Walker.Walk()
	}
}

type LabyrinthWalker interface {
	Walk()
	TurnLeft()
	Look() (left, front, right bool)
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

type WallDetector func(x, y int) bool

func NewBrightnessWallDetector(brightnessThreshold float64, img image.Image) WallDetector {
	return func(x, y int) bool {
		if !(image.Point{X: x, Y: y}).In(img.Bounds()) {
			return true
		}
		r, g, b, _ := img.At(x, y).RGBA()
		brightness := math.Sqrt(float64(r*r)+float64(g*g)+float64(b*b)) / math.Sqrt(3*float64(0xFFFF*0xFFFF))
		return brightness > brightnessThreshold
	}
}

type ImageWalker struct {
	image    image.Image
	wd       WallDetector
	pos, end *Vector2
	dir      *Vector2
}

func NewImageWalker(img image.Image, wd WallDetector, start, end *Vector2) *ImageWalker {
	return &ImageWalker{
		image: img,
		wd:    wd,
		end:   end,
		pos:   start,
		dir:   NewVector2(1, 0),
	}
}

func (iw *ImageWalker) Walk() {
	if _, front, _ := iw.Look(); !front {
		iw.pos.Point = iw.pos.Add(iw.dir.Point)
	}
}

func (iw *ImageWalker) TurnLeft() {
	iw.dir = iw.dir.RotateLeft()
}

func (iw *ImageWalker) Look() (left, front, right bool) {
	frontPos := iw.pos.Add(iw.dir.Point)
	leftPos := iw.pos.Add(iw.dir.RotateLeft().Point)
	rightPos := iw.pos.Add(iw.dir.RotateLeft().RotateLeft().RotateLeft().Point)
	return iw.wd(leftPos.X, leftPos.Y), iw.wd(frontPos.X, frontPos.Y), iw.wd(rightPos.X, rightPos.Y)
}

func (iw *ImageWalker) Done() bool {
	return iw.pos.Eq(iw.end.Point)
}

func copyImage(img image.Image) draw.Image {
	newimg := image.NewRGBA(img.Bounds())
	draw.Draw(newimg, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	return newimg
}
