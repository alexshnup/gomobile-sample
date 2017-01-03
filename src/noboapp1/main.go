// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux

// An app that draws a green triangle on a red background.
//
// Note: This demo is an early preview of Go 1.5. In order to build this
// program as an Android APK using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d golang.org/x/mobile/example/basic
//   $ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install golang.org/x/mobile/example/basic && basic
package main

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/alexshnup/material"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32
)

var (
	env = new(material.Environment)

	t112, t56, t45, t34, t24, t20, t16, t14, t12 *material.Button
)

func init() {
	env.SetPalette(material.Palette{
		Primary: material.BlueGrey500,
		Dark:    material.BlueGrey700,
		Light:   material.BlueGrey100,
		Accent:  material.DeepOrangeA200,
	})
}

func onStart(ctx gl.Context) {
	ctx.Enable(gl.BLEND)
	ctx.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	ctx.Enable(gl.CULL_FACE)
	ctx.CullFace(gl.BACK)

	env.Load(ctx)
	env.LoadGlyphs(ctx)

	t112 = env.NewButton(ctx)
	t112.SetTextColor(material.White)
	t112.SetText("AAAH`e_llo |jJ go 112px")
	t112.BehaviorFlags = material.DescriptorFlat

	t56 = env.NewButton(ctx)
	t56.SetTextColor(material.White)
	t56.SetText("Hello go 56px")
	t56.BehaviorFlags = material.DescriptorFlat

	t45 = env.NewButton(ctx)
	t45.SetTextColor(material.White)
	t45.SetText("Hello go 45px")
	t45.BehaviorFlags = material.DescriptorFlat

	t34 = env.NewButton(ctx)
	t34.SetTextColor(material.White)
	t34.SetText("Hello go 34px")
	t34.BehaviorFlags = material.DescriptorFlat

	t24 = env.NewButton(ctx)
	t24.SetTextColor(material.White)
	t24.SetText("Hello go 24px")
	t24.BehaviorFlags = material.DescriptorFlat

	t20 = env.NewButton(ctx)
	t20.SetTextColor(material.White)
	t20.SetText("Hello go 20px")
	t20.BehaviorFlags = material.DescriptorFlat

	t16 = env.NewButton(ctx)
	t16.SetTextColor(material.White)
	t16.SetText("Hello go 16px")
	t16.BehaviorFlags = material.DescriptorFlat

	t14 = env.NewButton(ctx)
	t14.SetTextColor(material.White)
	t14.SetText("Hello go 14px")
	t14.BehaviorFlags = material.DescriptorFlat

	t12 = env.NewButton(ctx)
	t12.SetTextColor(material.White)
	t12.SetText("Hello go 12px")
	t12.BehaviorFlags = material.DescriptorFlat
}

func onLayout(sz size.Event) {
	env.SetOrtho(sz)
	env.StartLayout()
	env.AddConstraints(
		t112.Width(1290), t112.Height(112), t112.Z(1), t112.StartIn(env.Box, env.Grid.Gutter), t112.TopIn(env.Box, env.Grid.Gutter),
		t56.Width(620), t56.Height(56), t56.Z(1), t56.StartIn(env.Box, env.Grid.Gutter), t56.Below(t112.Box, env.Grid.Gutter),
		t45.Width(500), t45.Height(45), t45.Z(1), t45.StartIn(env.Box, env.Grid.Gutter), t45.Below(t56.Box, env.Grid.Gutter),
		t34.Width(380), t34.Height(34), t34.Z(1), t34.StartIn(env.Box, env.Grid.Gutter), t34.Below(t45.Box, env.Grid.Gutter),
		t24.Width(270), t24.Height(24), t24.Z(1), t24.StartIn(env.Box, env.Grid.Gutter), t24.Below(t34.Box, env.Grid.Gutter),
		t20.Width(230), t20.Height(20), t20.Z(1), t20.StartIn(env.Box, env.Grid.Gutter), t20.Below(t24.Box, env.Grid.Gutter),
		t16.Width(180), t16.Height(16), t16.Z(1), t16.StartIn(env.Box, env.Grid.Gutter), t16.Below(t20.Box, env.Grid.Gutter),
		t14.Width(155), t14.Height(14), t14.Z(1), t14.StartIn(env.Box, env.Grid.Gutter), t14.Below(t16.Box, env.Grid.Gutter),
		t12.Width(135), t12.Height(12), t12.Z(1), t12.StartIn(env.Box, env.Grid.Gutter), t12.Below(t14.Box, env.Grid.Gutter),
	)
	log.Println("starting layout")
	t := time.Now()
	env.FinishLayout()
	log.Printf("finished layout in %s\n", time.Now().Sub(t))
}

var lastpaint time.Time
var fpsmaterial int

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				repaint(a)
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

// func onStart(glctx gl.Context) {
// 	var err error
// 	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
// 	if err != nil {
// 		log.Printf("error creating GL program: %v", err)
// 		return
// 	}
//
// 	buf = glctx.CreateBuffer()
// 	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
// 	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)
//
// 	position = glctx.GetAttribLocation(program, "position")
// 	color = glctx.GetUniformLocation(program, "color")
// 	offset = glctx.GetUniformLocation(program, "offset")
//
// 	images = glutil.NewImages(glctx)
// 	fps = debug.NewFPS(images)
// }

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	glctx.ClearColor(material.BlueGrey500.RGBA())
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	env.Draw(glctx)
	now := time.Now()
	fpsmaterial = int(time.Second / now.Sub(lastpaint))
	lastpaint = now

	glctx.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}
	// glctx.Uniform4f(color, 0, green, 0, 1)
	//
	// glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))
	//
	// glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	// glctx.EnableVertexAttribArray(position)
	// glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	// glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	// glctx.DisableVertexAttribArray(position)

	// fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.4, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
