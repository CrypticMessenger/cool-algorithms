package pkg

import (
	"image/color"
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PARTICLE_SIZE = 5.0
)

type FieldSimulation struct {
	WIDTH         int
	HEIGHT        int
	PARTICLES     []Particle
	VECTOR_FIELD  [][]complex128
	NUM_PARTICLES int
	COLS          int
	ROWS          int
	TRAIL_RECORD  int
	PARTICLE_SIZE float64
	DT            float64
	showField     bool
	particleColor color.RGBA
	trailColor    color.RGBA
	// d/dt [x, y] = [f(x, y), g(x, y)]
	F_XY func(complex128) float64
	G_XY func(complex128) float64
}
type Particle struct {
	isOutside bool
	Position  complex128
	Trail     []complex128
}

func NewFieldSimulation(width, height, NUM_PARTICLES, COLS, ROWS, TRAIL_RECORD int, F_XY func(complex128) float64, G_XY func(complex128) float64, particleSize float64, speedup float64, showField bool, particleColorHead color.RGBA, trailColorTail color.RGBA) *FieldSimulation {
	var spacing float64 = float64(width) / float64(COLS)
	var particles []Particle = make([]Particle, NUM_PARTICLES)
	for i := 0; i < NUM_PARTICLES; i++ {
		particles[i].Position = complex(math.Floor(rand.Float64()*spacing*float64(ROWS)), math.Floor(rand.Float64()*spacing*float64(COLS)))
		particles[i].Trail = make([]complex128, TRAIL_RECORD)
		for j := 0; j < TRAIL_RECORD; j++ {
			particles[i].Trail[j] = particles[i].Position
		}
	}
	var vector_field [][]complex128 = make([][]complex128, ROWS)
	for i := 0; i < ROWS; i++ {
		vector_field[i] = make([]complex128, COLS)
	}
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			var x float64 = float64(i) * spacing
			var y float64 = float64(j) * spacing
			vector_field[i][j] = complex(F_XY(complex(x, y)), G_XY(complex(x, y)))
			// normalize the vector
			var length float64 = cmplx.Abs(vector_field[i][j])
			vector_field[i][j] = vector_field[i][j] / complex(length, 0)
		}
	}

	return &FieldSimulation{
		WIDTH:         width,
		HEIGHT:        height,
		COLS:          COLS,
		ROWS:          ROWS,
		TRAIL_RECORD:  TRAIL_RECORD,
		NUM_PARTICLES: NUM_PARTICLES,
		PARTICLES:     particles,
		VECTOR_FIELD:  vector_field,
		PARTICLE_SIZE: particleSize,
		DT:            speedup,
		showField:     showField,
		particleColor: particleColorHead,
		trailColor:    trailColorTail,
		F_XY:          F_XY,
		G_XY:          G_XY,
	}
}

func DrawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	// Draw a line from (x1, y1) to (x2, y2)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)

}

func (fs *FieldSimulation) DrawVectorField(vector_field [][]complex128, screen *ebiten.Image) {
	var spacing = float64(fs.WIDTH) / float64(fs.COLS)
	for i := 0; i < fs.ROWS; i++ {
		for j := 0; j < fs.COLS; j++ {
			var x float64 = float64(i) * spacing
			var y float64 = float64(j) * spacing
			var v complex128 = vector_field[i][j]
			// Draw the vector
			DrawLine(screen, x, y, x+(spacing/2.0)*real(v), y+(spacing/2.0)*imag(v))

		}
	}

}

func (fs *FieldSimulation) Draw(screen *ebiten.Image) {
	if fs.showField {
		// Draw the vector field
		fs.DrawVectorField(fs.VECTOR_FIELD, screen)
	}
	var c1 = fs.particleColor
	var c2 = fs.trailColor
	var dr = (c1.R - c2.R) / uint8(fs.TRAIL_RECORD-1)
	var dg = (c1.G - c2.G) / uint8(fs.TRAIL_RECORD-1)
	var db = (c1.B - c2.B) / uint8(fs.TRAIL_RECORD-1)

	// Draw the particle
	for _, particle := range fs.PARTICLES {
		if particle.isOutside {
			continue
		}
		ebitenutil.DrawCircle(screen, real(particle.Position), imag(particle.Position), fs.PARTICLE_SIZE, c1)
	}

	for _, particle := range fs.PARTICLES {
		if particle.isOutside {
			continue
		}
		for i := 0; i < len(particle.Trail)-1; i++ {
			ebitenutil.DrawLine(screen, real(particle.Trail[i]), imag(particle.Trail[i]), real(particle.Trail[i+1]), imag(particle.Trail[i+1]), color.RGBA{c2.R + uint8(i)*dr, c2.G + uint8(i)*dg, c2.B + uint8(i)*db, 255})
		}
	}

}

func (fs *FieldSimulation) Update() {
	// Update the particles
	var spacing float64 = float64(fs.WIDTH) / float64(fs.COLS)
	for i := 0; i < fs.NUM_PARTICLES; i++ {
		var x = int(real(fs.PARTICLES[i].Position) / spacing)
		var y = int(imag(fs.PARTICLES[i].Position) / spacing)
		if x < 0 || x >= fs.ROWS || y < 0 || y >= fs.COLS {
			fs.PARTICLES[i].isOutside = true
			continue
		}
		var v = fs.VECTOR_FIELD[x][y]
		fs.PARTICLES[i].Position += v * complex(fs.DT, 0)
		fs.PARTICLES[i].Trail = append(fs.PARTICLES[i].Trail, fs.PARTICLES[i].Position)
		if len(fs.PARTICLES[i].Trail) > fs.TRAIL_RECORD {
			fs.PARTICLES[i].Trail = fs.PARTICLES[i].Trail[1:]
		}
	}
}

func (fs *FieldSimulation) GetTotalHeight() int {
	return fs.HEIGHT
}

func (fs *FieldSimulation) GetTotalWidth() int {
	return fs.WIDTH
}
