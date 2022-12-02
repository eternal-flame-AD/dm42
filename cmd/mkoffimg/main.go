package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"

	dither "github.com/esimov/dithergo"
	"github.com/eternal-flame-AD/dm42/offimg"
	"github.com/nfnt/resize"
)

var (
	flagInputFile   = flag.String("f", "", "input image file")
	flagOutputFile  = flag.String("o", "output.bmp", "output image file")
	flagDoDithering = flag.Bool("dither", true, "do dithering")

	flagCropPoints = flag.String("crop", "", "crop point x1:y1,x2:y2")
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

	if *flagCropPoints != "" {
		cropSpec := &offimg.CropSpec{}
		_, err := fmt.Sscanf(*flagCropPoints, "%d:%d,%d:%d", &cropSpec.DefinedPtOne.X, &cropSpec.DefinedPtOne.Y, &cropSpec.DefinedPtTwo.X, &cropSpec.DefinedPtTwo.Y)
		if err != nil {
			log.Panicf("invalid crop points: %s", *flagCropPoints)
		}
		img, err = cropSpec.CropResize(img, offimg.Width, offimg.Height)
		if err != nil {
			log.Panicf("crop/resize error: %v", err)
		}
	} else if img.Bounds().Dx() != offimg.Width || img.Bounds().Dy() != offimg.Height {
		img = resize.Resize(offimg.Width, offimg.Height, img, resize.Lanczos3)
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
