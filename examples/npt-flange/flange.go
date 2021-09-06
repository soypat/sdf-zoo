package main

import (
	"fmt"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const (
	// thread length
	tlen             = 18 / 25.4
	internalDiameter = 1.5 / 2.
	flangeH          = 7 / 25.4
	flangeD          = 60. / 25.4
	thread           = "npt_1/2"
	// internal diameter scaling.
	plaScale = 1.03
)

func main() {
	pipe, err := ThreadedPipe(&obj.NutParms{Thread: thread, Style: "pipe"})
	must(err)
	// PLA scaling to thread
	pipe = sdf.Transform3D(pipe, sdf.Scale3d(sdf.V3{plaScale, plaScale, 1}))
	flange, err := sdf.Cylinder3D(flangeH, flangeD/2, flangeH/8)
	must(err)
	hole, err := sdf.Cylinder3D(flangeH, internalDiameter/2, 0)
	must(err)
	flange = sdf.Difference3D(flange, hole)
	flange = sdf.Transform3D(flange, sdf.Translate3d(sdf.V3{0, 0, -tlen / 2}))
	pipe = sdf.Union3D(pipe, flange)

	render.RenderSTLSlow(pipe, 300, "npt_flange.stl")

}

func ThreadedPipe(k *obj.NutParms) (sdf.SDF3, error) {
	// validate parameters
	t, err := sdf.ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.Tolerance < 0 {
		return nil, sdf.ErrMsg("Tolerance < 0")
	}

	// nut body
	var nut sdf.SDF3
	nr := t.HexRadius()
	nh := t.HexHeight()
	plugExtraHeight := 0.
	switch k.Style {
	case "hex":
		nut, err = obj.HexHead3D(nr, nh+plugExtraHeight, "tb")
	case "knurl":
		nut, err = obj.KnurledHead3D(nr, nh+plugExtraHeight, nr*0.25)
	case "pipe":
		nut, err = sdf.Cylinder3D(nh+plugExtraHeight, nr*1.1, 0)
	default:
		return nil, sdf.ErrMsg(fmt.Sprintf("unknown style \"%s\"", k.Style))
	}
	if err != nil {
		return nil, err
	}
	nut = sdf.Transform3D(nut, sdf.Translate3d(sdf.V3{0, 0, plugExtraHeight / 2}))
	// internal thread
	isoThread, err := sdf.ISOThread(t.Radius+k.Tolerance, t.Pitch, false)
	if err != nil {
		return nil, err
	}
	thread, err := sdf.Screw3D(isoThread, nh, t.Taper, t.Pitch, 1)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(nut, thread), nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
