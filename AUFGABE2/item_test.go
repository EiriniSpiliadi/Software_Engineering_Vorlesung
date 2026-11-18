package main

import "testing"

func TestPrice(t *testing.T) {
	inner := NewWindow(Frame{}, Glass{}, Glass{}) // 25 + 40 + 40 = 105

	tests := []struct {
		name string
		item Item
		want float64
	}{
		{"glass", Glass{}, 40},
		{"frame", Frame{}, 25},
		{"single child", NewWindow(Glass{}), 40},
		{"flat window", NewWindow(Frame{}, Glass{}), 65},
		{"nested window", NewWindow(Frame{}, inner, Frame{}), 155},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Price(); got != tt.want {
				t.Errorf("Price() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}
