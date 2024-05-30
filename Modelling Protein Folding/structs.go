package main

type PotentialFunctionConfig struct {
	epsilon float64
	sigma   float64
}

type Context struct {
	potentialFunctionConfig *PotentialFunctionConfig
	learningRate            float64
	numIterations           int
	numPoints               int
	numDimensions           int
	inputPoints             [][]float64
}
