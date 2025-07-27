package zero

type Universe interface {
	AddCoord(x, y int, z any)
	GetX() []*Coordinate
	GetY(x int) []*Coordinate
	GetZ(x, y int) *Coordinate
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z any `json:"z"`
}

func NewUniverse() Universe {
	return &coordinates{
		coords: make(map[int]map[int]any),
	}
}

type coordinates struct {
	coords map[int]map[int]any
}

func (c *coordinates) AddCoord(x, y int, z any) {
	if _, exists := c.coords[x]; !exists {
		c.coords[x] = make(map[int]any)
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

func (c *coordinates) GetY(x int) []*Coordinate {
	result := make([]*Coordinate, 0)
	if row, exists := c.coords[x]; exists {
		for y, z := range row {
			result = append(result, &Coordinate{X: x, Y: y, Z: z})
		}
	}
	return result
}

func (c *coordinates) GetZ(x, y int) *Coordinate {
	if row, exists := c.coords[x]; exists {
		if z, ok := row[y]; ok {
			return &Coordinate{X: x, Y: y, Z: z}
		}
	}
	return &Coordinate{}
}
