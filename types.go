package main

import (
	"fmt"
	"image/color"
	"math"
)

type vector struct {
	x, y, z float64
}

func makeVector(x, y, z float64) vector {
	return vector{x, y, z}
}

func (v vector) dot(vec vector) float64 {
	return v.x*vec.x + v.y*vec.y + v.z*vec.z
}

func cross(v1, v2 vector) vector {
	x := (v1.y * v2.z) - (v1.z * v2.y)
	y := (v1.x * v2.z) - (v1.z * v2.x)
	z := (v1.x * v2.y) - (v1.y * v2.x)
	return makeVector(x, y, z)
}

func (v vector) mag() float64 {
	return math.Sqrt(v.dot(v))
}

func (v vector) direction() vector {
	mag := v.mag()
	if mag == 0 {
		return makeVector(0, 0, 0)
	}
	return v.mul(1 / mag)
}

func (v vector) mul(s float64) vector {
	v.x *= s
	v.y *= s
	v.z *= s
	return v
}

func (v vector) addScalar(k float64) vector {
	v.x += k
	v.y += k
	v.z += k
	return v
}

func (v vector) add(vec vector) vector {
	v.x += vec.x
	v.y += vec.y
	v.z += vec.z
	return v
}

func (v vector) sub(vec vector) vector {
	v.x -= vec.x
	v.y -= vec.y
	v.z -= vec.z
	return v
}

func multiplyNRGBA(col color.NRGBA, m float64) color.NRGBA {
	if m < 0 {
		return col
	}
	r := math.Min(255, float64(col.R)*m)
	g := math.Min(255, float64(col.G)*m)
	b := math.Min(255, float64(col.B)*m)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}

type screen struct {
	fov           float64
	aspectRatio   float64
	height, width float64

	// perspectivePoint vector

	pixels [][]color.NRGBA
}

func makeScreen(height, width, fov float64) (scrn screen) {
	scrn.height = height
	scrn.width = width
	scrn.fov = (math.Pi * fov / 180)
	scrn.aspectRatio = (width / height)
	// scrn.perspectivePoint = makeVector(0, (width/2)*math.Tan(scrn.fov/2), 0)

	scrn.pixels = make([][]color.NRGBA, int(widthRes), int(widthRes))
	for i := 0; i < int(widthRes); i++ {
		scrn.pixels[i] = make([]color.NRGBA, int(heightRes), int(heightRes))
		for j := 0; j < int(heightRes); j++ {
			scrn.pixels[i][j] = black
		}
	}
	return scrn
}

func (scrn screen) at(x, y int) color.Color {
	return scrn.pixels[x][y]
}

func (scrn screen) point(x, y int) vector {
	i := float64(x) / widthRes
	j := float64(y) / heightRes

	xdir := (2*i - 1) * math.Tan(scrn.fov/2) * scrn.aspectRatio
	zdir := (1 - (2 * j)) * math.Tan(scrn.fov/2)

	// v := scrn.topLeft.add(xdir.mul(i))
	// v = v.add(ydir.mul(j))

	v := makeVector(xdir, 0, zdir)
	return v
}

type camera struct {
	posn vector
}

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

func (sphr sphere) color() color.NRGBA {
	return sphr.col
}

type circle struct {
	center vector
	radius float64
	col    color.NRGBA
}

func (crcl circle) color() color.NRGBA {
	return crcl.col
}

func makeSphere(x, y, z, r float64, col color.NRGBA) (sphr sphere) {
	sphr.center = makeVector(x, y, z)
	sphr.radius = r
	sphr.col = col
	return sphr
}

func (crcl circle) obstruct(direction, origin vector, verbosity bool) float64 {
	s := crcl.center.y / direction.y
	p := direction.mul(s).add(origin)
	// fmt.Println(p)
	mag := (p.x-crcl.center.x)*(p.x-crcl.center.x) + (p.z-crcl.center.z)*(p.z-crcl.center.z) - crcl.radius*crcl.radius
	if mag < errorDelta { //&& mag > -errorDelta {
		// fmt.Println(mag)
		return s
	}
	return math.Inf(1)
}

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

// func (sphr sphere) obstruct(direction, origin vector) float64 {
// 	diff := sphr.center.sub(origin)
// 	v := diff.dot(direction)
// 	disc := (sphr.radius * sphr.radius) - (diff.dot(diff) - (v * v))
// 	if disc < 0 {
// 		return -1
// 	}
// 	return v - math.Sqrt(disc)
// }

type light struct {
	posn vector
}

func (sphr sphere) directIllumination(l light, point vector, objects []object) color.NRGBA {
	dir := l.posn.sub(point).direction()
	ambientFactor := 0.1
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
