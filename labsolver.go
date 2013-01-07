package main

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/voxelbrain/goptions"
)

type Crop [4]int

const (
	LEFT = iota
	TOP
	RIGHT
	BOTTOM
)

var (
	options = struct {
		Image               *os.File      `goptions:"-f, --file, obligatory, rdonly, description='Image file to read'"`
		Output              *os.File      `goptions:"-o, --output, obligatory, create, wronly, description='File to write to'"`
		Crop                *Crop         `goptions:"-c, --crop, description='Crop [left, top, right, bottom]'"`
		StartPosition       *Vector2      `goptions:"--start, obligatory, description='Start coordinates in pixels (pre-crop)'"`
		EndPosition         *Vector2      `goptions:"--end, obligatory, description='End coordinates in pixels (pre-crop)'"`
		BrightnessThreshold float64       `goptions:"-b, --brightness-threshold, description='Values above brightness threshold are walls'"`
		Help                goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		BrightnessThreshold: 0.5,
		Crop:                &Crop{0, 0, 0, 0},
		StartPosition:       NewVector2(0, 0),
	}
)

func init() {
	goptions.ParseAndFail(&options)
}

func main() {
	defer options.Image.Close()
	defer options.Output.Close()
	img, _, err := image.Decode(options.Image)
	if err != nil {
		log.Fatalf("Could not decode image: %s", err)
	}
	if ci, ok := img.(CropImage); ok {
		cpy := copyImage(img)
		img = ci.SubImage(image.Rect(options.Crop[LEFT],
			options.Crop[RIGHT],
			img.Bounds().Max.X-options.Crop[RIGHT],
			img.Bounds().Max.Y-options.Crop[BOTTOM]).Canon())
		iw := NewImageWalker(img, NewBrightnessWallDetector(options.BrightnessThreshold, img), options.StartPosition, options.EndPosition)

		ls := &LabyrinthSolver{NewDrawWalker(cpy, iw)}
		ls.Solve()
		err := png.Encode(options.Output, cpy)
		if err != nil {
			log.Fatalf("Could not save image: %s", err)
		}
		return
	}
	log.Fatalf("Could not crop image")
}

type CropImage interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}
