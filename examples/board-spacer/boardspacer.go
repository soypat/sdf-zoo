package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/soypat/sdf-zoo/helpers/correct"
)

const (
	baseWidth     = 100.0
	baseLength    = 100.0
	baseThickness = 2.4

	frontPanelThickness = 3.0
	frontPanelLength    = 170.0
	frontPanelHeight    = 50.0
	frontPanelYOffset   = 15.0

	holeWidth    = 4.0
	pillarHeight = 7
)

var material = correct.PLA

func main() {
	b := base()
	render.RenderSTL(b, 200, "base.stl")
}

// base returns the base mount.
func base() sdf.SDF3 {
	// base
	pp := &obj.PanelParms{
		Size:         sdf.V2{baseLength, baseWidth},
		CornerRadius: holeWidth * 1.2,
		HoleDiameter: material.InternalDimScale(holeWidth),
		HoleMargin:   [4]float64{4.5, 4.5, 4.5, 4.5},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0, err := obj.Panel2D(pp)

	if err != nil {
		panic(err)
	}

	s2 := sdf.Extrude3D(s0, baseThickness)
	xOfs := 0.5 * baseLength
	yOfs := 0.5 * baseWidth
	s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{xOfs, yOfs, 0}))

	// standoffs
	zOfs := 0.5 * (pillarHeight + baseThickness)
	m4Positions := sdf.V3Set{
		{4.5, 4.5, zOfs},
		{4.5, 95.5, zOfs},
		{95.5, 95.5, zOfs},
		{95.5, 4.5, zOfs},
		{60, 30, zOfs},
		{60, 70, zOfs},
	}
	m4Standoffs := standoffs(4, m4Positions)
	m3Positions := sdf.V3Set{
		{9, 35.5, zOfs},
		{9, 62.5, zOfs},
		{91, 64.5, zOfs},
		{91, 37.5, zOfs},
	}
	m3Standoffs := standoffs(3, m3Positions)
	s4 := sdf.Union3D(s2, m4Standoffs, m3Standoffs)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4
}

// multiple standoffs
func standoffs(holeWidth float64, positions sdf.V3Set) sdf.SDF3 {
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: holeWidth * 2,
		HoleDepth:      pillarHeight + baseThickness,
		HoleDiameter:   material.InternalDimScale(holeWidth),
	}

	// from the board mechanicals

	s, err := obj.Standoff3D(k)
	if err != nil {
		panic(err)
	}
	return sdf.Multi3D(s, positions)
}
