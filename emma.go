package main

import "image/color"

func e(col color.NRGBA) []object {
	bx1 := makeBox2(
		makeVector(-5.5, 100, 7),
		2, 1, 8,
		col,
	)
	bx2 := makeBox2(
		makeVector(-5.5, 100, 5),
		2, 1, 8,
		col,
	)
	bx3 := makeBox2(
		makeVector(-5.5, 100, 3),
		2, 1, 8,
		col,
	)
	bx4 := makeBox2(
		makeVector(-7, 100, 5),
		1, 5, 8,
		col,
	)

	objs := make([]object, 0)
	objs = append(objs, bx1)
	objs = append(objs, bx2)
	objs = append(objs, bx3)
	objs = append(objs, bx4)
	return objs
}

func m(posn vector, col color.NRGBA) []object {
	objs := make([]object, 0)

	bx1 := makeBox2(
		posn,
		1, 3, 8,
		col,
	)
	bx2 := makeBox2(
		posn.sub(makeVector(2, 0, 0)),
		1, 3, 8,
		col,
	)
	bx3 := makeBox2(
		posn.add(makeVector(2, 0, 0)),
		1, 3, 8,
		col,
	)
	sphr1 := makeSphere(posn.x-1, posn.y-5, posn.z+1.5, .7, col)
	sphr2 := makeSphere(posn.x+1, posn.y-5, posn.z+1.5, .7, col)

	objs = append(objs, bx1)
	objs = append(objs, bx2)
	objs = append(objs, bx3)
	objs = append(objs, sphr1)
	objs = append(objs, sphr2)

	return objs
}

func a(posn vector) []object {
	objs := make([]object, 0)

	bx1 := makeBox2(
		posn.add(makeVector(0, 0, .3)),
		1, .75, 8,
		pink,
	)
	bx2 := makeBox2(
		posn.sub(makeVector(.75, 0, 0)),
		.75, 3, 8,
		pink,
	)
	bx3 := makeBox2(
		posn.add(makeVector(.75, 0, 0)),
		.75, 3, 8,
		pink,
	)

	sphr1 := makeSphere(posn.x, posn.y-3.2, posn.z+2, 1, pink)

	objs = append(objs, bx1)
	objs = append(objs, bx2)
	objs = append(objs, bx3)
	objs = append(objs, sphr1)
	return objs
}

func emma() []object {
	objs := make([]object, 0)
	objs = append(objs, e(purple)...)
	objs = append(objs, m(makeVector(5.5, 100, 5), orange)...)
	objs = append(objs, m(makeVector(-5.5, 100, -5), orange)...)
	objs = append(objs, a(makeVector(5.5, 100, -5))...)
	return objs
}
