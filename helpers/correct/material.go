package correct

import "github.com/deadsy/sdfx/sdf"

var (
	PLA = Material{shrink: 0.2e-2} // 0.2% shrinkage
)

type Material struct {
	shrink float64
}

func (m Material) Scale(s sdf.SDF3) sdf.SDF3 {
	scale := 1 / (1 - m.shrink)
	return sdf.ScaleUniform3D(s, scale)

}
