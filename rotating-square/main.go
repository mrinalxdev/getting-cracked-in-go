// package main

// import (
// 	"github.com/gen2brain/raylib-go/raylib"
// )

// func main() {
// 	// Initialize the window and set the width, height, and title
// 	rl.InitWindow(800, 600, "3D Rotating Cube in Golang")

// 	// Set the camera position and where it is looking at
// 	camera := rl.Camera{
// 		Position:   rl.NewVector3(4.0, 4.0, 4.0),   // Camera position in space
// 		Target:     rl.NewVector3(0.0, 0.0, 0.0),   // The point the camera is looking at
// 		Up:         rl.NewVector3(0.0, 1.0, 0.0),   // Camera's up direction
// 		Fovy:       45.0,                           // Field of view in degrees
// 		Type:       rl.CameraPerspective,           // Perspective camera type
// 	}

// 	rl.SetTargetFPS(60) // Set the frame rate

// 	rotationAngle := float32(0.0) // Variable to track the cube's rotation

// 	// Main game loop
// 	for !rl.WindowShouldClose() {
// 		// Update the rotation angle of the cube
// 		rotationAngle += 0.5

// 		// Begin drawing the 3D scene
// 		rl.BeginDrawing()

// 		rl.ClearBackground(rl.RayWhite) // Clear the window with a white background

// 		// Start 3D mode
// 		rl.BeginMode3D(camera)

// 		// Draw a rotating cube
// 		rl.DrawCube(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 2.0, rl.Blue) // Cube body
// 		rl.DrawCubeWires(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 2.0, rl.Black) // Cube edges

// 		// Apply the rotation on the cube
// 		rl.DrawCubeV(rl.NewVector3(0.0, 0.0, 0.0), rl.NewVector3(2.0, 2.0, 2.0), rl.Blue)
// 		rl.DrawCubeWiresV(rl.NewVector3(0.0, 0.0, 0.0), rl.NewVector3(2.0, 2.0, 2.0), rl.Black)

// 		// Rotate around the Y-axis
// 		rl.DrawCube(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 2.0, rl.NewColor(255, 0, 0, 255))
// 		rl.DrawCubeWires(rl.NewVector3(
// 		// Rotate the cube by applying a rotation matrix
// 		rl.PushMatrix(),     
// 		rl.RotateCubeV(rotationAngle, rl.NewVector3(0, 1, 0)), // Rotate around Y-axis
// 		rl.DrawCube(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 2.0, rl.Red),// Draw the rotating cube
// 		rl.PopMatrix(),                          // Restore the original transformation matrix

// 		// End 3D mode
// 		rl.EndMode3D(),

// 		// Draw some 2D text
// 		rl.DrawText("3D Rotating Cube in Golang", 10, 10, 20, rl.Black),

// 		// End drawing
// 		rl.EndDrawing(),
// 	},

// 	// Close the window and clean up resources when the program ends
// 	rl.CloseWindow()
// }

package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize the window
	rl.InitWindow(800, 600, "3D Rotating Cube in Golang")

	// Set the camera
	camera := rl.Camera3D{
		Position: rl.NewVector3(4.0, 4.0, 4.0),
		Target:   rl.NewVector3(0.0, 0.0, 0.0),
		Up:       rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:      45.0,
		Type:     rl.CameraPerspective,
	}

	// Set the frame rate
	rl.SetTargetFPS(60)

	// Initialize rotation angle
	rotationAngle := float32(0.0)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update rotation angle
		rotationAngle += 0.5

		// Begin drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Start 3D mode
		rl.BeginMode3D(camera)

		// Push transformation matrix
		rl.PushMatrix()

		// Apply rotation
		rl.RotateY(rotationAngle)

		// Draw rotating cube
		rl.DrawCube(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 2.0, rl.Red)
		rl.DrawCubeWires(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 2.0, rl.Black)


		rl.PopMatrix()

		// End 3D mode
		rl.EndMode3D()

		// Draw 2D text
		rl.DrawText("3D Rotating Cube in Golang", 10, 10, 20, rl.Black)

		// End drawing
		rl.EndDrawing()
	}
    

	// Close window and clean up resources
	rl.CloseWindow()
}
