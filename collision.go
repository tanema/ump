package ump

type (
	// Point is used as data points in a collision. This could represent a touch
	// or a normal or something else.
	Point struct {
		X, Y float32
	}
	// Collision represents a touch of two objects. The intersection is the fraction
	// of the movement where the two items touched. Distanse is how far the two object
	// are from each other. Body is the other body that was collided.
	//
	// - Move is the amount the body needs to move to resolve the collision
	// - Normal is the normal of the collision
	// - Touch is the point where the two bodies first touched
	// - Data is used for misc data that can be pass in a filter
	// - RespType describes which filter was used to resolve the collision
	Collision struct {
		Intersection float32
		Distance     float32
		Body         *Body
		Move         Point
		Normal       Point
		Touch        Point
		Data         Point
		RespType     string
	}
)

type collisionsByDistance []*Collision

func (s collisionsByDistance) Len() int {
	return len(s)
}

func (s collisionsByDistance) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s collisionsByDistance) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Intersection == b.Intersection {
		return a.Distance < b.Distance
	}
	return a.Intersection < b.Intersection
}
