package elements

import "html/template"

type other struct{}

func (o *other) Nav(attrs map[string]string) template.HTML {
	return ClosedTag("nav", attrs)
}

func (o *other) Button(label string) template.HTML {
	return Tag("button", label)
}
func (o *other) Code(code string) template.HTML {
	return Tag("code", code)
}
func (o *other) Table(cols uint8, rows uint64, data [][]string) template.HTML {
	table := "<table>"
	for _, row := range data {
		table += "<tr>"
		for _, cell := range row {
			table += string(Tag("td", cell))
		}
		table += "</tr>"
	}
	table += "</table>"
	return template.HTML(table)
}
