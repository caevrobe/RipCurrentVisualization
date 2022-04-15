package main

import (
	"PA1/volumes"
	"fmt"
)

// TODO try reading more than one frame into buffer at a time
// TODO add New() function for scalar and vector

var ManresaScalar = volumes.Scalar{
	Volume: volumes.NewVolume("manresa_scalar_512.raw", 1280, 720, 512),
}

var ManresaVector = volumes.Vector{
	Volume: volumes.NewVolume("manresa_vector_512-001.raw", 1280, 720, 512),
}

func main() {
	//fmt.Println(ManresaScalar.PullFrames(512))
	//fmt.Println(ManresaScalar.HorizontalTimestack(512))
	fmt.Println(ManresaScalar.AverageFrames(512))

	//fmt.Println(ManresaVector.PullFrames(1))
	//ManresaScalar.AverageFrames(512)
	//ManresaScalar.AverageFrames2(512)

	/* ManresaScalar.AverageFrames(512)
	ManresaScalar.PullFrames(50) */
	//ManresaScalar.HorizontalTimestack(512)
	//ManresaVector.pullFrames(50)
}
