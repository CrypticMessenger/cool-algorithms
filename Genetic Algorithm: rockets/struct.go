package main

import (
	"image/color"
	"math/cmplx"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DNA struct {
	Genes [GENERATION_LIFETIME]complex128
}

func (d *DNA) getNewDNA() {
	for i := 0; i < GENERATION_LIFETIME; i++ {
		d.Genes[i] = complex((rand.Float64()*2 - 0.95), (rand.Float64()*2 - 0.95))
		d.Genes[i] /= complex(cmplx.Abs(d.Genes[i]), 0)
		d.Genes[i] *= complex(MAX_FORCE*rand.Float64(), 0)
	}
}

type Rocket struct {
	Dna     *DNA
	Fitness float64
	Mass    float64
	Pos     complex128
	Vel     complex128
	Acc     complex128
}

type Simulation struct {
	Rockets [POPULATION]*Rocket
}

type Pair struct {
	First  int
	Second float64
}



func StartSimulation() *Simulation {
	var rockets [POPULATION]*Rocket

	for i := 0; i < POPULATION; i++ {
		var dna *DNA = &DNA{}
		dna.getNewDNA()
		rockets[i] = &Rocket{
			Fitness: 0,
			Dna:     dna,
			Mass:    ROCKET_MASS,
			Pos:     initialPos,
			Vel:     0,
			Acc:     0,
		}
	}
	return &Simulation{Rockets: rockets}
}

func (r *Rocket) ApplyForce(f complex128) {
	r.Acc += f / complex(r.Mass, 0)
}

func (r *Rocket) Update() {
	r.Vel += r.Acc * complex(DT, 0)
	r.Pos += r.Vel * complex(DT, 0)
	r.Acc = 0
}

func (r *Rocket) Draw(screen *ebiten.Image, color color.Color) {
	ebitenutil.DrawCircle(screen, real(r.Pos), imag(r.Pos), 2, color)
}
