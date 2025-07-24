package main

import (
	"boids/sim"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(sim.WindowWidth, sim.WindowHeight, "Boid Simulation")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	// Initialise all boids before beginning render loop
	var flock []*sim.Boid

	for range sim.TotalFlockSize {
		boid := sim.NewBoid()
		flock = append(flock, boid)
	}

	for !rl.WindowShouldClose() {
		for _, boid := range flock {
			boid.Update(flock)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for _, boid := range flock {
			boid.Draw(sim.BoidSize, rl.RayWhite)
		}

		rl.EndDrawing()
	}
}
