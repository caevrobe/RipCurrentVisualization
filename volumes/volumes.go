package volumes

import (
	"PA1/framebuffer"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

const (
	SZ_UINT32 = 4
	SZ_UCHAR  = 1
)

var encoder = &png.Encoder{CompressionLevel: png.NoCompression}

type Pixel struct {
	R uint64
	G uint64
	B uint64
}

type Volume struct {
	filepath string
	*os.File
	width  int
	height int
	depth  int
}

func NewVolume(filepath string, width, height, depth int) *Volume {
	return &Volume{
		filepath,
		nil,
		width,
		height,
		depth,
	}
}

// TODO test if v.File methods are accessible from outside package
func (v *Volume) open() error {
	file, err := os.OpenFile(v.filepath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	v.File = file

	return nil
}

// creates image canvas with volume dimensions
func (v Volume) createImage() *image.RGBA {
	return image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{v.width, v.height},
	})
}

/* Scalar Volume ************************************************************/

type Scalar struct {
	*Volume
}

// TODO check for EOF or make sure numFrames doesn't exceed v.depth
func (v *Scalar) PullFrames(numFrames int) error {
	if err := v.open(); err != nil {
		return err
	}
	defer v.Close()

	frame := framebuffer.New(v.width, v.height, SZ_UCHAR, 3)

	for currFrame := 0; currFrame < numFrames; currFrame++ {
		frame.Reset()
		if _, err := v.Read(frame.Buffer); err != nil {
			return err
		}

		// TODO check read correct # of bytes

		img := v.createImage()
		color := color.RGBA{0, 0, 0, 0xFF}
		for y := 0; y < v.height; y++ {
			for x := 0; x < v.width; x++ {
				color.R = byte(frame.Next().(int))
				color.G = byte(frame.Next().(int))
				color.B = byte(frame.Next().(int))
				img.Set(x, y, color)
			}
		}

		out, _ := os.Create(fmt.Sprintf("temp/scalar%d.png", currFrame))
		encoder.Encode(out, img)
	}

	return nil
}

func (v *Scalar) HorizontalTimestack(numFrames int) error {
	if err := v.open(); err != nil {
		return err
	}
	defer v.Close()

	frame := framebuffer.New(v.width, v.height, SZ_UCHAR, 3)

	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{numFrames, v.height},
	})
	color := color.RGBA{0, 0, 0, 0xFF}
	cursor := 0

	for currFrame := 0; currFrame < numFrames; currFrame++ {
		if _, err := v.Read(frame.Buffer); err != nil {
			return err
		}
		xIndex := 900
		for y := 0; y < v.height; y++ {
			color.R = frame.Buffer[xIndex]
			color.G = frame.Buffer[xIndex+1]
			color.B = frame.Buffer[xIndex+2]
			img.Set(cursor, y, color)
			xIndex += v.width * 3
		}
		cursor++
	}

	out, _ := os.Create(fmt.Sprintf("temp/hstack%d.png", numFrames))
	encoder.Encode(out, img)

	return nil
}

func (v *Scalar) AverageFrames(numFrames int) error {
	if err := v.open(); err != nil {
		return err
	}
	defer v.Close()

	// slice containing RGB values of each pixel for one frame
	frame := make([]byte, v.height*v.width*3)

	avg := make([]Pixel, v.width*v.height)

	img := v.createImage()

	for currFrame := 0; currFrame < numFrames; currFrame++ {
		if _, err := v.Read(frame); err != nil {
			return err
		}

		// TODO check correct # of bytes read

		for idx := 0; idx < v.height*v.width*3; idx += 3 {
			currentPixel := &avg[idx/3]
			currentPixel.R += uint64(frame[idx])
			currentPixel.G += uint64(frame[idx+1])
			currentPixel.B += uint64(frame[idx+2])
		}
	}

	color := color.RGBA{0, 0, 0, 0xFF}
	cursor := 0
	numFrames_64 := uint64(numFrames)
	for y := 0; y < v.height; y++ {
		for x := 0; x < v.width; x++ {
			color.R = uint8(avg[cursor].R / numFrames_64)
			color.G = uint8(avg[cursor].G / numFrames_64)
			color.B = uint8(avg[cursor].B / numFrames_64)
			img.Set(x, y, color)

			cursor++
		}
	}

	out, _ := os.Create(fmt.Sprintf("temp/scalar_avg%d.png", numFrames))
	encoder.Encode(out, img)

	return nil
}

/* Vector Volume ************************************************************/

type Vector struct {
	*Volume
}

func (v *Vector) PullFrames(numFrames int) error {
	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{v.width, v.height},
	})
	color := color.RGBA{0, 255, 0, 0xFF}
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color)
	draw2dkit.Circle(gc, 345, 546, 200)
	gc.Fill()

	out, _ := os.Create(fmt.Sprintf("temp/scalar%d.png", 2))
	png.Encode(out, img)

	return nil
}
