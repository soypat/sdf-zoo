package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	const (
		spacerHeight     float64 = 5 // mm
		holeDiam         float64 = 3 // mm
		plaHoleCorrected         = holeDiam*1.05 + .4
	)
	ftf := math.Ceil(plaHoleCorrected*1.4) - .15 // Face to face hex distance
	hexRadius := getHexRadiusFromFTF(ftf)
	fmt.Printf("FTF:%.2g, holeCorrected:%.1g\n", ftf, plaHoleCorrected)
	sp, err := obj.HexHead3D(hexRadius, spacerHeight, "")
	must(err)
	hole, err := sdf.Cylinder3D(spacerHeight, plaHoleCorrected/2, 0)
	must(err)
	sp = sdf.Difference3D(sp, hole)
	render.RenderSTLSlow(sp, 100, fmt.Sprintf("spacer%gx%g.stl", math.Round(holeDiam), math.Round(spacerHeight)))
}

func getHexRadiusFromFTF(ftf float64) (radius float64) {
	return ftf / math.Cos(30.*math.Pi/180.) / 2
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
