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

	LocalFlock      []*Boid
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

func (boid *Boid) DrawPerceptionField(color rl.Color) {
	// Total visible angle: 360° - blind spot
	visibleAngle := 360.0 - BlindSpotAngle

	// Start the arc centered around boid's facing angle
	boidAngleRad := math.Atan2(float64(boid.Velocity.Y), float64(boid.Velocity.X))
	boidAngleDeg := boidAngleRad * rl.Rad2deg

	startAngle := boidAngleDeg - visibleAngle/2

	rl.DrawCircleSectorLines(
		boid.CurPos,
		PerceptionRadius,
		float32(startAngle),
		float32(startAngle+visibleAngle),
		60,
		color,
	)
}

func (boid *Boid) Update(flock []*Boid) {
	boid.FindLocalFlock(flock)
	boid.align(boid.avgVelocity())

	amountSteeringVecs := len(boid.SteeringVectors)
	avgSteeringVector := rl.Vector2Zero()

	// Only steer boid if steering vectors are applied
	if amountSteeringVecs != 0 {
		for _, vec := range boid.SteeringVectors {
			avgSteeringVector = rl.Vector2Add(avgSteeringVector, vec)
		}

		boid.Acceleration = rl.Vector2Scale(
			avgSteeringVector,
			1.0/float32(amountSteeringVecs),
		)
	}
	// Remove all applied steering vectors
	boid.SteeringVectors = boid.SteeringVectors[:0]

	boid.Velocity = rl.Vector2Add(boid.Velocity, boid.Acceleration)

	// Smoothly de-accelerate towards velocity limit
	speed := rl.Vector2Length(boid.Velocity)
	if speed > VelocityLimit {
		target := rl.Vector2Scale(rl.Vector2Normalize(boid.Velocity), VelocityLimit)
		step := (speed - VelocityLimit) * DampingFactor

		boid.Velocity = rl.Vector2MoveTowards(boid.Velocity, target, step)
	}

	boid.CurPos = rl.Vector2Add(boid.CurPos, boid.Velocity)

	boid.wrap(10.0)
}

func (boid *Boid) wrap(padding float32) {
	if boid.CurPos.X > WindowWidth {
		boid.CurPos.X -= WindowWidth + padding
	} else if boid.CurPos.X < 0 {
		boid.CurPos.X += WindowWidth + padding
	}

	if boid.CurPos.Y > WindowHeight {
		boid.CurPos.Y -= WindowHeight + padding
	} else if boid.CurPos.Y < 0 {
		boid.CurPos.Y += WindowHeight + padding
	}
}

func (boid *Boid) FindLocalFlock(flock []*Boid) {
	// Remove all references before re-calculating local flock
	boid.LocalFlock = boid.LocalFlock[:0]

	for _, other_boid := range flock {
		if boid == other_boid {
			continue
		}

		if rl.Vector2Distance(boid.CurPos, other_boid.CurPos) > PerceptionRadius {
			continue
		}

		angle_to_other_boid := rl.Vector2Angle(boid.Velocity, other_boid.Velocity)
		is_visible := angle_to_other_boid < BlindSpotAngle || angle_to_other_boid > 360-BlindSpotAngle

		if is_visible {
			boid.LocalFlock = append(boid.LocalFlock, other_boid)
		}
	}
}

func (boid Boid) avgVelocity() rl.Vector2 {
	localFlockSize := len(boid.LocalFlock)
	avgVelocity := rl.Vector2Zero()

	if localFlockSize == 0 {
		return avgVelocity
	}

	for _, other_boid := range boid.LocalFlock {
		avgVelocity = rl.Vector2Add(avgVelocity, other_boid.Velocity)
	}

	avgVelocity = rl.Vector2Scale(avgVelocity, 1.0/float32(localFlockSize))
	return avgVelocity
}

func (boid Boid) avgPosition() rl.Vector2 {
	localFlockSize := len(boid.LocalFlock)
	avgPosition := rl.Vector2Zero()

	if localFlockSize == 0 {
		return avgPosition
	}

	for _, other_boid := range boid.LocalFlock {
		avgPosition = rl.Vector2Add(avgPosition, other_boid.CurPos)
	}

	avgPosition = rl.Vector2Scale(avgPosition, 1.0/float32(localFlockSize))
	return avgPosition
}

func (boid *Boid) align(avgVelocity rl.Vector2) {
	if rl.Vector2Length(avgVelocity) != 0 {
		steeringVec := rl.Vector2Subtract(avgVelocity, boid.Velocity)
		boid.SteeringVectors = append(boid.SteeringVectors, steeringVec)
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

	if Debug {
		boid.DrawPerceptionField(rl.Red)
	}
}
