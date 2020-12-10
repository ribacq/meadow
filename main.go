package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unsafe"
)

func compileShader(xtype uint32, src string) uint32 {
	id := gl.CreateShader(xtype)
	srcPtr, free := gl.Strs(src + "\000")
	defer free()
	gl.ShaderSource(id, 1, srcPtr, nil)
	gl.CompileShader(id)

	var result int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)
		gl.DeleteShader(id)
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

	vs := compileShader(gl.VERTEX_SHADER, vertexShaderSrc)
	log.Println("vs ok")
	fs := compileShader(gl.FRAGMENT_SHADER, fragmentShaderSrc)
	log.Println("fs ok")

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(infoLog))

		log.Fatalln("could not link program", infoLog)
	}

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
	glfw.WindowHint(glfw.Resizable, glfw.False)
	w, err := glfw.CreateWindow(500, 500, "Hello, world!", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	w.MakeContextCurrent()

	// init gl
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	log.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	// positions
	positions := []float32{
		-0.5, -0.5,
		0.5, -0.5,
		0.5, 0.5,
		-0.5, 0.5,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	// create a vertex buffer
	var vbuffer uint32
	gl.GenBuffers(1, &vbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*int(unsafe.Sizeof(float32(0.0))), gl.Ptr(positions), gl.STATIC_DRAW)

	// create an index buffer
	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	// load shaders
	shaderProgram := createVFShaderFromFiles("res/shaders/vert00.glsl", "res/shaders/frag00.glsl")
	defer gl.DeleteProgram(shaderProgram)
	gl.UseProgram(shaderProgram)

	// render loop
	for !w.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		w.SwapBuffers()
		glfw.PollEvents()
	}
}
