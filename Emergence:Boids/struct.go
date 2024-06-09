package main

import (
	"image/color"
	"math/cmplx"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Boid struct {
	Position complex128
	Velocity complex128
}

type Simulation struct {
	Boids           []Boid
	ProtectedRange  float64 // Boids avoid each other if they are closer than this distance
	AvoidanceFactor float64 // How much to avoid other boids
	VisibleRange    float64 // Boids try to match the velocity of boids within this range
	MatchingFactor  float64 // How much to match the velocity of other boids
	CenteringFactor float64 // How much to move towards the center of mass of other boids
	LEFT_MARGIN     int     // Boids turn around if they reach the margin
	RIGHT_MARGIN    int     // Boids turn around if they reach the margin
	TOP_MARGIN      int     // Boids turn around if they reach the margin
	BOTTOM_MARGIN   int     // Boids turn around if they reach the margin
	TurnFactor      float64 // How much to turn around when reaching the margin
	MinSpeed        float64 // Minimum speed of the boids
	MaxSpeed        float64 // Maximum speed of the boids
	DT              float64 // Time step
}

func NewSimulation(numBoids int, Dt float64) *Simulation {
	var boids []Boid = make([]Boid, numBoids)
	for i := 0; i < numBoids; i++ {
		boids[i] = Boid{
			Position: complex((rand.Float64() * float64(WIDTH)), rand.Float64()*float64(HEIGHT)),
			Velocity: complex(3, 2),
		}
	}
	return &Simulation{
		Boids:           boids,
		ProtectedRange:  10,
		AvoidanceFactor: 0.05,
		VisibleRange:    20,
		MatchingFactor:  0.0005,
		CenteringFactor: 0.05,
		LEFT_MARGIN:     WIDTH - 50,
		RIGHT_MARGIN:    50,
		TOP_MARGIN:      50,
		BOTTOM_MARGIN:   HEIGHT - 50,
		TurnFactor:      0.2,
		MinSpeed:        2,
		MaxSpeed:        3,
		DT:              Dt,
	}

}

func (s *Simulation) Update() {
	// Rule 1 Separation: Boids avoid running into each other
	for i := 0; i < len(s.Boids); i++ {
		var close_dx, close_dy float64
		var B1 = s.Boids[i]
		for j := 0; j < len(s.Boids); j++ {
			if i == j {
				continue
			}
			var B2 = s.Boids[j]
			var displacementVector = B1.Position - B2.Position
			var distance float64 = cmplx.Abs(displacementVector)
			if distance < s.ProtectedRange {
				close_dx += real(displacementVector)
				close_dy += imag(displacementVector)
			}
		}
		B1.Velocity += complex(close_dx, close_dy) * complex(s.AvoidanceFactor, 0)
		s.Boids[i] = B1

	}
	// fmt.Println(s.Boids[0].Velocity)

	// Rule 2 Alignment: Boids try to match the velocity of their neighbors
	for i := 0; i < len(s.Boids); i++ {
		var vel_avg complex128 = complex(0, 0)
		var neighbouringBoids int
		var B1 = s.Boids[i]
		for j := 0; j < len(s.Boids); j++ {
			if i == j {
				continue
			}
			var B2 = s.Boids[j]
			var displacementVector = B1.Position - B2.Position
			var distance float64 = cmplx.Abs(displacementVector)
			if distance < s.VisibleRange {
				vel_avg += B2.Velocity
				neighbouringBoids++
			}
		}
		if neighbouringBoids > 0 {
			vel_avg /= complex(float64(neighbouringBoids), 0)
			B1.Velocity += (vel_avg - B1.Velocity) * complex(s.MatchingFactor, 0)
		}
		s.Boids[i] = B1

	}

	// Rule 3 Cohesion: Boids try to move towards the center of mass of their neighbors
	for i := 0; i < len(s.Boids); i++ {
		var pos_avg complex128 = complex(0, 0)
		var neighbouringBoids int
		var B1 = s.Boids[i]
		for j := 0; j < len(s.Boids); j++ {
			if i == j {
				continue
			}
			var B2 = s.Boids[j]
			var displacementVector = B1.Position - B2.Position
			var distance float64 = cmplx.Abs(displacementVector)
			if distance < s.VisibleRange {

				pos_avg += B2.Position
				neighbouringBoids++
			}
		}
		if neighbouringBoids > 0 {
			pos_avg /= complex(float64(neighbouringBoids), 0)
			B1.Velocity += (pos_avg - B1.Position) * complex(s.CenteringFactor, 0)
		}
		s.Boids[i] = B1

	}

	// Smooth Turns
	for i := 0; i < len(s.Boids); i++ {
		var B1 = s.Boids[i]
		if real(B1.Position) < float64(s.LEFT_MARGIN) {
			B1.Velocity += complex(s.TurnFactor, 0)
		}
		if real(B1.Position) > float64(s.RIGHT_MARGIN) {
			B1.Velocity -= complex(s.TurnFactor, 0)
		}
		if imag(B1.Position) < float64(s.TOP_MARGIN) {
			B1.Velocity += complex(0, s.TurnFactor)
		}
		if imag(B1.Position) > float64(s.BOTTOM_MARGIN) {
			B1.Velocity -= complex(0, s.TurnFactor)
		}
		s.Boids[i] = B1

	}

	// make sure speed is capped between s.MinSpeed and s.MaxSpeed
	for i := 0; i < len(s.Boids); i++ {
		var B1 = s.Boids[i]
		var speed = cmplx.Abs(B1.Velocity)
		if speed > s.MaxSpeed {
			B1.Velocity = B1.Velocity / complex(speed, 0)
			var newVelo = complex(real(B1.Velocity)*s.MaxSpeed, imag(B1.Velocity)*s.MinSpeed)
			B1.Velocity = newVelo

		}
		if speed < s.MinSpeed {
			B1.Velocity = B1.Velocity / complex(speed, 0) * complex(s.MinSpeed, 0)

		}
		s.Boids[i] = B1
	}

	// Update position
	for i := 0; i < len(s.Boids); i++ {
		var B1 = s.Boids[i]
		B1.Position += B1.Velocity * complex(s.DT, 0)
		s.Boids[i] = B1

	}
	// fmt.Println(s.Boids[0].Position)

}

func (s *Simulation) Draw(screen *ebiten.Image) {
	for _, boid := range s.Boids {
		ebitenutil.DrawCircle(screen, real(boid.Position), imag(boid.Position), 2, color.White)
	}

}
