package volumes

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// TODO add scalar timestack function

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

// TODO test if file closes when v.file does
// TODO test if v.File methods are accessible from outside package
func (v *Volume) Open() error {
	file, err := os.OpenFile(v.filepath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	v.File = file

	return nil
}

/* Scalar Volume ************************************************************/

type Scalar struct {
	*Volume
}

// TODO check for EOF or make sure numFrames doesn't exceed v.depth
func (v *Scalar) PullFrames(numFrames int) error {
	if err := v.Open(); err != nil {
		return err
	}
	defer v.Close()

	frame := make([]byte, v.height*v.width*3)

	for currFrame := 0; currFrame < numFrames; currFrame++ {
		cursor := 0
		if _, err := v.Read(frame); err != nil {
			return err
		}

		// TODO check read correct # of bytes

		img := image.NewRGBA(image.Rectangle{
			image.Point{0, 0},
			image.Point{v.width, v.height},
		})
		color := color.RGBA{0, 0, 0, 0xFF}
		for y := 0; y < v.height; y++ {
			for x := 0; x < v.width; x++ {
				color.R = frame[cursor]
				color.G = frame[cursor+1]
				color.B = frame[cursor+2]
				cursor += 3
				img.Set(x, y, color)
			}
		}

		out, _ := os.Create(fmt.Sprintf("temp/scalar%d.png", currFrame))
		png.Encode(out, img)
	}

	return nil
}

func (v *Scalar) AverageFrames(numFrames int) error {
	if err := v.Open(); err != nil {
		return err
	}
	defer v.Close()

	// slice containing RGB values of each pixel for one frame
	frame := make([]byte, v.height*v.width*3)

	avg := make([]Pixel, v.width*v.height)
	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{v.width, v.height},
	})

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
	png.Encode(out, img) // TODO check error

	return nil
}

/* Vector Volume ************************************************************/

type Vector struct {
	*Volume
}
