package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const (
	// Gabinete has 140x100 (widthxheight) space available
	// Width,  panelH
	panelW, panelH                = 100., 70
	holeDiam                      = 4.0
	panelThickness                = 4.
	bananaSpacing                 = 20.
	vDispH, vDispW                = 25.7, 45
	regBlockDepth, regBlockMargin = 4., 4.
)

var (
	// Depending on material printed. 1.12 and 0.95 for PLA
	thruHole   = dimensionCorrector{multiplier: 1.12}
	threadHole = dimensionCorrector{multiplier: 0.95}

	onButton, _        = sdf.Circle2D(thruHole.apply(20. / 2))
	bananaPlugBig, _   = sdf.Circle2D(thruHole.apply(6.5 / 2))
	bananaPlugSmall, _ = sdf.Circle2D(thruHole.apply(6. / 2))
)

func main() {

	panel, err := obj.Panel2D(&obj.PanelParms{
		Size:         sdf.V2{panelW, panelH},
		CornerRadius: .8 * holeDiam,
		HoleDiameter: threadHole.apply(holeDiam),
		HoleMargin:   [4]float64{1.5 * holeDiam, 1.5 * holeDiam, 1.5 * holeDiam, 1.5 * holeDiam},
		HolePattern:  [4]string{"xxx", "xxx", "xxx", "xxx"},
	})
	must(err)

	// LED Switch board buttons
	// | ALL ON |5V Periph|
	// | 12V BUS| 5V OBC  |
	const onSpacingMult = 1.4
	onButtons := sdf.Array2D(onButton, sdf.V2i{2, 2}, sdf.V2{20 * onSpacingMult, 20 * onSpacingMult})
	onButtons = sdf.Transform2D(onButtons, sdf.Translate2d(sdf.V2{-panelW/2 + 20, -panelH/2 + 20}))
	// Outputs banana plugs
	outputs := sdf.Array2D(bananaPlugBig, sdf.V2i{2, 3}, sdf.V2{bananaSpacing, bananaSpacing})
	outputs = sdf.Transform2D(outputs, sdf.Translate2d(sdf.V2{18, -panelH/2 + 15}))

	panel = sdf.Difference2D(panel, onButtons)
	panel = sdf.Difference2D(panel, outputs)

	// Generate model
	model := sdf.Extrude3D(panel, panelThickness)
	render.RenderSTLSlow(model, 200, "ledswitchboard.stl")
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
