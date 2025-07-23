package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main()  {
	rl.InitWindow(800, 450, "Boid Simulation")
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawText("This will have some boids soon", 190, 200, 20, rl.RayWhite)

		rl.EndDrawing()
	}
}
