package main

import (
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

// scrn -> *image.RGBA
// returns an *image.RGBA given a screen, such that the image.RGBA pixel colors and dimensions
// correspond to the colors and dimensions stored in the screen
func screenToImage(scrn screen) (img *image.RGBA) {
	img = image.NewRGBA(image.Rect(0, 0, int(widthRes), int(heightRes)))
	for y := 0; y < int(widthRes); y++ {
		for x := 0; x < int(heightRes); x++ {
			// img.Set(x, y, scrn.at(x, y))
			img.Set(x, y, scrn.pixels[x][y])
		}
	}
	return img
}

// str, screen ->
// Write data stored in a screen to a png file of given name.
// (creates file if it does not already exist)
func writeToPNG(fileName string, scrn screen) {
	img := screenToImage(scrn)

	f, err := os.Create(fileName)
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

// -> screen, light, []object
// Returns the variables that define a worldspace:
// a perspective screen, lighting, and objects
func initializeScene() (screen, light, []object) {
	// creating the perspective point, and plane to be colored
	var scrn = makeScreen(10, 10, 80.0)

	var sun = light{makeVector(-40, 0, 10)}

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

	return scrn, sun, objects
}

// screen, x, y -> vector, vector
// given a perspective screen and index for a pixel in the screen,
// returns the pixels corresponding point in the worldspace and direction
// of a vector through that point from the perspective origin
func getDirection(scrn screen, x, y int) (vector, vector) {
	screenPoint := scrn.point(x, y)                                 // worldspace coordinate for a given x,y pixel coordinate
	forwardComponent := (scrn.width / 2) / (math.Tan(scrn.fov / 2)) // forward component for the camera
	dir := (screenPoint.add(makeVector(0, forwardComponent, 0))).direction()
	return screenPoint, dir
}

// ->
// Initializes a worldspace, including a screen
// Determines the color to display for each pixel in the screen by:
// finding a ray through the pixel's corresponding point in worldspace, then
// finding the color of the object with which that ray collides first.
func main() {
	scrn, sun, objects := initializeScene()

	for x := 0; x < widthRes; x++ {
		for y := 0; y < widthRes; y++ {
			screenPoint, dir := getDirection(scrn, x, y)

			minMag := math.Inf(1)
			var closestObj object
			for _, obj := range objects {
				obstuctionMag := obj.obstruct(dir, screenPoint, false)
				if obstuctionMag < minMag {
					minMag = obstuctionMag
					closestObj = obj
				}
			}
			if minMag < math.Inf(1) {
				scrn.pixels[x][y] = (closestObj).directIllumination(sun, dir.mul(minMag).add(screenPoint), objects)
			}
		}
	}

	writeToPNG("image.png", scrn)
}
