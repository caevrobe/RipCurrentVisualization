package volumes

import (
	"PA1/framebuffer"
	"bufio"
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"

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
				color.R = byte(frame.Next().(int32))
				color.G = byte(frame.Next().(int32))
				color.B = byte(frame.Next().(int32))
				img.Set(x, y, color)
			}
		}

		out, _ := os.Create(fmt.Sprintf("temp/scalar%d.png", currFrame))
		encoder.Encode(out, img)
	}

	return nil
}

func (v *Scalar) HorizontalTimestack(numFrames int, xIndex int) error {
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
		x := (xIndex - 1) * 3
		for y := 0; y < v.height; y++ {
			color.R = frame.Buffer[x]
			color.G = frame.Buffer[x+1]
			color.B = frame.Buffer[x+2]
			img.Set(cursor, y, color)
			x += v.width * 3
		}
		cursor++
	}

	out, _ := os.Create(fmt.Sprintf("temp/hstack%d_%d.png", numFrames, xIndex))
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

/* Composite Volume *********************************************************/

type Point struct {
	X float64
	Y float64
}

type Composite struct {
	*Scalar
	*Vector
	LeftEnd, RightEnd Point
}

// b = y - mx
// m = (y2-y1)/(x2-x1)
// line length (d) = sqrt((x2-x1)^2+(y2-y1)^2)
// image top left is (0, 0)
// distance between buoys = line length/num buoys

func (v *Composite) Timeline(numFrames, numBuoys int) error {
	if err := v.Scalar.open(); err != nil {
		return err
	}
	defer v.Scalar.Close()
	if err := v.Vector.open(); err != nil {
		return err
	}
	defer v.Vector.Close()

	v_reader := csv.NewReader(bufio.NewReader(v.Vector.File))
	type temp struct {
		X int
		Y int
	}

	s_frame := framebuffer.New(v.Scalar.width, v.Scalar.height, SZ_UCHAR, 3)
	v_frame := framebuffer.New(v.Vector.width, v.Vector.height, SZ_UINT32, 2)

	m := (float64(v.RightEnd.Y) - float64(v.LeftEnd.Y)) / (float64(v.RightEnd.X) - float64(v.LeftEnd.X))
	b := float64(v.RightEnd.Y - m*v.RightEnd.X)
	dist := (v.RightEnd.X - v.LeftEnd.X) / (float64(numBuoys - 1))

	var points []Point
	for i := float64(0); int(i) < numBuoys; i++ {
		x := float64(v.LeftEnd.X) + dist*i
		points = append(points, Point{math.Round(x), math.Round(m*x + b)})
	}

	_color := color.RGBA{0, 0, 0, 0xFF}

	for currFrame := 0; currFrame < numFrames; currFrame++ {
		s_frame.Reset()
		v_frame.Reset()
		if _, err := v.Scalar.Read(s_frame.Buffer); err != nil {
			return err
		}

		/* if _, err := v.Vector.Read(v_frame.Buffer); err != nil {
			return err
		} */

		// TODO check read correct # of bytes

		img := v.Scalar.createImage()
		for y := 0; y < v.Scalar.height; y++ {
			for x := 0; x < v.Scalar.width; x++ {
				_color.R = byte(s_frame.Next().(int32))
				_color.G = byte(s_frame.Next().(int32))
				_color.B = byte(s_frame.Next().(int32))
				img.Set(x, y, _color)
			}
		}

		var pixels []temp
		for n := 0; n < v.Vector.height*v.Vector.width; n++ {
			px, _ := v_reader.Read()
			x, _ := strconv.Atoi(px[0])
			y, _ := strconv.Atoi(px[1])
			pixels = append(pixels, temp{x, y})
		}

		gc := draw2dimg.NewGraphicContext(img)
		gc2 := draw2dimg.NewGraphicContext(img)
		gc2.SetStrokeColor(color.RGBA{0xFF, 0, 0, 0xFF})
		gc2.SetLineWidth(4)
		gc.SetFillColor(color.RGBA{0xFF, 0, 0, 0xFF})
		gc.SetLineWidth(8)

		for i, old := range points {
			gc2.LineTo(old.X, old.Y)
			gc2.Stroke()
			gc2.MoveTo(old.X, old.Y)

			draw2dkit.Circle(gc, old.X, old.Y, 5)
			gc.Fill()

			point := &points[i]
			px := (int(point.X) + int(point.Y)*v.Vector.width)

			newX := pixels[px].X
			newY := pixels[px].Y
			if newX != 0 {
				point.X += float64(newX / 333333)
			}
			if newY != 0 {
				point.Y += float64(newY / 333333)
			}
		}

		//fmt.Printf("p %+v\n", points[0])

		out, _ := os.Create(fmt.Sprintf("temp/composite%d.png", currFrame))
		encoder.Encode(out, img)
	}

	return nil
}
