package offimg

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

const Width = 400
const Height = 240

type BitmapFileHeader struct {
	BfType      [2]byte
	BfSize      uint32
	BfReserved1 uint16
	BfReserved2 uint16
	BfOffBits   uint32
}

type BitmapInfoHeader struct {
	BfSize           uint32
	BitWidth         uint32
	BitHeight        uint32
	BiPlanes         uint16
	BiBitCount       uint16
	BiCompressing    uint32
	BiSizeImage      uint32
	BiXPixelPerMeter uint32
	BiYPixelPerMeter uint32
	BiClrUsed        uint32
	BiClrImportant   uint32
}

type BitmapColors []RGBQuad

type RGBQuad struct {
	Blue, Green, Red, Reserved uint8
}

func White() RGBQuad {
	return RGBQuad{255, 255, 255, 0}
}

func Black() RGBQuad {
	return RGBQuad{0, 0, 0, 0}
}

func WriteImage(img image.Image, w io.Writer) error {
	if img.Bounds().Dx() != Width {
		return fmt.Errorf("width is %d, expect %d px", img.Bounds().Dx(), Width)
	}
	if img.Bounds().Dy() != Height {
		return fmt.Errorf("height is %d, expect %d px", img.Bounds().Dy(), Height)
	}

	fileHeader := BitmapFileHeader{
		BfType:    [2]byte{'B', 'M'},
		BfSize:    0, //TO BE FILLED sizeof(file)
		BfOffBits: 0, //TO BE FILLED sizeof(fileheader+infoHeader+palette)
	}
	fileInfoHeader := BitmapInfoHeader{
		BfSize:           0, //TO BE FILLED sizeof(infoheader)
		BitWidth:         Width,
		BitHeight:        Height,
		BiPlanes:         1,
		BiBitCount:       1,
		BiCompressing:    0,
		BiSizeImage:      0,
		BiXPixelPerMeter: 2834, // taken from smpl image
		BiYPixelPerMeter: 2834,
		BiClrUsed:        0,
		BiClrImportant:   0,
	}
	fileInfoHeader.BfSize = uint32(binary.Size(fileInfoHeader))
	colors := [2]RGBQuad{White(), Black()}
	fileHeader.BfOffBits = uint32(binary.Size(fileHeader)) + fileInfoHeader.BfSize + uint32(binary.Size(colors))
	bytesPerLine := (Width + 31) / 32 * 32 / 8
	imageSize := bytesPerLine * Height
	fileInfoHeader.BiSizeImage = uint32(imageSize)
	fileHeader.BfSize = fileHeader.BfOffBits + uint32(imageSize)
	paddingBytes := (fileHeader.BfSize+7)/8*8 - fileHeader.BfSize
	fileHeader.BfSize += paddingBytes
	fileInfoHeader.BiSizeImage += paddingBytes

	if err := binary.Write(w, binary.LittleEndian, fileHeader); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, fileInfoHeader); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, colors); err != nil {
		return err
	}

	for y := img.Bounds().Max.Y - 1; y >= img.Bounds().Min.Y; y-- {
		lineBuffer := make([]byte, bytesPerLine)
		bitPlace := 0
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if (r+g+b)/3 < (a)/2 /* black -> 1 */ {
				lineBuffer[bitPlace/8] += 1 << (7 - bitPlace%8)
			}
			bitPlace++
		}
		if _, err := w.Write(lineBuffer); err != nil {
			return err
		}
	}
	pad := make([]byte, paddingBytes)
	if _, err := w.Write(pad); err != nil {
		return err
	}
	return nil
}
