package io

type Coordinates interface {
	Add(x, y string, z any)
	GetX() []string
	GetY(x string) []string
	GetZ(x, y string) any
}

func NewCoordinates() Coordinates {
	return &coordinates{
		coords: make(map[string]map[string]any),
	}
}

type coordinates struct {
	coords map[string]map[string]any
}

func (c *coordinates) Add(x, y string, z any) {
	if _, exists := c.coords[x]; !exists {
		c.coords[x] = make(map[string]any)
	}
	c.coords[x][y] = z
}

func (c *coordinates) GetX() []string {
	xs := make([]string, 0, len(c.coords))
	for x := range c.coords {
		xs = append(xs, x)
	}
	return xs
}

func (c *coordinates) GetY(x string) []string {
	ys := make([]string, 0)
	if row, exists := c.coords[x]; exists {
		for y := range row {
			ys = append(ys, y)
		}
	}
	return ys
}

func (c *coordinates) GetZ(x, y string) any {
	if row, exists := c.coords[x]; exists {
		if z, ok := row[y]; ok {
			return z
		}
	}
	return nil
}
