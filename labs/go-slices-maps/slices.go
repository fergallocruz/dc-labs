package main

import (
	"code.google.com/p/go-tour/pic"
)

func Pic(dx, dy int) [][]uint8 {
	// el := se usa si la inicializas
	image := make([][]uint8, dx) //este es el slice
	for x := 0; x < dx; x++ {
		fila := make([]uint8, dy, dy)
		for y := 0; y < dy; y++ {
			// en cada fila se llenan las columnas con valores uint8
			fila[y] = uint8((x + y) / 2)
			fila[len(fila)-y-1] = fila[y] ^ uint8(x*y)
		}
		image[x] = fila
	}
	return image
}

func main() {
	pic.Show(Pic)
}
