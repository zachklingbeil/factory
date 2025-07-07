package universe

import (
	"github.com/zachklingbeil/factory/universe/elements"
	"github.com/zachklingbeil/factory/universe/frame"
)

type Universe struct {
	Elements *elements.Elements
	Frame    *frame.Component
}
