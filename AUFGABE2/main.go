package main

import "fmt"

func main() {
	inner := NewWindow(Frame{}, Glass{}, Glass{})
	w := NewWindow(Frame{}, Glass{}, inner, Frame{})

	w.Display(0)
	fmt.Printf("\nTotal price: %.2f\n", w.Price())
}
