package spacing

import (
	"github.com/deadsy/sdfx/sdf"
)

// Spaces sdf3's located at one point.
func Directional(shapes []sdf.SDF3, dir sdf.V3, baseSep float64) (spacings []sdf.V3) {
	if dir.Equals(sdf.V3{}, 0) {
		panic("direction must not be zero")
	}
	spacings = make([]sdf.V3, len(shapes))
	norm := dir.Normalize()
	sep := norm.MulScalar(baseSep)
	for i := 1; i < len(shapes); i++ {
		prevCenter := shapes[i-1].BoundingBox().Center()
		box := shapes[i].BoundingBox().Size()
		boxPrev := shapes[i-1].BoundingBox().Size()
		spacings[i] = norm.MulScalar((box.Dot(norm) + boxPrev.Dot(norm)) / 2)
		spacings[i] = spacings[i].Add(sep.Add(prevCenter))
		shapes[i] = sdf.Transform3D(shapes[i], sdf.Translate3d(spacings[i]))
	}
	return spacings
}
