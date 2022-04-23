package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	const (
		// Calculate candle holder internal volume specifications
		volRatio     = 1.05
		candleDiam   = 35.0 // mm
		candleHeight = 50.0
		holderDiam   = candleDiam + 4.0 // 4mm extra diameter for eazy fitting of candle.
		// solve for holderHeight
		// ratio*pi*(Dc^2)*Hc/4 = pi*(Dh^2)*Hh/4
		// Hh = ratio*(Dc/Dh)^2*Hc
		diamRatio    = candleDiam / holderDiam
		holderHeight = volRatio * diamRatio * diamRatio * candleHeight
		// Candle Holder mechanical specifications
		holderWallThick = 3.0
		holderBaseThick = 5.0
	)
	holder, _ := sdf.Cylinder3D(holderBaseThick+holderHeight, holderDiam/2+holderWallThick, 2)
	hole, _ := sdf.Cylinder3D(holderHeight+1, holderDiam/2, 1)
	hole = sdf.Transform3D(hole, sdf.Translate3d(sdf.V3{Z: holderBaseThick / 2}))
	holder = sdf.Difference3D(holder, hole)
	render.RenderSTL(holder, 100, "candleholder.stl")
}
