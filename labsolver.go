package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/voxelbrain/goptions"
)

type Crop [4]int

var (
	options = struct {
		Image              *os.File      `goptions:"-f, --file, obligatory, rdonly, description='Image file to read'"`
		Crop               *Crop         `goptions:"-c, --crop, description='Crop [left, top, right, bottom]'"`
		StartPosition      *Vector2      `goptions:"--start, obligatory, description='Start coordinates in pixels (post-crop)'"`
		EndPosition        *Vector2      `goptions:"--end, obligatory, description='End coordinates in pixels (post-crop)'"`
		Invert             bool          `goptions:"-i, --invert, description='Invert walls and paths'"`
		LightnessThreshold float32       `goptions:"-l, --lightness-threshold, description='HSL-Values above lightness threshold are walls'"`
		Help               goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		LightnessThreshold: 0.5,
		Crop:               &Crop{0, 0, 0, 0},
		StartPosition:      &Vector2{0, 0},
	}
)

func init() {
	goptions.ParseAndFail(&options)
}

func main() {
	img, _, err := image.Decode(options.Image)
	if err != nil {
		log.Fatalf("Could not decode image: %s", err)
	}
	_ = img
}
