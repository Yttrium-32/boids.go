package sim

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Boid struct {
	CurPos rl.Vector2
	// Angle stores the boid's current Angle in degrees
	Angle  float32

	Velocity     rl.Vector2
	Acceleration rl.Vector2

	LocalFlock      []Boid
	SteeringVectors []rl.Vector2
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
		CurPos: randPos,
		Angle: randRange(0, 360),
	}

	return &newBoid
}

func (boid Boid) Draw(offset float32, color rl.Color) {
	angleRad := boid.Angle * rl.Deg2rad
	halfSize := offset / 2

	localVertices := []rl.Vector2 {
		rl.NewVector2(0, -offset),
		rl.NewVector2(-halfSize, halfSize),
		rl.NewVector2(halfSize, halfSize),
	}

	for i := range localVertices {
		localVertices[i] = rl.Vector2Add(
			boid.CurPos,
			rl.Vector2Rotate(localVertices[i], angleRad),
		)
	}

	rl.DrawTriangle(localVertices[0], localVertices[1], localVertices[2], color)
}
