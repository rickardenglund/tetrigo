package camera

import (
	"Tetrigo/glhelpers/keyhandler"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type CamState struct {
	pos, front, up mgl32.Vec3
	yaw, pitch     float64
	win            *glfw.Window
	keyHandler     keyhandler.Handler
	first          bool
	lastX, lastY   float64
}

const mouseSensitivity = 0.02

func New(win *glfw.Window) CamState {
	keyHandler, keyCallback := keyhandler.New()
	win.SetKeyCallback(keyCallback)
	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	c := CamState{
		up:         mgl32.Vec3{0, 1, 0},
		win:        win,
		keyHandler: keyHandler,
		first:      true,
	}

	return c
}

func (c *CamState) Update() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	c.inputHandler()
	x, y := c.win.GetCursorPos()
	if c.first {
		c.lastX = x
		c.lastY = y
		c.first = false
	}
	xd := (x - c.lastX) * mouseSensitivity
	yd := (y - c.lastY) * -1 * mouseSensitivity
	c.lastX = x
	c.lastY = y

	c.pitch += yd
	if c.pitch > 89 {
		c.pitch = 89
	}
	if c.pitch < -89 {
		c.pitch = -89
	}
	c.yaw += xd
	c.front = getCameraDir(c.yaw, c.pitch)

	return c.pos, c.pos.Add(c.front), c.up
}

func (c *CamState) inputHandler() {
	front := mgl32.Vec3{c.front.X(), 0, c.front.Z()}
	h := &c.keyHandler

	const camSpeed = 0.05

	if h.IsPressed(glfw.KeyW) {
		c.pos = front.Mul(camSpeed).Add(c.pos)
	}
	if h.IsPressed(glfw.KeyS) {
		c.pos = front.Mul(-camSpeed).Add(c.pos)
	}
	if h.IsPressed(glfw.KeyA) {
		c.pos = front.Cross(c.up).Normalize().Mul(-camSpeed).Add(c.pos)
	}
	if h.IsPressed(glfw.KeyD) {
		c.pos = front.Cross(c.up).Normalize().Mul(camSpeed).Add(c.pos)
	}
}

func getCameraDir(yaw, pitch float64) mgl32.Vec3 {
	x := math.Cos(ToRadians(yaw)) * math.Cos(ToRadians(pitch))
	y := math.Sin(ToRadians(pitch))
	z := math.Sin(ToRadians(yaw)) * math.Cos(ToRadians(pitch))

	return mgl32.Vec3{float32(x), float32(y), float32(z)}.Normalize()
}

func ToRadians(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}
