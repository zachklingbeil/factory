package one

import (
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	*zero.Zero
	*fx.Fx
}

func NewOne() *One {
	one := &One{
		Zero: zero.NewZero(),
		Fx:   fx.InitFx(),
	}
	return one
}

func (o *One) InitPathless(body zero.One, css, js string) {
	o.BuildPathless(body, css, js)
	o.AddFrame(o.Pathless, o.Router)
}

func (o *One) BuildFrame(name string, elements []zero.One) {
	frame := o.Build(elements)
	final := o.Final(name, frame)
	o.AddFrame(final, o.Router)
}
