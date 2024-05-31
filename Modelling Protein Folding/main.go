package main

import (
	"fmt"
	"math"
	"math/rand"
)

func lennardJonesPotential(r float64, config *PotentialFunctionConfig) float64 {
	var sigma float64 = config.sigma
	var epsilon float64 = config.epsilon
	return 4 * epsilon * (math.Pow(sigma/r, 12) - math.Pow(sigma/r, 6))
}

func lennardJonesPotentialGradient(r float64, config *PotentialFunctionConfig) float64 {
	var sigma float64 = config.sigma
	var epsilon float64 = config.epsilon
	return -24 * epsilon * (math.Pow(sigma, 6)/math.Pow(r, 8) - 2*math.Pow(sigma, 12)/math.Pow(r, 14))
}

func getDistance(point1 []float64, point2 []float64) float64 {
	var distance float64 = 0
	for i := 0; i < len(point1); i++ {
		distance += math.Pow(point1[i]-point2[i], 2)
	}
	return math.Sqrt(distance)
}

func totalPotentialEnergy(inputPoints [][]float64, config *PotentialFunctionConfig) float64 {
	var numPoints int = len(inputPoints)
	var totalEnergy float64 = 0
	for i := 0; i < numPoints; i++ {
		for j := i + 1; j < numPoints; j++ {
			var distance float64 = getDistance(inputPoints[i], inputPoints[j])
			totalEnergy += lennardJonesPotential(distance, config)
		}
	}
	return totalEnergy
}

func getGradientwrt(index int, inputPoints [][]float64, config *PotentialFunctionConfig) []float64 {
	var numPoints int = len(inputPoints)
	var numDimensions int = len(inputPoints[0])
	var gradient []float64 = make([]float64, numDimensions)
	for i := 0; i < numPoints; i++ {
		if i == index {
			continue
		}
		var distance float64 = getDistance(inputPoints[index], inputPoints[i])

		var lj_force float64 = lennardJonesPotentialGradient(distance, config)
		for j := 0; j < numDimensions; j++ {
			gradient[j] += lj_force * (inputPoints[index][j] - inputPoints[i][j])
		}
	}
	return gradient
}

func gradientDescent(context *Context) {
	var numPoints int = len(context.inputPoints)
	var numDimensions int = len(context.inputPoints[0])
	for i := 0; i < context.numIterations; i++ {
		for j := 0; j < numPoints; j++ {
			var gradient []float64 = getGradientwrt(j, context.inputPoints, context.potentialFunctionConfig)
			for k := 0; k < numDimensions; k++ {
				context.inputPoints[j][k] += context.learningRate * gradient[k]
			}
		}
		if (i+1)%50 == 0 {
			fmt.Println("Iteration:", i+1, "Potential energy:", totalPotentialEnergy(context.inputPoints, context.potentialFunctionConfig))
		}
	}
}

func main() {
	source := rand.NewSource(42)
	generator := rand.New(source)

	var context *Context = setupContext(
		4,         //  points
		1,         //  dimensions
		0.001,     // learning rate
		3000,       // iterations
		generator, // random number generator
	)

	fmt.Println("Initial points:")
	for i := 0; i < context.numPoints; i++ {
		fmt.Println(context.inputPoints[i])
	}

	fmt.Println("Initial total potential energy:", totalPotentialEnergy(context.inputPoints, context.potentialFunctionConfig))

	// perform gradient descent
	gradientDescent(context)

	fmt.Println("\nOutput points:")
	for i := 0; i < context.numPoints; i++ {
		fmt.Println(context.inputPoints[i])
	}
	fmt.Println("Final total potential energy:", totalPotentialEnergy(context.inputPoints, context.potentialFunctionConfig))

}
