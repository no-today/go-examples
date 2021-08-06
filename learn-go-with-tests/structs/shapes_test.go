package main

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{Width: 10.0, Height: 10.0}
	got := rectangle.Perimeter()
	want := 40.0
	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	areaTests := []struct {
		name     string
		shape    Shape
		wantArea float64
	}{
		{"Rectangle", Rectangle{10.0, 15.0}, 150.0},
		{"Circle", Circle{10.0}, 314.1592653589793},
		{"Triangle", Triangle{12, 6}, 36},
	}

	for _, tt := range areaTests {
		// using tt.name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.wantArea {
				t.Errorf("%#v got %.2f want %.2f", tt.shape, got, tt.wantArea)
			}
		})
	}
}
