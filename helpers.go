package main

import (
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
