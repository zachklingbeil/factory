package element

import "strings"

func simplify(keys []Frame) string {
	var style, script, html strings.Builder

	for _, item := range keys {
		s := string(item.Render())
		for {
			start, end := strings.Index(s, "<style>"), strings.Index(s, "</style>")
			if start < 0 || end <= start {
				break
			}
			style.WriteString(s[start+len("<style>") : end])
			s = s[:start] + s[end+len("</style>"):]
		}
		for {
			start, end := strings.Index(s, "<script>"), strings.Index(s, "</script>")
			if start < 0 || end <= start {
				break
			}
			script.WriteString(s[start+len("<script>") : end])
			s = s[:start] + s[end+len("</script>"):]
		}
		html.WriteString(s)
	}

	var output strings.Builder
	if style.Len() > 0 {
		output.WriteString("<style>")
		output.WriteString(style.String())
		output.WriteString("</style>")
	}
	if script.Len() > 0 {
		output.WriteString("<script>")
		output.WriteString(script.String())
		output.WriteString("</script>")
	}
	output.WriteString(html.String())
	return output.String()
}
