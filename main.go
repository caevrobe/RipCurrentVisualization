package main

import (
	"PA1/volumes"
	"log"
)

// TODO try reading more than one frame into buffer at a time
// TODO add New() function for scalar and vector

// buoy line: (150, 675) to (1200, 330)
var ManresaScalar = volumes.Scalar{
	Volume: volumes.NewVolume("TEST_DATA/Manresa/manresa_scalar_512.raw", 1280, 720, 512),
}

var ManresaVector = volumes.Vector{
	Volume: volumes.NewVolume("TEST_DATA/Manresa/manresa_vector_512-001.raw", 1280, 720, 512),
}

var ManresaComposite = volumes.Composite{
	Scalar: &ManresaScalar,
	Vector: &ManresaVector,
}

var SeabrightScalar = volumes.Scalar{
	Volume: volumes.NewVolume("TEST_DATA/Seabright/seabright_scalar.raw", 1280, 720, 512),
}

// buoy line: (100, 485) to (1200, 410)
var SeabrightVector = volumes.Vector{
	Volume: volumes.NewVolume("TEST_DATA/Seabright/seabright_vector.raw", 1280, 720, 512),
}

func main() {
	//SeabrightScalar.PullFrames(512)

	for x := 1; x <= 1280; x += 64 {
		if err := ManresaScalar.HorizontalTimestack(512, x); err != nil {
			log.Fatal(err)
		}
	}

	//fmt.Println(SeabrightScalar.PullFrames(1))
	//fmt.Println(SeabrightScalar.AverageFrames(50))

	//fmt.Println(ManresaScalar.PullFrames(1))
	//fmt.Println(ManresaScalar.HorizontalTimestack(512))
	//fmt.Println(ManresaScalar.AverageFrames(512))
	//fmt.Println(SeabrightScalar.AverageFrames(512))

	//fmt.Println(ManresaVector.PullFrames(1))
	//ManresaScalar.AverageFrames(512)
	//ManresaScalar.AverageFrames2(512)

	/* ManresaScalar.AverageFrames(512)
	ManresaScalar.PullFrames(50) */
	//ManresaScalar.HorizontalTimestack(512)
	//ManresaVector.pullFrames(50)
}
