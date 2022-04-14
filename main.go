package main

import (
	"PA1/volumes"
)

// TODO try reading more than one frame into buffer at a time
// TODO add New() function for scalar and vector

var ManresaScalar = volumes.Scalar{
	Volume: volumes.NewVolume("manresa_scalar_512.raw", 1280, 720, 512),
}

var ManresaVector = volumes.Vector{
	Volume: volumes.NewVolume("manresa_scalar_512-001.raw", 1280, 720, 512),
}

/* func (v *VectorVolume) pullFrames(numFrames int) error {
	file, err := os.OpenFile(v.filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	frame := make([]byte, v.height*v.width*2)

	for curr := 0; curr < numFrames; curr++ {
		cursor := 0
		if _, err := file.Read(frame); err != nil {
			return err
		}

		// todo check read correct # of bytes
		img := image.NewRGBA(image.Rectangle{
			image.Point{0, 0},
			image.Point{v.width, v.height},
		})
		color := color.RGBA{0, 0, 0, 0xFF}
		for y := 0; y < v.height; y++ {
			for x := 0; x < v.width; x++ {
				color.R = frame[cursor]
				color.G = frame[cursor+1]
				cursor += 2
				img.Set(x, y, color)
			}
		}

		out, _ := os.Create(fmt.Sprintf("temp/vector%d.png", curr))
		png.Encode(out, img) // todo check error
	}

	return nil
} */

func main() {
	/* for x := 0; x < 5; x++ {
		start := time.Now()
		fmt.Println(time.Since(start).Seconds())

		ManresaScalar.pullFrames(10)
	} */

	//ManresaScalar.AverageFrames(512)
	//ManresaScalar.PullFrames(50)
	ManresaScalar.HorizontalTimestack(512)
	//ManresaVector.pullFrames(50)
}
