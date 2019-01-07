package main

type plane struct {
	topLeft     vector
	topRight    vector
	bottomLeft  vector
	bottomRight vector

	normal vector
}

func makePlane(topLeft, topRight, bottomLeft, bottomRight vector) (p plane) {
	p.topLeft = topLeft
	p.topRight = topRight
	p.bottomLeft = bottomLeft
	p.bottomRight = bottomRight

	p.normal = cross(topRight.sub(topLeft), topRight.sub(bottomRight))

	return p
}

func (p plane) obstruct(direction, origin vector) float64 {
	d := p.normal.dot(p.topRight)
	// origin.add(direction.mul(s)) = d
	// direction.x*s*p.normal.x + direction.y*s*p.normal.y + direction.z*s*p.normal.z = d - origin.x - origin.y - origin.z
	denominator := (direction.x*p.normal.x + direction.y*p.normal.y + direction.z*p.normal.z)
	if denominator == 0 {
		return -1
	}
	s := (d - origin.x - origin.y - origin.z) / (direction.x*p.normal.x + direction.y*p.normal.y + direction.z*p.normal.z)
	// point := direction.mul(s).add(origin)
	return s
}
