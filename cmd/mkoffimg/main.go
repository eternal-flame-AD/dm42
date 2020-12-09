package main

import (
	"flag"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	dither "github.com/esimov/dithergo"
	"github.com/eternal-flame-AD/dm42/offimg"
)

var (
	flagInputFile   = flag.String("f", "", "input image file")
	flagOutputFile  = flag.String("o", "output.bmp", "output image file")
	flagDoDithering = flag.Bool("dither", true, "do dithering")
)

var ditherer = dither.Dither{
	"FloydSteinberg",
	dither.Settings{
		[][]float32{
			[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
			[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
			[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
		},
	},
}

func main() {
	flag.Parse()

	f, err := os.Open(*flagInputFile)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	if *flagDoDithering {
		img = ditherer.Monochrome(img, 1.18)
	}

	f, err = os.Create(*flagOutputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var encodeErr error
	switch strings.ToLower(filepath.Ext(*flagOutputFile)) {
	case ".bmp":
		encodeErr = offimg.WriteImage(img, f)
	case ".png":
		encodeErr = png.Encode(f, img)
	case ".jpg", ".jpeg":
		encodeErr = jpeg.Encode(f, img, nil)
	}
	if encodeErr != nil {
		panic(encodeErr)
	}
}
