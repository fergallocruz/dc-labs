// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Point struct{ x, y float64 }

func (p Point) X(q Point) float64 {
	return q.x
}
func (p Point) Y(q Point) float64 {
	return q.y
}

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.
type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	fmt.Printf("- Figure's Perimeter \n - ")
	valuesText := []string{}
	for i := range path {
		if i > 0 {
			num := path[i-1].Distance(path[i])
			sum += num
			text := fmt.Sprintf("%f", num)
			valuesText = append(valuesText, text)
		}
	}

	// Join our string slice.
	result := strings.Join(valuesText, " + ")
	fmt.Print(result)
	return sum
}

//!-path

func main() {
	var a int
	a, _ = strconv.Atoi(os.Args[1])
	fmt.Printf("- Generating a [%d] sides figure \n", a)
	var p Path
	min := -100.0
	max := 100.0
	for i := 0; i < a; i++ {
		p1 := Point{(rand.Float64() * max) + min, (rand.Float64() * max) + min}
		p = append(p, p1)

	}
	fmt.Println("- Figure's vertices")
	for _, n := range p {
		fmt.Printf("	- ( %f, %f)\n", n.x, n.y)
	}
	distance := p.Distance()
	fmt.Printf(" = %f ", distance)

}
