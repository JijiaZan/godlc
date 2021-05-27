package utils

import "github.com/JijiaZan/godml/pyserver"

func CalGradient(a *pyserver.Layer, b *pyserver.Layer) *pyserver.Layer{

	gradient := &pyserver.Layer{}
	gradient.Name = a.GetName()
	gradient.DimBias = a.GetDimBias()
	gradient.Bias = calArray(a.GetBias(), b.GetBias())
	gradient.DimKernel = a.GetDimKernel()
	gradient.Kernel = calArray(a.GetKernel(), b.GetKernel())
	return gradient
}

func calArray(a []float32, b []float32) []float32{
	res := make([]float32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = a[i] - b[i]
	}
	return res
}