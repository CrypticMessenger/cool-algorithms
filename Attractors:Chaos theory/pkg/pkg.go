package pkg

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Attractor interface {
	Draw(*ebiten.Image)
	GetTotalHeight() int
	GetTotalWidth() int
	Update()
}
