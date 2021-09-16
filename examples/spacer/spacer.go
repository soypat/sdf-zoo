package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	holeDiam := 5.0 // mm
	ftfDiam := holeDiam * 1.8
	spacerHeight := 27.0 // mm
	sp, err := obj.HexHead3D(ftfDiam/2, spacerHeight, "")
	must(err)
	hole, err := sdf.Cylinder3D(spacerHeight, holeDiam/2, 0)
	must(err)
	sp = sdf.Difference3D(sp, hole)
	render.RenderSTLSlow(sp, 200, "spacer.stl")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
