package main

import rl "github.com/gen2brain/raylib-go/raylib"
import "boids/sim"

func main() {
	rl.InitWindow(sim.WindowWidth, sim.WindowHeight, "Boid Simulation")
	defer rl.CloseWindow()

	// Initialise all the boids before beginning render loop
	var flock []sim.Boid

	for range sim.TotalFlockSize {
		boid := sim.NewBoid()
		flock = append(flock, *boid)
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for _, boid := range flock {
			boid.Draw(7, rl.RayWhite)
		}

		rl.EndDrawing()
	}
}
