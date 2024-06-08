package pkg

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ReccurencePlots struct {
	DELTA        float64
	WIDTH        int
	HEIGHT       int
	EXTRA_HEIGHT int
	f_n          func(int) float64
}

func NewReccurencePlots() *ReccurencePlots {
	return &ReccurencePlots{
		DELTA:        6,
		WIDTH:        620,
		HEIGHT:       620,
		EXTRA_HEIGHT: 180, // for wave plot
		f_n:          func(x int) float64 { return 50 * math.Sin(2.0*math.Pi*(6+5*math.Sin(2*2.0*math.Pi*float64(x)/620))) },
	}
}

func (rp *ReccurencePlots) Draw(screen *ebiten.Image) {
	// Recurrence plot: A data series f(t) is plotted in a coordinate system on both the x and y axes. All values ​​are compared with each other. If their difference falls below a certain value (delta), a point is placed at this point. The pattern created in this way is the recurrence plot.
	for i := 0; i < rp.WIDTH; i++ {
		for j := rp.HEIGHT - 1; j >= 0; j-- {
			if math.Abs(rp.f_n(i)-rp.f_n(rp.HEIGHT-1-j)) < rp.DELTA {
				c_r := uint8(255 * rand.Float64())
				c_g := uint8(255 * rand.Float64())
				c_b := uint8(255 * rand.Float64())
				screen.Set(i, j, color.RGBA{c_r, c_g, c_b, 255})
				// screen.Set(i, j, color.RGBA{0, 255, 0, 255})
			}
		}

	}
	// Wave plot
	for i := 0; i < rp.WIDTH-1; i++ {
		ebitenutil.DrawLine(screen, float64(i), float64(2*rp.HEIGHT+rp.EXTRA_HEIGHT)/2+rp.f_n(i), float64(i+1), float64(2*rp.HEIGHT+rp.EXTRA_HEIGHT)/2+rp.f_n(i+1), color.RGBA{0, 255, 0, 255})
		// screen.Set(i, (2*rp.HEIGHT+rp.EXTRA_HEIGHT)/2+int(rp.f_n(i)), color.RGBA{255, 0, 0, 255})
	}
}
func (rp *ReccurencePlots) GetTotalHeight() int {
	return rp.HEIGHT + rp.EXTRA_HEIGHT
}
func (rp *ReccurencePlots) GetTotalWidth() int {
	return rp.WIDTH
}

func (rp *ReccurencePlots) Update() {
}
