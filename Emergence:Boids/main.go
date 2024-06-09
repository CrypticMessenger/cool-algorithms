package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WIDTH  int = 400
	HEIGHT int = 400
)

type Game struct {
	simulation *Simulation
}

// 0 -> WIDTH |
// |   Board  |
// v		  |
// HEIGHT	  |
func (g *Game) Update() error {
	g.simulation.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.simulation.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {

	var g *Game = &Game{
		simulation: NewSimulation(50, 1),
	}
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Emergence: Boids")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
