package universe

import (
	"github.com/zachklingbeil/factory/universe/constant"
	"github.com/zachklingbeil/factory/universe/element"
)

type Universe struct {
	Element  *element.Element
	Constant *constant.Head
}

func New() *Universe {
	return &Universe{
		Element:  element.NewElements(),
		Constant: constant.NewHead(),
	}
}
