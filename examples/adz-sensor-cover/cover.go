package main

import (
	"fmt"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/soypat/sdf-zoo/helpers/correct"
)

func main() {
	// For ADZ Nagano pressure sensor type (SML-x)
	material := correct.PLA
	const (
		nozzleDiam           = 0.55
		adzDiam      float64 = 22
		connectorDim float64 = 15.6
		// Cover dimensions
		coverThick   float64 = 2 * nozzleDiam
		coverProtude float64 = 3
		round                = coverThick / 2
	)
	var dim float64
	// First create connector-facing part
	dim = connectorDim + coverThick*2
	cover, err := sdf.Box3D(sdf.V3{X: dim, Y: dim, Z: coverThick + coverProtude}, round)
	must(err)
	dim = material.InternalDimScale(connectorDim)
	empty, err := sdf.Box3D(sdf.V3{X: dim, Y: dim, Z: coverProtude}, 0)
	must(err)
	empty = sdf.Transform3D(empty, sdf.Translate3d(sdf.V3{Z: coverProtude / 2}))
	cover = sdf.Difference3D(cover, empty)

	// We now create sensor-facing part
	sensorCover, err := sdf.Cylinder3D(coverProtude+coverThick, adzDiam/2+coverThick, round)
	must(err)
	dim = material.InternalDimScale(adzDiam / 2)
	empty, err = sdf.Cylinder3D(coverProtude*2, dim, round)
	empty = sdf.Transform3D(empty, sdf.Translate3d(sdf.V3{Z: -coverProtude}))
	sensorCover = sdf.Difference3D(sensorCover, empty)
	sensorCover = sdf.Transform3D(sensorCover, sdf.Translate3d(sdf.V3{Z: -coverProtude}))
	cover = sdf.Union3D(cover, sensorCover)

	// Make hole for connector pins.
	dim = 10
	hole, err := sdf.Box3D(sdf.V3{X: dim, Y: dim, Z: 4 * coverThick}, 0)
	must(err)
	hole = sdf.Transform3D(hole, sdf.Translate3d(sdf.V3{Z: -2 * coverThick}))
	cover = sdf.Difference3D(cover, hole)
	cover = material.Scale(cover)
	render.RenderSTLSlow(cover, 250, fmt.Sprintf("cover.stl"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
