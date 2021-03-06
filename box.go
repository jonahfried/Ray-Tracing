package main

import (
	"fmt"
	"image/color"
	"math"
)

// box is a struct type
// a box is an axis-aligned rectangular prism
// implements the object interface
type box struct {
	min, max vector
	center   vector
	col      color.NRGBA
}

// float64, float64, float64, float64, float64, float64, color.NRGBA -> box
// box constructor. returns a box given:
// x, y, z components of the points closest and furtherst from the origin.
func makeBox(min1, min2, min3, max1, max2, max3 float64, col color.NRGBA) (b box) {
	b.min = makeVector(min1, min2, min3)
	b.max = makeVector(max1, max2, max3)
	diff := b.max.sub(b.min)
	diff = diff.mul(.5)
	b.center = b.min.add(diff)
	b.col = col
	return b
}

// vector, float64, float64, float64, color.NRGBA -> box
// alternate box constructor. returns a box given:
// a vector representing the center of the box, a width, height and depth
func makeBox2(center vector, width, height, depth float64, col color.NRGBA) (b box) {
	dimensions := makeVector(width, depth, height).mul(.5)

	b.min = center.sub(dimensions)
	b.max = center.add(dimensions)
	b.center = center
	b.col = col
	return b
}

// box -> color.NRGBA
// returns a given box's color
func (b box) color() color.NRGBA {
	return b.col
}

// box, vector, vector, (verbosity) -> float64
// returns the magnitude of the ray at point of intersection with the prism.
// returns infinity if no intersection occurs.
// Determines the magnitudes at which the ray intersects planes that comprise the prism.
// looks for the smallest of these that satisfies the requirements of being on the prism.
func (b box) obstruct(direction, origin vector, verbosity bool) float64 {
	if verbosity {
		fmt.Println(b.min.x-origin.x, b.min.x, origin.x)
	}

	t0x := (b.min.x - origin.x) / direction.x
	t1x := (b.max.x - origin.x) / direction.x
	t0y := (b.min.y - origin.y) / direction.y
	t1y := (b.max.y - origin.y) / direction.y

	if verbosity {
		fmt.Println(b.min.z-origin.z, b.min.z, origin.z)
	}

	t0z := (b.min.z - origin.z) / direction.z
	t1z := (b.max.z - origin.z) / direction.z

	if verbosity {
		fmt.Println(t0x, t1x, t0y, t1y, t0z, t1z)
	}

	tmin := math.Min(t0x, t1x)
	tmax := math.Max(t0x, t1x)
	tminY := math.Min(t0y, t1y)
	tmaxY := math.Max(t0y, t1y)
	tminZ := math.Min(t0z, t1z)
	tmaxZ := math.Max(t0z, t1z)

	if verbosity {
		fmt.Println(tmin, tmax, tminY, tmaxY, tminZ, tmaxZ)
	}

	if (tmin > tmaxY) || (tminY > tmax) {
		return math.Inf(1)
	}

	tmin = math.Max(tmin, tminY)
	tmax = math.Min(tmax, tmaxY)

	if (tmin > tmaxZ) || (tminZ > tmax) {
		return math.Inf(1)
	}

	tmin = math.Max(tmin, tminZ)
	tmax = math.Min(tmax, tmaxZ)

	if tmin < errorDelta {
		if tmax < errorDelta {
			return math.Inf(1)
		}
		return tmax
	}
	if verbosity {
		fmt.Println(tmin)
	}
	return tmin
}

func (b box) directIllumination(l light, point vector, objects []object) color.NRGBA {
	dir := l.posn.sub(point).direction()
	ambientFactor := 0.2
	var col = b.col
	for _, obj := range objects {
		stoppingPoint := obj.obstruct(dir, point, false)
		if stoppingPoint != math.Inf(1) {
			return multiplyNRGBA(col, ambientFactor)
		}
	}

	shader := b.determineColorBrightness(ambientFactor, point, dir)

	return multiplyNRGBA(col, shader)
}

// box, float64, light, vector, vector -> float64
// determines the proportion of lighting that is non-ambient
// calculates the unit normal vector to the sphere at intersection point
// using the direction of incoming light, and normal vector, determines amount of non-ambient
// lighting to utilize.
// returns sum of ambient and non-ambient lighting
func (b box) determineColorBrightness(ambientFactor float64, point, photonDir vector) float64 {
	diffuseFactor := 1 - ambientFactor
	normal := b.findNormal(point)
	shadeFactor := math.Max(0, photonDir.dot(normal))
	colorFactor := (ambientFactor + diffuseFactor*shadeFactor)
	return colorFactor
}

// box, vector -> vector
// returns the normal to the face of a prism at a point.
// attempts to determine on which face of axis-aligned prism the point sits.
// uses the face to calculate the normal.
func (b box) findNormal(point vector) vector {
	normal := makeVector(0, 0, 0) //point.sub(b.center).direction()
	if math.Abs(point.x-b.min.x) < errorDelta {
		normal = makeVector(-1, 0, 0)
	} else if math.Abs(point.x-b.max.x) < errorDelta {
		normal = makeVector(1, 0, 0)
	} else if math.Abs(point.y-b.min.y) < errorDelta {
		normal = makeVector(0, -1, 0)
	} else if math.Abs(point.y-b.max.y) < errorDelta {
		normal = makeVector(0, 1, 0)
	} else if math.Abs(point.z-b.min.z) < errorDelta {
		normal = makeVector(0, 0, -1)
	} else if math.Abs(point.z-b.max.z) < errorDelta {
		normal = makeVector(0, 0, 1)
	} else {
		fmt.Println("Point not placed:", math.Abs(point.x-b.min.x) < errorDelta)
	}
	return normal
}
