package zero

type Universe interface {
	AddCoord(x, y string, z any)
	GetX() []*Coordinate
	GetY(x string) []*Coordinate
	GetZ(x, y string) *Coordinate
}

type Coordinate struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z any    `json:"z"`
}

func NewUniverse() Universe {
	return &coordinates{
		coords: make(map[string]map[string]any),
	}
}

type coordinates struct {
	coords map[string]map[string]any
}

func (c *coordinates) AddCoord(x, y string, z any) {
	if _, exists := c.coords[x]; !exists {
		c.coords[x] = make(map[string]any)
	}
	c.coords[x][y] = z
}

func (c *coordinates) GetX() []*Coordinate {
	result := make([]*Coordinate, 0)
	for x, row := range c.coords {
		for y, z := range row {
			result = append(result, &Coordinate{X: x, Y: y, Z: z})
		}
	}
	return result
}

func (c *coordinates) GetY(x string) []*Coordinate {
	result := make([]*Coordinate, 0)
	if row, exists := c.coords[x]; exists {
		for y, z := range row {
			result = append(result, &Coordinate{X: x, Y: y, Z: z})
		}
	}
	return result
}

func (c *coordinates) GetZ(x, y string) *Coordinate {
	if row, exists := c.coords[x]; exists {
		if z, ok := row[y]; ok {
			return &Coordinate{X: x, Y: y, Z: z}
		}
	}
	return &Coordinate{}
}
