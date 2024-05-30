package main

import "math/rand"

func newPotentialFunctionConfig() *PotentialFunctionConfig {
	return &PotentialFunctionConfig{
		epsilon: 1,
		sigma:   1,
	}
}
func getRandomPoints(numPoints int, numDimensions int, generator *rand.Rand) [][]float64 {
	var points [][]float64 = make([][]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = make([]float64, numDimensions)
		for j := 0; j < numDimensions; j++ {
			points[i][j] = generator.Float64()
		}
	}
	return points
}
func setupContext(numPoints int, numDimensions int, learningRate float64, numIterations int, generator *rand.Rand) *Context {
	return &Context{
		potentialFunctionConfig: newPotentialFunctionConfig(),
		learningRate:            learningRate,
		numIterations:           numIterations,
		numPoints:               numPoints,
		numDimensions:           numDimensions,
		inputPoints:             getRandomPoints(numPoints, numDimensions, generator),
	}
}
