package main

import (
	"fmt"
	"strings"
)

// Item is the common interface for leaves and composites.
type Item interface {
	Price() float64
	Display(indent int)
}

const (
	glassPrice = 40.0
	framePrice = 25.0
)

func line(indent int, format string, args ...any) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", indent), fmt.Sprintf(format, args...))
}

// Glass is a leaf with a fixed price.
type Glass struct{}

func (Glass) Price() float64 { return glassPrice }
func (g Glass) Display(indent int) {
	line(indent, "Glass (%.2f)", g.Price())
}

// Frame is a leaf with a fixed price.
type Frame struct{}

func (Frame) Price() float64 { return framePrice }
func (f Frame) Display(indent int) {
	line(indent, "Frame (%.2f)", f.Price())
}

// Window is the composite: it holds at least one child Item,
// which may itself be a Window.
type Window struct {
	items []Item
}

// NewWindow requires at least one child, enforced by the signature.
func NewWindow(first Item, rest ...Item) *Window {
	return &Window{items: append([]Item{first}, rest...)}
}

// Price is the sum of the children's prices.
func (w *Window) Price() float64 {
	var total float64
	for _, it := range w.items {
		total += it.Price()
	}
	return total
}

// Display prints the tree, indenting children one level deeper.
func (w *Window) Display(indent int) {
	line(indent, "Window (%.2f)", w.Price())
	for _, it := range w.items {
		it.Display(indent + 1)
	}
}
