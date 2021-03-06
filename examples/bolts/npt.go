package main

import (
	"fmt"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/soypat/sdf-zoo/helpers/correct"
)

func main() {
	const (
		mmPerInch = 25.4
		tlen      = 15 / mmPerInch
		shank     = 3 / mmPerInch
		thread    = "npt_3/4"
		cells     = 200
	)
	bolt, err := obj.Bolt(&obj.BoltParms{Thread: thread, Style: "knurl", TotalLength: tlen, ShankLength: shank})
	must(err)
	bbz := bolt.BoundingBox().Size().Z
	bbx := bolt.BoundingBox().Size().X
	hollow, err := sdf.Cylinder3D(tlen+shank, bbx/3.5, 3/mmPerInch)
	hollow = sdf.Transform3D(hollow, sdf.Translate3d(sdf.V3{0, 0, bbz - tlen/2 - shank/2}))
	must(err)
	bolt = sdf.Difference3D(bolt, hollow)
	nut, err := NutPlug(&obj.NutParms{Thread: thread, Style: "hex"})
	must(err)

	nut = sdf.ScaleUniform3D(nut, mmPerInch)
	bolt = sdf.ScaleUniform3D(bolt, mmPerInch)
	render.RenderSTLSlow(correct.PLA.Scale(nut), cells, "npt_nut.stl")
	render.RenderSTLSlow(correct.PLA.Scale(bolt), cells, "npt_bolt.stl")
}

// must asserts there is no error. if error encountered terminate program
func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Nut returns a simple nut suitable for 3d printing.
func NutPlug(k *obj.NutParms) (sdf.SDF3, error) {
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
	nr := t.HexRadius() * 1.1
	nh := t.HexHeight()
	plugExtraHeight := nh * 0.2
	switch k.Style {
	case "hex":
		nut, err = obj.HexHead3D(nr, nh+plugExtraHeight, "tb")
	case "knurl":
		nut, err = obj.KnurledHead3D(nr, nh+plugExtraHeight, nr*0.25)
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
