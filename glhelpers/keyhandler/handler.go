package keyhandler

import "github.com/go-gl/glfw/v3.3/glfw"

type Handler struct {
	pressed map[glfw.Key]bool
}

func New() (Handler, glfw.KeyCallback) {
	h := Handler{pressed: make(map[glfw.Key]bool)}
	return h, h.keyCallback
}

func (h *Handler) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		h.pressed[key] = true
	}

	if action == glfw.Release {
		h.pressed[key] = false
	}
}

func (h *Handler) IsPressed(key glfw.Key) bool {
	return h.pressed[key]
}
