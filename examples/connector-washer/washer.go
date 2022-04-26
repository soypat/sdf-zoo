package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/soypat/sdf-zoo/helpers/correct"
)

func main() {
	// For PLT connectors (microphone connector)
	material := correct.PLA
	const (
		panelHoleD      float64 = 16
		connectorD      float64 = 15.4
		washerHeight    float64 = 1.5
		washerRimHeight float64 = .7
		washerMainDiam          = 1.5 * panelHoleD
		washerRimDiam           = .99 * panelHoleD
	)
	var (
		washerHoleDiam = material.InternalDimScale(connectorD)
	)

	washer, err := sdf.Cylinder3D(washerHeight, washerMainDiam/2, 0)
	must(err)
	washerRim, err := sdf.Cylinder3D(washerHeight+washerRimHeight, washerRimDiam/2, 0)
	must(err)
	washerRim = sdf.Transform3D(washerRim, sdf.Translate3d(sdf.V3{Z: washerRimHeight / 2}))
	washer = sdf.Union3D(washerRim, washer)
	// Make hole
	washerHole, err := sdf.Cylinder3D(2*(washerHeight+washerRimHeight), washerHoleDiam/2, 0)
	must(err)
	washer = sdf.Difference3D(washer, washerHole)
	render.RenderSTLSlow(washer, 250, fmt.Sprintf("washer%gx%g.stl", math.Round(washerHoleDiam), math.Round(washerHeight)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
