package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// Constants:
const errorDelta = 1e-8
const widthRes = 2000
const heightRes = 2000

var xTrack = int(500) //(50 * math.Sqrt(2)))
var yTrack = int(500) //(50 * math.Sqrt(2)))

var blue = color.NRGBA{
	R: uint8(0),
	G: uint8(0),
	B: uint8(255),
	A: 255,
}
var orange = color.NRGBA{
	R: uint8(250),
	G: uint8(70),
	B: uint8(10),
	A: 255,
}
var purple = color.NRGBA{
	R: uint8(175),
	G: uint8(90),
	B: uint8(210),
	A: 255,
}
var white = color.NRGBA{
	R: uint8(255),
	G: uint8(255),
	B: uint8(255),
	A: 255,
}
var black = color.NRGBA{
	R: uint8(0),
	G: uint8(0),
	B: uint8(0),
	A: 255,
}
var pink = color.NRGBA{
	R: uint8(140),
	G: uint8(20),
	B: uint8(20),
	A: 255,
} // End of constants

// str, screen -> write(png)
// Write stored in a screen struct to a file
func createImage(fileName string, scrn screen) {
	img := image.NewRGBA(image.Rect(0, 0, int(widthRes), int(heightRes)))
	for y := 0; y < int(widthRes); y++ {
		for x := 0; x < int(heightRes); x++ {
			// img.Set(x, y, scrn.at(x, y))
			img.Set(x, y, scrn.pixels[x][y])
		}
	}
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func initializeScene() (screen, light, []object) {
	// creating the perspective point, and plane to be colored
	var scrn = makeScreen(500, 500, 80.0)

	// var sun = light{makeVector(-40, -1, -10)}
	var sun = light{makeVector(-40, 0, 10)}
	// var sun = light{makeVector(500, -1000, 20)}
	// var sun = light{makeVector(0, -1, 0)}
	// var sun = light{makeVector(0, 200, 50)}
	// var sun = light{makeVector(0, 90, 0)}

	bx := makeBox(
		2, 90, 2,
		6, 110, 6,
		white,
	)
	bx2 := makeBox(
		-8, 100, -8,
		-2, 120, -5,
		blue,
	)
	sphr := makeSphere(7, 110, 3, 3, purple) // an object in the worldspace
	objects := make([]object, 0)
	objects = append(objects, bx)
	objects = append(objects, bx2)
	objects = append(objects, sphr)
	objects = append(objects, makeSphere(-3, 100, 3, 4, blue))
	objects = append(objects, makeSphere(3, 130, -2, 2, pink))

	// objects := emma()
	return scrn, sun, objects
}

func getDirection(scrn screen, x, y int) (vector, vector) {
	screenPoint := scrn.point(x, y) // worldspace coordinate for a given x,y pixel coordinate
	zoom := 10.0                    //(scrn.width / 2) * math.Tan(scrn.fov/2)
	dir := (screenPoint.sub(makeVector(0, -zoom, 0))).direction()
	return screenPoint, dir
}

func main() {
	scrn, sun, objects := initializeScene()

	for x := 0; x < widthRes; x++ {
		for y := 0; y < widthRes; y++ {
			screenPoint, dir := getDirection(scrn, x, y)

			minMag := math.Inf(1)
			for _, obj := range objects {
				// dir * obstructionMag + screenPoint = P
				obstuctionMag := obj.obstruct(dir, screenPoint, (x == xTrack && y == yTrack))
				// if obstuctionMag != math.Inf(1) {
				// 	fmt.Println(obstuctionMag)
				// }
				if obstuctionMag < minMag { //<= 60 { //
					// fmt.Println(dir.mul(obstuctionMag).add(screenPoint).mag())
					minMag = obstuctionMag
					scrn.pixels[x][y] = obj.directIllumination(sun, dir.mul(minMag).add(screenPoint), objects) // Inefficient to calc before knowing min
				}
			}

			if x == xTrack && y == yTrack {
				endPoint := dir.mul(minMag).add(screenPoint)
				fmt.Printf("(%d,%d)\nscreenPoint:%v,\ndir:%v,\nendPoint:%v,\nminMag:%f \n", xTrack, yTrack, screenPoint, dir, endPoint, minMag)

				scrn.pixels[x][y] = pink
			}
		}
	}

	createImage("image.png", scrn)
}
