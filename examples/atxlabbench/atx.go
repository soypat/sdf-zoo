package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const (
	// Width
	atxW = 150.
	// height
	atxH           = 85.
	panelThickness = 4.
	bananaSpacing  = 20.
)

var (
	// Depending on material printed. 1.12 and 0.95 for PLA
	thruHole   = dimensionCorrector{multiplier: 1.12}
	threadHole = dimensionCorrector{multiplier: 0.95}

	onButton, _        = sdf.Circle2D(thruHole.apply(20 / 2))
	bananaPlugBig, _   = sdf.Circle2D(thruHole.apply(6.5 / 2))
	bananaPlugSmall, _ = sdf.Circle2D(2)

	voltageDisplay = sdf.Box2D(sdf.V2{45, 25.7}, 0)
)

func main() {

	panel, err := obj.Panel2D(&obj.PanelParms{
		Size:         sdf.V2{atxW, atxH},
		CornerRadius: 4,
		HoleDiameter: threadHole.apply(4),
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"xxx", "xxx", "xxx", "xxx"},
	})
	must(err)

	// Outputs banana plugs
	outputs := sdf.Array2D(bananaPlugBig, sdf.V2i{4, 2}, sdf.V2{-bananaSpacing, bananaSpacing})
	outputs = sdf.Transform2D(outputs, sdf.Translate2d(sdf.V2{atxW/2 - 15, -atxH/2 + 15}))

	panel = sdf.Difference2D(panel, sdf.Transform2D(onButton, sdf.Translate2d(sdf.V2{atxW/2 - 25, atxH/2 - 15})))
	panel = sdf.Difference2D(panel, outputs)

	// Begin working on regulated step-down block

	//	WIP

	model := sdf.Extrude3D(panel, panelThickness)
	render.RenderSTLSlow(model, 200, "atx_bench.stl")
}

type dimensionCorrector struct {
	multiplier float64
}

func (d dimensionCorrector) apply(dim float64) (corrected float64) { return d.multiplier * dim }

func must(err error) {
	if err != nil {
		panic(err)
	}
}
