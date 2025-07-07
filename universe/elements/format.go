package elements

import "html/template"

type format struct{}

func (f *format) Strong(s string) template.HTML { return Tag("strong", s) }
func (f *format) Em(s string) template.HTML     { return Tag("em", s) }
func (f *format) Small(s string) template.HTML  { return Tag("small", s) }
func (f *format) Mark(s string) template.HTML   { return Tag("mark", s) }
func (f *format) Del(s string) template.HTML    { return Tag("del", s) }
func (f *format) Ins(s string) template.HTML    { return Tag("ins", s) }
func (f *format) Sub(s string) template.HTML    { return Tag("sub", s) }
func (f *format) Sup(s string) template.HTML    { return Tag("sup", s) }
func (f *format) Kbd(s string) template.HTML    { return Tag("kbd", s) }
func (f *format) Samp(s string) template.HTML   { return Tag("samp", s) }
func (f *format) Var(s string) template.HTML    { return Tag("var", s) }
func (f *format) Abbr(s string) template.HTML   { return Tag("abbr", s) }
func (f *format) Time(s string) template.HTML   { return Tag("time", s) }
