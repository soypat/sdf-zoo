package main

import (
	"fmt"
	"sdftest/helpers/hio"
	"sdftest/helpers/spacing"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/golang/freetype/truetype"
)

func main() {
	const (
		boxThick     = 3.0
		letterHeight = 4.0
		solidmode    = true
	)
	var boxHeight, boxLength float64
	var panel, holes, text3d sdf.SDF3
	diams := []float64{1, 1.5, 2, 2.5, 3, 4, 5, 6, 8}
	var holeSpacing []sdf.V3
	// Create hole cylinders
	{
		sholes := make([]sdf.SDF3, len(diams))
		for i, d := range diams {
			c, err := sdf.Circle2D(d)
			must(err)
			sholes[i] = sdf.Extrude3D(c, 5)
		}
		holeSpacing = spacing.Directional(sholes, sdf.V3{1, 0, 0}, 1)
		holes = sdf.Union3D(sholes...)
		holes = sdf.Transform3D(holes, sdf.Translate3d(sdf.V3{0, letterHeight / 2, 0}))
		if solidmode {
			render.RenderSTLSlow(holes, 400, "cylinders.stl")
		}
	}
	// create box
	{
		// create box
		bb := holes.BoundingBox()
		sz := bb.Size()
		sq := sdf.Box2D(sdf.V2{X: sz.X + 5, Y: sz.Y + 5 + letterHeight}, 4)
		sq = sdf.Transform2D(sq, sdf.Translate2d(sdf.V2{35, 0}))
		panel = sdf.Extrude3D(sq, boxThick)
		boxHeight, boxLength = panel.BoundingBox().Size().Y, panel.BoundingBox().Size().X
	}

	// Create Text
	{
		ttf, err := truetype.Parse(hio.FontJetBrainsMonoBold)
		must(err)
		var texts []sdf.SDF3
		for i := range diams {
			txtobj := sdf.NewText(fmt.Sprintf("%g", diams[i]))
			textsdf, _ := sdf.TextSDF2(ttf, txtobj, 4)
			sdf.Extrude3D(textsdf, 2)
			pos := holeSpacing[i]
			if i%2 == 0 { // stagger hole labels
				pos.Y -= letterHeight
			}
			texts = append(texts, sdf.Transform3D(sdf.Extrude3D(textsdf, 2), sdf.Translate3d(pos)))

		}
		text3d = sdf.Union3D(texts...)
		must(err)
		text3d = sdf.Transform3D(text3d, sdf.Translate3d(sdf.V3{boxLength * 0, -boxHeight/2 + letterHeight*1.6, boxThick / 2}))
	}

	result := sdf.Difference3D(panel, holes)
	result = sdf.Difference3D(result, text3d)
	render.RenderSTLSlow(result, 400, "hole_panel.stl")

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
