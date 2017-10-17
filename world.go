package ump

import (
	"math"
	"sort"
)

const defaultFilter = "slide"

type (
	// World is the virtual world in which all these collisions happen. The world
	// contains a grid, which contains several cells, which contains collidable bodies.
	//
	// A world also has registered responses to filter collisions please see Resp for this.
	World struct {
		grid      *grid
		responses map[string]Resp
	}
	// Resp is a function that will handle and resolve a collision. For instance
	// the bound filter will return the bounce goal gx gy, and then project for the
	// new direction to make sure there are not collisions in that direction.
	Resp func(world *World, col *Collision, body *Body, goalX, goalY float32) (gx, gy float32, cols []*Collision)
)

// NewWorld builds a physics world with the provided cell size. A good default is 64.
// It represents the size of the sides of the (squared) cells that will be used
// internally to provide the data. In tile based games, it's usually a multiple
// of the tile side size. So in a game where tiles are 32x32, cellSize will be 32,
// 64 or 128. In more sparse games, it can be higher.
func NewWorld(cellSize int) *World {
	world := &World{
		grid:      newGrid(cellSize),
		responses: map[string]Resp{},
	}
	world.AddResponse("touch", touchFilter)
	world.AddResponse("cross", crossFilter)
	world.AddResponse("slide", slideFilter)
	world.AddResponse("bounce", bounceFilter)
	return world
}

// Add will create a new Body to be tracked in this world. The tag is important.
// It is used to decided which filter to use but also for you to decide what to
// do when a body collides with your new body.
//
// left, top, w, and h describe a rectangle for the body to inhabit.
func (world *World) Add(tag string, left, top, w, h float32) *Body {
	return newBody(world, tag, left, top, w, h)
}

// QueryRect will take the rectangle arguments and return any bodies that are in
// that rectangle
//
// If tags are passed into the query then it will only return the bodies with those
// tags.
func (world *World) QueryRect(x, y, w, h float32, tags ...string) []*Body {
	return world.getBodiesInCells(world.grid.cellsInRect(x, y, w, h), tags...)
}

// QueryPoint will return any bodies that are underneathe the point.
//
// If tags are passed into the query then it will only return the bodies with those
// tags.
func (world *World) QueryPoint(x, y float32, tags ...string) []*Body {
	bodies := []*Body{}
	c := world.grid.cellAt(x, y, false)
	if c == nil {
		return []*Body{}
	}
	for _, body := range c.bodies {
		if body.HasTag(tags...) && body.containsPoint(x, y) {
			bodies = append(bodies, body)
		}
	}
	return bodies
}

// QuerySegment will return any bodies that are underneathe the segment/line.
//
// If tags are passed into the query then it will only return the bodies with those
// tags.
func (world *World) QuerySegment(x1, y1, x2, y2 float32, tags ...string) []*Body {
	bodies := []*Body{}
	visited := map[*Body]bool{}
	cells := world.grid.getCellsTouchedBySegment(x1, y1, x2, y2)
	bodiesOnSegment := world.getBodiesInCells(cells)
	distances := map[uint32]float32{}
	for _, body := range bodiesOnSegment {
		if _, ok := visited[body]; !ok && body.HasTag(tags...) {
			visited[body] = true
			fraction, _, _ := body.getRayIntersectionFraction(x1, y1, x2-x1, y2-y1)
			if fraction != inf {
				bodies = append(bodies, body)
				distances[body.ID] = fraction
			}
		}
	}

	bodiesBy(func(b1, b2 *Body) bool {
		return distances[b1.ID] < distances[b2.ID]
	}).Sort(bodies)

	return bodies
}

func (world *World) getBodiesInCells(cells []*cell, tags ...string) []*Body {
	dict := make(map[uint32]bool)
	bodies := []*Body{}
	for _, c := range cells {
		for id, body := range c.bodies {
			if _, ok := dict[id]; !ok && body.HasTag(tags...) {
				bodies = append(bodies, body)
				dict[id] = true
			}
		}
	}
	return bodies
}

// Project will project the goal location of the provided body but not move it.
// This is good for checking a future location of a body and see if there are any
// collisions in that space.
func (world *World) Project(body *Body, goalX, goalY float32) []*Collision {
	collisions := []*Collision{}

	tl := float32(math.Min(float64(goalX), float64(body.x)))
	tt := float32(math.Min(float64(goalY), float64(body.y)))
	tr := float32(math.Max(float64(goalX+body.w), float64(body.x+body.w)))
	tb := float32(math.Max(float64(goalY+body.h), float64(body.y+body.h)))

	visited := map[*Body]bool{}
	bodies := world.getBodiesInCells(world.grid.cellsInRect(tl, tt, tr-tl, tb-tt))
	for _, other := range bodies {
		if _, ok := visited[other]; !ok {
			visited[other] = true
			if col := body.collide(other, goalX, goalY); col != nil {
				collisions = append(collisions, col)
			}
		}
	}

	sort.Sort(collisionsByDistance(collisions))

	return collisions
}

// AddResponse will add a new filter response for this world. This is helpful if
// you are creating custom reactions in your world.
func (world *World) AddResponse(name string, response Resp) {
	world.responses[name] = response
}
