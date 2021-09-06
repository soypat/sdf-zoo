package main

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	const (
		tlen   = .6
		shank  = 0.1
		thread = "npt_1/2"
	)
	// s1, err := nutAndBolt("npt_1/2", 10, 1)
	s1, err := obj.Bolt(&obj.BoltParms{
		Thread:      thread,
		Style:       "hex",
		Tolerance:   0.01,
		TotalLength: tlen + shank,
		ShankLength: shank,
	})
	t, _ := sdf.ThreadLookup(thread)

	holeRadius := math.Min(t.Radius/2, math.Max(t.Radius-4, 1))
	hole, _ := sdf.Cylinder3D(tlen, holeRadius, 0)
	print(t.HexHeight(), holeRadius)
	hole = sdf.Transform3D(hole, sdf.Translate3d(sdf.V3{0, 0, t.HexHeight() + shank + .2}))
	sresult := sdf.Difference3D(s1, hole)
	if err != nil {
		panic(err)
	}
	render.RenderSTLSlow(sresult, 100, "nutandbolt.stl")
}

func nutAndBolt(
	name string, // name of thread
	totalLength float64, // threaded length + shank length
	shankLength float64, //  non threaded length
) (sdf.SDF3, error) {

	// bolt
	boltParms := obj.BoltParms{
		Thread:      name,
		Style:       "hex",
		TotalLength: totalLength,
		ShankLength: shankLength,
	}
	bolt, err := obj.Bolt(&boltParms)
	if err != nil {
		return nil, err
	}

	// nut
	nutParms := obj.NutParms{
		Thread: name,
		Style:  "hex",
	}
	nut, err := obj.Nut(&nutParms)
	if err != nil {
		return nil, err
	}

	zOffset := totalLength * 1.5
	nut = sdf.Transform3D(nut, sdf.Translate3d(sdf.V3{0, 0, zOffset}))

	return sdf.Union3D(nut, bolt), nil
}
