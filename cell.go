package ump

type cell struct {
	bodies    map[uint32]*Body
	itemCount int
}

func (c *cell) enter(body *Body) {
	c.bodies[body.ID] = body
	body.cells = append(body.cells, c)
	c.itemCount = len(c.bodies)
}

func (c *cell) leave(body *Body) {
	delete(c.bodies, body.ID)
	c.itemCount = len(c.bodies)
}
