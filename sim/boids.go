package sim

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Boid struct {
	curPos       rl.Vector2
	velocity     rl.Vector2
	acceleration rl.Vector2

	localFlock      []Boid
	steeringVectors []rl.Vector2
}

func randRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func NewBoid() *Boid {
	randPos := rl.NewVector2(
		randRange(0, WindowWidth),
		randRange(0, WindowHeight),
	)

	newBoid := Boid{
		curPos: randPos,
	}

	return &newBoid
}

func (boid Boid) Draw(offset float32, color rl.Color) {
	halfSize := offset / 2

	vertex1 := rl.NewVector2(boid.curPos.X, boid.curPos.Y - offset)
	vertex2 := rl.NewVector2(boid.curPos.X - halfSize, boid.curPos.Y + halfSize)
	vertex3 := rl.NewVector2(boid.curPos.X + halfSize, boid.curPos.Y + halfSize)

	rl.DrawTriangle(vertex1, vertex2, vertex3, color)
}
