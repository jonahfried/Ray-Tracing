package main

import (
	"fmt"
	"image/color"
	"math"
)

type object interface {
	obstruct(direction, origin vector, verbosity bool) float64
	directIllumination(l light, point vector, objects []object) color.NRGBA
	color() color.NRGBA
}

type sphere struct {
	center vector
	radius float64
	col    color.NRGBA
}

// sphere constructor
func (sphr sphere) color() color.NRGBA {
	return sphr.col
}

// sphere constructor
func makeSphere(x, y, z, r float64, col color.NRGBA) (sphr sphere) {
	sphr.center = makeVector(x, y, z)
	sphr.radius = r
	sphr.col = col
	return sphr
}

// sphere, vector, vector, (verbosity) ->  float64
// determines if a ray will intersect a given sphere, returning the length of ray at point of intersection
// returns infinity if ray does not intersect sphere
func (sphr sphere) obstruct(direction, origin vector, verbosity bool) float64 {
	//	|(s*d+p) - c|^2 = r^2
	//	s = mag
	//	p = starting point
	//	d = unit direction vec
	//	c = center
	//	r = radius sphere
	offset := origin.sub(sphr.center)
	// s^2(d.x^2+d.y^2+d.z^2) + s*2(d.x(offset.x)+d.y(offset.y)+d.z(offset.z))+(|offset|^2-radius^2) = 0
	a := (direction.x * direction.x) + (direction.y * direction.y) + (direction.z * direction.z)
	b := 2 * (direction.x*(offset.x) + direction.y*(offset.y) + direction.z*(offset.z))
	c := (math.Pow((offset.x), 2) + math.Pow((offset.y), 2) + math.Pow((offset.z), 2) - (sphr.radius * sphr.radius))

	// s = (-b ± sqrt(b*b - 4*(a)*(c)))/2a
	discriminant := ((b * b) - (4 * a * c))

	if verbosity {
		fmt.Printf("a:%f,\nb:%f,\nc:%f \n", a, b, c)
	}
	if discriminant < 0 {
		if verbosity {
			fmt.Println(math.Inf(1))
		}
		return math.Inf(1)
	} else if discriminant == 0 {
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
	} else {
		s1 := (-b - math.Sqrt(discriminant)) / (2 * a)
		s2 := (-b + math.Sqrt(discriminant)) / (2 * a)
		// use .000000001 (1e-9), instead of 0,
		// to mitigate rounding errors.
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
			// fmt.Println("test")
			return s
		}
	}
}

// returns a unit ray in a random direction in a hemisphere around given normal
// func sampleHemisphere(n vector, phi, sina, cosa float64) vector {
// 	w := n.direction()
// 	u :=

// 	return theta, phi
// }

type light struct {
	posn vector
}

// determines what color to display at given point
func (sphr sphere) directIllumination(l light, point vector, objects []object) color.NRGBA {
	dir := l.posn.sub(point).direction()
	ambientFactor := 0.2
	var col = sphr.col
	for _, obj := range objects {
		stoppingPoint := obj.obstruct(dir, point, false)
		if stoppingPoint != math.Inf(1) {
			return multiplyNRGBA(col, ambientFactor)
		}
	}
	normal := point.sub(sphr.center).direction()
	diffuseFactor := 1 - ambientFactor
	shadeFactor := math.Max(0, dir.dot(normal))
	colorFactor := (ambientFactor + diffuseFactor*shadeFactor)

	return multiplyNRGBA(col, colorFactor)
}

// type circle struct {
// 	center vector
// 	radius float64
// 	col    color.NRGBA
// }

// // color accessor
// func (crcl circle) color() color.NRGBA {
// 	return crcl.col
// }

// //
// func (crcl circle) obstruct(direction, origin vector, verbosity bool) float64 {
// 	s := crcl.center.y / direction.y
// 	p := direction.mul(s).add(origin)
// 	// fmt.Println(p)
// 	mag := (p.x-crcl.center.x)*(p.x-crcl.center.x) + (p.z-crcl.center.z)*(p.z-crcl.center.z) - crcl.radius*crcl.radius
// 	if mag < errorDelta { //&& mag > -errorDelta {
// 		// fmt.Println(mag)
// 		return s
// 	}
// 	return math.Inf(1)
// }
