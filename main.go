package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/gl/v4.6-core/gl"

	"unsafe"
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

	// positions
	positions := [6]float32{
		-0.5, -0.5,
		 0.0,  0.5,
		 0.5, -0.5,
	}

	// create a vertex buffer
	var buffer uint32
	var zero int
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*int(unsafe.Sizeof(float32(0.0))), unsafe.Pointer(&positions), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*int32(unsafe.Sizeof(float32(0.0))), unsafe.Pointer(&zero))
	gl.EnableVertexAttribArray(0)

	// render loop
	for !w.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		w.SwapBuffers()
		glfw.PollEvents()
	}
}
