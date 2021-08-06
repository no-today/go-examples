package main

import "math"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

// Perimeter perimeter = 2(height + width)
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// area = width * height
func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Triangle struct {
	Base   float64
	Height float64
}

// area = 2/1(ah)
func (t Triangle) Area() float64 {
	return (t.Base * t.Height) * 0.5
}
