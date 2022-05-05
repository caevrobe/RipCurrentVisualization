package main

import (
	"PA1/volumes"
)

// TODO try reading more than one frame into buffer at a time
// TODO add New() function for scalar and vector

// buoy line: (150, 675) to (1200, 330)
var ManresaScalar = volumes.Scalar{
	Volume: volumes.NewVolume("TEST_DATA/Manresa/manresa_scalar_512.raw", 1280, 720, 512),
}

var ManresaVector = volumes.Vector{
	Volume: volumes.NewVolume("TEST_DATA/Manresa/vector.csv", 1280, 720, 512),
	//Volume: volumes.NewVolume("TEST_DATA/Manresa/manresa_vector_512-001.raw", 1280, 720, 512),
}

var ManresaComposite = volumes.Composite{
	Scalar:   &ManresaScalar,
	Vector:   &ManresaVector,
	LeftEnd:  volumes.Point{150, 675},
	RightEnd: volumes.Point{1200, 330},
	/* LeftEnd:  volumes.Point{50, 600},
	RightEnd: volumes.Point{1100, 300}, */
	/* LeftEnd:  volumes.Point{50, 700},
	RightEnd: volumes.Point{1100, 350}, */
}

var SeabrightScalar = volumes.Scalar{
	Volume: volumes.NewVolume("TEST_DATA/Seabright/seabright_scalar.raw", 1280, 720, 512),
}

// buoy line: (100, 485) to (1200, 410)
var SeabrightVector = volumes.Vector{
	Volume: volumes.NewVolume("TEST_DATA/Seabright/vector.csv", 1280, 720, 512),
	//Volume: volumes.NewVolume("TEST_DATA/Seabright/seabright_vector.raw", 1280, 720, 512),
}

var SeabrightComposite = volumes.Composite{
	Scalar: &SeabrightScalar,
	Vector: &SeabrightVector,
	/* LeftEnd:  volumes.Point{100, 485},
	RightEnd: volumes.Point{1200, 410}, */
	LeftEnd:  volumes.Point{100, 350},
	RightEnd: volumes.Point{1200, 300},
}

func main() {

	// manresa timestack
	/* for x := 1; x <= 1280; x += 64 {
		if err := ManresaScalar.HorizontalTimestack(512, x); err != nil {
			log.Fatal(err)
		}
	} */

	// seabright timestack
	/* for x := 1; x <= 1280; x += 64 {
		if err := SeabrightScalar.HorizontalTimestack(512, x); err != nil {
			log.Fatal(err)
		}
	} */

	// manresa timeline
	//ManresaComposite.Timeline(127, 100)

	// seabright timeline
	//SeabrightComposite.Timeline(512, 100)

	// manresa average
	//ManresaScalar.AverageFrames(512)

	// seabright average
	//SeabrightScalar.AverageFrames(512)
}
