package sim

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Boid struct {
	CurPos       rl.Vector2
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

	randAngle := randRange(0, 360) * rl.Deg2rad
	dir_vec := rl.NewVector2(
		float32(math.Cos(float64(randAngle))),
		float32(math.Sin(float64(randAngle))),
	)
	speed := randRange(0.1, VelocityLimit)

	newBoid := Boid{
		CurPos:   randPos,
		Velocity: rl.Vector2Scale(dir_vec, speed),
	}

	return &newBoid
}

func (boid *Boid) Update() {
	boid.CurPos = rl.Vector2Add(boid.CurPos, boid.Velocity)
	boid.Velocity = rl.Vector2Add(boid.Velocity, boid.Acceleration)

	// Smoothly de-accelerate towards velocity limit
	speed := rl.Vector2Length(boid.Velocity)
	if speed > VelocityLimit {
		target := rl.Vector2Scale(rl.Vector2Normalize(boid.Velocity), VelocityLimit)
		step := (speed - VelocityLimit) * DampingFactor

		boid.Velocity = rl.Vector2MoveTowards(boid.Velocity, target, step)
	}

	boid.wrap()
}

func (boid *Boid) wrap()  {
	if boid.CurPos.X > WindowWidth {
		boid.CurPos.X -= WindowWidth
	} else if boid.CurPos.X < 0 {
		boid.CurPos.X += WindowWidth
	}

	if boid.CurPos.Y > WindowHeight {
		boid.CurPos.Y -= WindowHeight
	} else if boid.CurPos.Y < 0 {
		boid.CurPos.Y += WindowHeight
	}
}

func (boid Boid) Draw(offset float32, color rl.Color) {
	// Triangle is defined with tip pointing up (-Y). To align with velocity
	// direction, compute angle from velocity vector and rotate triangle by
	// (angle + pi/2) to correct for orientation mismatch.
	angleRad := float32(math.Atan2(float64(boid.Velocity.Y), float64(boid.Velocity.X))) + rl.Pi/2
	halfSize := offset / 2

	relativeVertices := []rl.Vector2{
		rl.NewVector2(0, -offset),
		rl.NewVector2(-halfSize, halfSize),
		rl.NewVector2(halfSize, halfSize),
	}

	for i := range relativeVertices {
		relativeVertices[i] = rl.Vector2Add(
			boid.CurPos,
			rl.Vector2Rotate(relativeVertices[i], angleRad),
		)
	}

	rl.DrawTriangle(
		relativeVertices[0],
		relativeVertices[1],
		relativeVertices[2],
		color,
	)
}
