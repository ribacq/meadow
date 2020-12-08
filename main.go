package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"fmt"
	"io/ioutil"
	"log"
	"unsafe"
)

func compileShader(xtype uint32, src string) uint32 {
	id := gl.CreateShader(xtype)
	srcPtr, free := gl.Strs(src)
	defer free()
	gl.ShaderSource(id, 1, srcPtr, nil)
	gl.CompileShader(id)

	var result int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)
		messageStr := string(make([]rune, length))
		message := gl.Str(messageStr)
		gl.GetShaderInfoLog(id, length, &length, message)
		log.Fatalln("[Fatal] could not compile shader", gl.GoStr(message))
	}

	return id
}

func createVFShaderFromFiles(vertexShaderFilename, fragmentShaderFilename string) uint32 {
	vertexShaderBytes, err := ioutil.ReadFile(vertexShaderFilename)
	if err != nil {
		log.Fatal(err)
	}
	vertexShaderSrc := fmt.Sprintf("%s", vertexShaderBytes)

	fragmentShaderBytes, err := ioutil.ReadFile(fragmentShaderFilename)
	if err != nil {
		log.Fatal(err)
	}
	fragmentShaderSrc := fmt.Sprintf("%s", fragmentShaderBytes)

	program := gl.CreateProgram()
	vs := compileShader(gl.VERTEX_SHADER, vertexShaderSrc)
	fs := compileShader(gl.FRAGMENT_SHADER, fragmentShaderSrc)

	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return program
}

func main() {
	// init glfw
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	// open a window
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
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
		0.0, 0.5,
		0.5, -0.5,
	}

	// create a vertex buffer
	var buffer uint32
	var zero int
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*int(unsafe.Sizeof(float32(0.0))), unsafe.Pointer(&positions), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*int32(unsafe.Sizeof(float32(0.0))), unsafe.Pointer(&zero))

	// load shaders
	shaderProgram := createVFShaderFromFiles("shaders/vert00.glsl", "shaders/frag00.glsl")
	gl.UseProgram(shaderProgram)
	defer gl.DeleteProgram(shaderProgram)

	// render loop
	for !w.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		w.SwapBuffers()
		glfw.PollEvents()
	}
}
