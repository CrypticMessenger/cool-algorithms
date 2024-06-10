package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// statistics
var framecount int = 0
var generation int = 0
var AllGenerationBestFitness float64 = 0
var bestRocket *Rocket = nil

// parameters of the simulation
const (
	WIDTH               int        = 400
	HEIGHT              int        = 400
	DT                  float64    = 0.3
	POPULATION          int        = 50
	MUTATION_RATE       float64    = 0.2
	ROCKET_MASS         float64    = 1
	MAX_FORCE           float64    = 1.5
	GENERATION_LIFETIME int        = 200
	ERROR               float64    = 2
	GENERATION_LIMIT    int        = 50
	target              complex128 = complex(float64(WIDTH/2), 50)
	initialPos          complex128 = complex(50, 100)
)

// fitness function
func fitnessFunction(r *Rocket) float64 {
	return math.Pow(1/cmplx.Abs(r.Pos-target), 2)
}

// Selection function
func (s *Simulation) Selection() []DNA {

	var matingPoolSize int = 1000
	var matingPool []DNA = make([]DNA, 0, matingPoolSize)
	var topKDistribution []float64 = []float64{0.7, 0.2, 0.1}
	var topK []Pair = make([]Pair, 0, POPULATION)
	for i := 0; i < POPULATION; i++ {
		topK = append(topK, Pair{i, s.Rockets[i].Fitness})
	}
	// sort topK by fitness
	for i := 0; i < POPULATION; i++ {
		for j := i + 1; j < POPULATION; j++ {
			if topK[i].Second < topK[j].Second {
				topK[i], topK[j] = topK[j], topK[i]
			}
		}
	}
	// get mating pool
	for i := 0; i < POPULATION; i++ {
		for j := 0; j < len(topKDistribution); j++ {
			for k := 0; k < int(topKDistribution[j]*float64(matingPoolSize)); k++ {
				matingPool = append(matingPool, *s.Rockets[topK[i].First].Dna)
			}
		}
	}

	return matingPool
}

// Reproduction function
func (s *Simulation) Reproduce(matingPool []DNA) [POPULATION]*Rocket {
	var newRockets [POPULATION]*Rocket
	// Reproduction
	for i := 0; i < POPULATION; i++ {
		var parentA DNA = matingPool[int(rand.Float64()*float64(len(matingPool)))]
		var parentB DNA = matingPool[int(rand.Float64()*float64(len(matingPool)))]
		var childDNA DNA = DNA{}
		// CrossOver
		for j := 0; j < GENERATION_LIFETIME; j++ {
			if rand.Float64() < 0.5 {
				childDNA.Genes[j] = parentA.Genes[j]
			} else {
				childDNA.Genes[j] = parentB.Genes[j]
			}
			// Mutation
			if rand.Float64() < MUTATION_RATE {
				childDNA.Genes[j] = complex((rand.Float64()*2 - 0.95), (rand.Float64()*2 - 0.95))
				childDNA.Genes[j] /= complex(cmplx.Abs(childDNA.Genes[j]), 0)
				childDNA.Genes[j] *= complex(MAX_FORCE*rand.Float64(), 0)
			}
		}
		newRockets[i] = &Rocket{
			Fitness: 0,
			Dna:     &childDNA,
			Mass:    ROCKET_MASS,
			Pos:     initialPos,
			Vel:     0,
			Acc:     0,
		}
	}
	return newRockets
}

type Game struct {
	simulation *Simulation
}

// 0 -> WIDTH |
// |   Board  |
// v		  |
// HEIGHT	  |
func (g *Game) Update() error {
	for i := 0; i < POPULATION; i++ {
		g.simulation.Rockets[i].ApplyForce(g.simulation.Rockets[i].Dna.Genes[framecount%GENERATION_LIFETIME])
		g.simulation.Rockets[i].Update()
		g.simulation.Rockets[i].Fitness = fitnessFunction(g.simulation.Rockets[i])

	}
	framecount++
	if framecount == GENERATION_LIFETIME {

		var matingPool []DNA
		var newRockets [POPULATION]*Rocket

		matingPool = g.simulation.Selection()           // Selection
		newRockets = g.simulation.Reproduce(matingPool) // consists of crossover and mutation
		g.simulation.Rockets = newRockets

		generation++
		framecount = framecount % GENERATION_LIFETIME
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// if generation limit reached, then stop
	if generation > GENERATION_LIMIT {
		ebitenutil.DebugPrintAt(screen, "Generation Limit Reached", WIDTH/2, HEIGHT/2)
		fmt.Println("Target Reached")
		fmt.Println((*bestRocket).Dna.Genes)
		os.Exit(0)
		return
	}

	// if within error range of target, then stop
	if bestRocket != nil && cmplx.Abs(bestRocket.Pos-target) < ERROR {
		ebitenutil.DebugPrintAt(screen, "Target Reached", WIDTH/2, HEIGHT/2)
		fmt.Println("Target Reached")
		fmt.Println((*bestRocket).Dna.Genes)
		os.Exit(0)
		return
	}

	// display information
	ebitenutil.DebugPrintAt(screen, "Frame: "+fmt.Sprint(framecount), 0, 5)
	ebitenutil.DebugPrintAt(screen, "Generation: "+fmt.Sprint(generation), WIDTH-100, 5)

	// get best rocket and fitness in the generation
	var best_fitness float64 = -1
	var best_rocket_ind int = 0
	for i := 0; i < POPULATION; i++ {
		if g.simulation.Rockets[i].Fitness > best_fitness {
			best_fitness = g.simulation.Rockets[i].Fitness
			best_rocket_ind = i
		}
	}

	// update global best fitness
	ebitenutil.DebugPrintAt(screen, "best_fitness: "+fmt.Sprint(best_fitness), 0, 20)
	if best_fitness > AllGenerationBestFitness {
		AllGenerationBestFitness = best_fitness
		bestRocket = g.simulation.Rockets[best_rocket_ind]
	}

	// display all generation best fitness
	ebitenutil.DebugPrintAt(screen, "All generation best fitness: "+fmt.Sprint(AllGenerationBestFitness), 0, HEIGHT-20)

	// draw src,target and rockets
	screen.Set(int(real(initialPos)), int(imag(initialPos)), color.RGBA{0, 0, 255, 255})
	screen.Set(int(real(target)), int(imag(target)), color.RGBA{255, 0, 0, 255})
	for i := 0; i < POPULATION; i++ {
		if i == best_rocket_ind {
			g.simulation.Rockets[i].Draw(screen, color.RGBA{0, 255, 0, 255})
		} else {
			g.simulation.Rockets[i].Draw(screen, color.RGBA{255, 255, 255, 255})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {

	var g *Game = &Game{simulation: StartSimulation()}
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Genetic Algorithm: rockets")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
