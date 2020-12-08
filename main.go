package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/gl/v4.6-core/gl"

	"log"
)

func main() {
	log.Print("main()")

	// init glfw
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	// open a window
	w, err := glfw.CreateWindow(500, 500, "Hello, world!", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	w.MakeContextCurrent()

	// init gl
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	// render loop
	for !w.ShouldClose() {
		w.SwapBuffers()
		glfw.PollEvents()
	}
}
