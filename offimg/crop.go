package offimg

import (
	"fmt"
	"image"
	"reflect"

	"github.com/nfnt/resize"
)

type CropSpec struct {
	DefinedPtOne image.Point
	DefinedPtTwo image.Point
}

func (c *CropSpec) CropResize(img image.Image, width uint, height uint) (image.Image, error) {
	asp := float64(width) / float64(height)
	bndBox, err := c.resolve(img, asp)
	if err != nil {
		return nil, err
	}
	subImage := reflect.ValueOf(img).MethodByName("SubImage").Call([]reflect.Value{reflect.ValueOf(*bndBox)})[0].Interface().(image.Image)
	return resize.Resize(width, height, subImage, resize.Lanczos3), nil
}

func (c *CropSpec) Crop(img image.Image, asp float64) (image.Image, error) {
	bndBox, err := c.resolve(img, asp)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(img).MethodByName("SubImage").Call([]reflect.Value{reflect.ValueOf(*bndBox)})[0].Interface().(image.Image), nil
}

func (c *CropSpec) resolve(img image.Image, asp float64) (*image.Rectangle, error) {
	if c.DefinedPtOne.X == c.DefinedPtTwo.X {
		ySpan := c.DefinedPtTwo.Y - c.DefinedPtOne.Y
		if ySpan < 0 {
			ySpan = -ySpan
			c.DefinedPtOne.Y, c.DefinedPtTwo.Y = c.DefinedPtTwo.Y, c.DefinedPtOne.Y
		} else if ySpan == 0 {
			return nil, fmt.Errorf("crop: both points are the same")
		}
		xSpan := int(float64(ySpan) * asp)
		c.DefinedPtTwo.X += xSpan
	} else if c.DefinedPtOne.Y == c.DefinedPtTwo.Y {
		xSpan := c.DefinedPtTwo.X - c.DefinedPtOne.X
		if xSpan < 0 {
			xSpan = -xSpan
			c.DefinedPtOne.X, c.DefinedPtTwo.X = c.DefinedPtTwo.X, c.DefinedPtOne.X
		}
		ySpan := int(float64(xSpan) / asp)
		c.DefinedPtTwo.Y += ySpan
	}
	return &image.Rectangle{
		Min: c.DefinedPtOne,
		Max: c.DefinedPtTwo,
	}, nil
}
