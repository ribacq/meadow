package main

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/renderer/shaders"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"

	"log"
	"math/rand"
	"time"
)

var plainFragmentShaderSource string = `
out vec4 FragColor;

void main()
{
	FragColor = vertex_color;
}
`

func init() {
	shaders.AddShader("plain_fragment", plainFragmentShaderSource)
	shaders.AddProgram("plain", "standard_vertex", "plain_fragment")
}

func main() {
	log.Print("birb: g3n test")
	log.Print(shaders.Shaders())

	// init random
	rand.Seed(time.Now().UnixNano())

	// app and scene
	a := app.App()
	scene := core.NewNode()

	// scene managed by gui manager
	gui.Manager().Set(scene)

	// perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// orbit control for camera
	camera.NewOrbitControl(cam)

	// handle window resize
	onResize := func(evname string, ev interface{}) {
		// get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// update camera’s aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// create mesh
	geom := geometry.NewSphere(1.0, 16, 16)
	mat := material.NewStandard(math32.NewColor("DarkGreen"))
	//mat.SetShader("plain")
	//mat.SetOpacity(0.8)
	//mat.SetSide(material.SideDouble)
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 1.0))
	pointLight := light.NewPoint(&math32.Color{1, 1, 0.8}, 1.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// create and add axis
	scene.Add(helper.NewAxes(0.5))

	// set grey background
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// run the app!
	a.Run(func(r *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		r.Render(scene, cam)
	})
}
