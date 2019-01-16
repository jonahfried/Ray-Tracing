package main

import (
	"image/color"
	"math"
)

// Stores an x, y, and z float64 values representing coordinates in a cartesian coordinate system
type vector struct {
	x, y, z float64
}

// float64, float64, float64 -> vector
// Constructs a vector with given x, y, and z values
func makeVector(x, y, z float64) vector {
	return vector{x, y, z}
}

// vector, vector -> float64
// Performs a dot product of two vectors
func (v vector) dot(vec vector) float64 {
	return v.x*vec.x + v.y*vec.y + v.z*vec.z
}

// vector, vector -> vector
// returns the cross product of two given vectors
// (v1.x * v2.z) - (v1.z * v2.x)
func cross(v1, v2 vector) vector {
	x := (v1.y * v2.z) - (v1.z * v2.y)
	y := (v1.z * v2.x) - (v1.x * v2.z)
	z := (v1.x * v2.y) - (v1.y * v2.x)
	return makeVector(x, y, z)
}

// vector -> float64
// returns the magnitude of a given vector.
// √(x^2 + y^2 + z^2)
func (v vector) mag() float64 {
	return math.Sqrt(v.dot(v))
}

// vector -> vector
// returns the unit vector in the direction of given vector
// (returns zero vector if given vector has magnitude == 0)
func (v vector) direction() vector {
	mag := v.mag()
	if mag == 0 {
		return makeVector(0, 0, 0)
	}
	return v.mul(1 / mag)
}

// vector, float64 -> vector
// returns the result of scalar multiplication on a given vector
func (v vector) mul(s float64) vector {
	v.x *= s
	v.y *= s
	v.z *= s
	return v
}

// vector, vector -> vector
// returns the result of one vector added to another vector
func (v vector) add(vec vector) vector {
	v.x += vec.x
	v.y += vec.y
	v.z += vec.z
	return v
}

// vector, vector -> vector
// returns the result of one vector (vec) subtracted from to the first vector (v)
func (v vector) sub(vec vector) vector {
	v.x -= vec.x
	v.y -= vec.y
	v.z -= vec.z
	return v
}

// color.NRGBA, float64 -> color.NRGBA
// returns the result of multiplying the R, G, and B elements of an NRGBA struct
// (values capped at 255 before conversion back to uint8 unsigned integers)
func multiplyNRGBA(col color.NRGBA, m float64) color.NRGBA {
	if m < 0 {
		return col
	}
	r := math.Min(255, float64(col.R)*m)
	g := math.Min(255, float64(col.G)*m)
	b := math.Min(255, float64(col.B)*m)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}

// screen is a struct type representing a "camera lens"
// storing:
// 	fov ("field of view")
// 	aspectRatio
// 	height, width of camera in worldspace
// as well as a 2D array of color.NRGBA values corresponding to colors to be displayed on screen
type screen struct {
	fov           float64
	aspectRatio   float64
	height, width float64

	// perspectivePoint vector

	pixels [][]color.NRGBA
}

// float64, float64, float64 -> screen
// constructor for a screen. Takes in the desired height, width, and fov (in degrees), calculating and
// initializing the other values stored in a screen.
func makeScreen(height, width, fov float64) (scrn screen) {
	scrn.height = height
	scrn.width = width
	scrn.fov = (math.Pi * fov / 180) // convert from degrees to radians
	scrn.aspectRatio = (width / height)

	scrn.pixels = initializePixels()
	return scrn
}

// -> [][]color.NRGBA
// initializePixels allocates memory for a [][]color.NRGBA, and then allocates memory for each row's
// []color.NRGBA, and sets the default value for each pixel.
// Dimensions of the arrays are defined by the constants widthRes and heighRes.
func initializePixels() [][]color.NRGBA {
	pixels := make([][]color.NRGBA, int(widthRes), int(widthRes))
	for i := 0; i < int(widthRes); i++ {
		pixels[i] = make([]color.NRGBA, int(heightRes), int(heightRes))
		for j := 0; j < int(heightRes); j++ {
			pixels[i][j] = black
		}
	}
	return pixels
}

// screen, x, y -> vector
// returns the worldspace point relative to the given pixelspace coordinate
func (scrn screen) point(x, y int) vector {
	i := float64(x) / widthRes
	j := float64(y) / heightRes

	xdir := (2*i - 1) * math.Tan(scrn.fov/2) * scrn.aspectRatio
	zdir := (1 - (2 * j)) * math.Tan(scrn.fov/2)

	return makeVector(xdir, 0, zdir)
}
