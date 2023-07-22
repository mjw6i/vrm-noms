package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var blocks uint

func init() {
	runtime.LockOSThread()
	flag.UintVar(&blocks, "blocks", 10, "number of 64MB blocks of vram to reserve")
}

func main() {
	flag.Parse()
	flag.PrintDefaults()
	fmt.Printf("Blocking %d MB of VRAM\n", blocks*64)
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(640, 480, "Nope", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	for i := 0; i < int(blocks); i++ {
		id := newTexture64mb()
		defer gl.DeleteTextures(1, &id)
	}

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newTexture64mb() uint32 {
	return newTexture(1024 * 8)
}

func newTexture(imageSide int) uint32 {
	fakeImage := make([]uint8, imageSide*imageSide)

	var id uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.ALPHA,
		int32(imageSide),
		int32(imageSide),
		0,
		gl.ALPHA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(fakeImage))

	return id
}
