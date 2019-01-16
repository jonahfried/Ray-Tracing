package main

import (
	"fmt"
	"image/color"
	"math"
)

// object is an interface type
// represents a physical mass in the worldspace that interacts with light
type object interface {
	obstruct(direction, origin vector, verbosity bool) float64
	directIllumination(l light, point vector, objects []object) color.NRGBA
	color() color.NRGBA
}

// sphere is a struct type
// implements the object interface
// stores a center, radius, and color
type sphere struct {
	center vector
	radius float64
	col    color.NRGBA
}

// sphr -> color.NRGBA
// returns a given sphere's color
func (sphr sphere) color() color.NRGBA {
	return sphr.col
}

// float64, float64, float64, float64, float64, color.NRGBA -> sphere
// sphere constructor, returning an instance of a sphere struct
func makeSphere(x, y, z, r float64, col color.NRGBA) (sphr sphere) {
	sphr.center = makeVector(x, y, z)
	sphr.radius = r
	sphr.col = col
	return sphr
}

// sphere, vector, vector, (verbosity) ->  float64
//	|(s*d+p) - c|^2 = r^2
// 	s^2(d.x^2+d.y^2+d.z^2) + s*2(d.x(offset.x)+d.y(offset.y)+d.z(offset.z))+(|offset|^2-radius^2) = 0
//	s = mag
//	p = starting point
//	d = unit direction vec
//	c = center
//	r = radius sphere
// analyzing the quadratic equation to determine whether there exists a solution
// to the ray-sphere intersection. If there is an intersection, returns magnitude of ray
// to point of intersection. If there are two solutions, returns the smaller positive solution.
// If there is no intersection returns infinity.
// (accepts a verbosity parameter for testing; prints calculated values)
func (sphr sphere) obstruct(direction, origin vector, verbosity bool) float64 {

	a, b, c := findCoefficients(sphr, direction, origin, verbosity)
	discriminant := findDiscriminant(a, b, c)

	if discriminant < 0 { // NO INTERSECTIONS
		if verbosity {
			fmt.Println(math.Inf(1))
		}
		return math.Inf(1)

	} else if discriminant == 0 { // ONE INTERSECTION
		s := -b / (2 * a)
		if s <= errorDelta {
			if verbosity {
				fmt.Println(math.Inf(1))
			}
			return math.Inf(1)
		}
		if verbosity {
			fmt.Println(s)
		}
		return s

	} else { // TWO INTERSECTIONS
		s1 := (-b - math.Sqrt(discriminant)) / (2 * a)
		s2 := (-b + math.Sqrt(discriminant)) / (2 * a)
		// use .000000001 (1e-9), instead of 0,
		// to prevent rounding errors.
		if s1 <= errorDelta && s2 <= errorDelta {
			if verbosity {
				fmt.Println(math.Inf(1))
			}
			return math.Inf(1)
		} else if s1 <= errorDelta {
			if verbosity {
				fmt.Println(s2)
			}
			return s2
		} else if s2 <= errorDelta {
			if verbosity {
				fmt.Println(s1)
			}
			return s1
		} else {
			s := math.Min(s1, s2)
			if verbosity {
				fmt.Println(s)
			}
			return s
		}
	}
}

// sphere, vector, vector, (verbosity) -> float64, float64, float64
// returns the a, b, and c coefficients of the polynomial
// 	s^2(d.x^2+d.y^2+d.z^2) + s*2(d.x(offset.x)+d.y(offset.y)+d.z(offset.z))+(|offset|^2-radius^2) = 0
// where s is the magnitude of the given ray, and the solutions are
// solutions to the intersection equation
func findCoefficients(sphr sphere, direction, origin vector, verbosity bool) (float64, float64, float64) {
	offset := origin.sub(sphr.center)
	a := (direction.x * direction.x) + (direction.y * direction.y) + (direction.z * direction.z)
	b := 2 * (direction.x*(offset.x) + direction.y*(offset.y) + direction.z*(offset.z))
	c := (math.Pow((offset.x), 2) + math.Pow((offset.y), 2) + math.Pow((offset.z), 2) - (sphr.radius * sphr.radius))

	if verbosity {
		fmt.Printf("a:%f,\nb:%f,\nc:%f \n", a, b, c)
	}
	return a, b, c
}

// float64, float64, float64 -> float64
// returns the discrimanant of a quadratic polynomial
func findDiscriminant(a, b, c float64) float64 {
	// s = (-b ± sqrt(b*b - 4*(a)*(c)))/2a
	return (b * b) - (4 * a * c)
}

// light is a struct type representing the source of lighting in a worldspace
type light struct {
	posn vector
}

// sphr, light, vector, []objects -> color.NRGBA
// determines what color to display at given point:
// Ambiently colors all points some fraction of their color.
// Determines if the light can shine on given point.
// If light hits, determines what additional portion of color to show, and returns scaled color.
func (sphr sphere) directIllumination(l light, point vector, objects []object) color.NRGBA {
	ambientFactor := 0.2 // Cosmetic. Edit based on preference

	dir := l.posn.sub(point).direction() // unit vector from point on sphere to light
	var col = sphr.col
	for _, obj := range objects {
		stoppingPoint := obj.obstruct(dir, point, false)
		if stoppingPoint != math.Inf(1) {
			return multiplyNRGBA(col, ambientFactor)
		}
	}

	shader := sphr.determineColorBrightness(ambientFactor, point, dir)

	return multiplyNRGBA(col, shader)
}

// sphere, float64, light, vector, vector -> float64
// determines the proportion of lighting that is non-ambient
// calculates the unit normal vector to the sphere at intersection point
// using the direction of incoming light, and normal vector, determines amount of non-ambient
// lighting to utilize.
// returns sum of ambient and non-ambient lighting
func (sphr sphere) determineColorBrightness(ambientFactor float64, point, photonDir vector) float64 {
	diffuseFactor := 1 - ambientFactor
	normal := point.sub(sphr.center).direction() // unit vector from sphere center to point on sphere
	shadeFactor := math.Max(0, photonDir.dot(normal))
	colorFactor := (ambientFactor + diffuseFactor*shadeFactor)

	return colorFactor
}
