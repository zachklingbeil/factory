package zero

import (
	"html/template"
	"os"
	"regexp"
	"strings"
)

func (f *zero) BuildFrameFromHTMLFile() error {
	content, err := os.ReadFile("./embed/coordinate.html")
	if err != nil {
		return err
	}
	htmlStr := string(content)

	// Regex patterns to extract full <style>, <template>, <script> blocks (including tags)
	styleRe := regexp.MustCompile(`(?s)<style.*?>.*?</style>`)
	templateRe := regexp.MustCompile(`(?s)<template.*?>.*?</template>`)
	scriptRe := regexp.MustCompile(`(?s)<script.*?>.*?</script>`)

	var b strings.Builder

	// Append <style> blocks
	for _, match := range styleRe.FindAllString(htmlStr, -1) {
		b.WriteString(match)
	}

	// Append <template> blocks
	for _, match := range templateRe.FindAllString(htmlStr, -1) {
		b.WriteString(match)
	}

	// Append <script> blocks
	for _, match := range scriptRe.FindAllString(htmlStr, -1) {
		b.WriteString(match)
	}

	final := One(template.HTML(b.String()))
	f.frames = append(f.frames, &final)
	f.count++
	return nil
}

// type Coord struct {
// 	X int
// 	Y int
// 	Z struct {
// 		Peer  string
// 		Time  string
// 		Value string
// 	}
// }

// type PlaneRow struct {
// 	Negatives []Coord
// 	YInt      *Coord
// 	Positives []Coord
// 	RowIndex  int
// }

// func (f *zero) CoordinatePlane(coords []Coord) {
// 	var b strings.Builder
// 	b.WriteString(`<style>`)
// 	b.WriteString(f.CoordinateCSS())
// 	b.WriteString(`</style>`)
// 	b.WriteString(renderPlane(coords))
// 	final := One(template.HTML(b.String()))
// 	f.frames = append(f.frames, &final)
// 	f.count++
// }

// func renderPlane(coords []Coord) string {
// 	rows := groupCoordsByRow(coords)
// 	var b strings.Builder
// 	b.WriteString(`<div class="coordinate-plane" id="coordinate-plane">`)
// 	for _, row := range rows {
// 		b.WriteString(renderPlaneRow(row))
// 	}
// 	b.WriteString(`</div>`)
// 	return b.String()
// }

// func groupCoordsByRow(coords []Coord) []PlaneRow {
// 	rowMap := make(map[int]*PlaneRow)
// 	for _, c := range coords {
// 		if rowMap[c.Y] == nil {
// 			rowMap[c.Y] = &PlaneRow{RowIndex: c.Y}
// 		}
// 		switch {
// 		case c.X < 0:
// 			rowMap[c.Y].Negatives = append(rowMap[c.Y].Negatives, c)
// 		case c.X == 0:
// 			rowMap[c.Y].YInt = &c
// 		case c.X > 0:
// 			rowMap[c.Y].Positives = append(rowMap[c.Y].Positives, c)
// 		}
// 	}
// 	var rows []PlaneRow
// 	for _, row := range rowMap {
// 		sort.Slice(row.Negatives, func(i, j int) bool { return row.Negatives[i].X < row.Negatives[j].X })
// 		sort.Slice(row.Positives, func(i, j int) bool { return row.Positives[i].X < row.Positives[j].X })
// 		rows = append(rows, *row)
// 	}
// 	sort.Slice(rows, func(i, j int) bool { return rows[i].RowIndex < rows[j].RowIndex })
// 	return rows
// }

// func renderPlaneRow(row PlaneRow) string {
// 	return fmt.Sprintf(
// 		`<div class="row">%s%s%s</div>`,
// 		renderAxis("negative", row.Negatives),
// 		renderYInt(row.YInt, row.RowIndex),
// 		renderAxis("positive", row.Positives),
// 	)
// }

// func renderAxis(axisType string, coords []Coord) string {
// 	var b strings.Builder
// 	b.WriteString(fmt.Sprintf(`<div class="axis %s">`, axisType))
// 	for _, c := range coords {
// 		b.WriteString(renderCoordinate(c))
// 	}
// 	b.WriteString(`</div>`)
// 	return b.String()
// }

// func renderYInt(yint *Coord, rowIndex int) string {
// 	if yint != nil {
// 		return fmt.Sprintf(`<span class="yint">%s</span>`, template.HTMLEscapeString(yint.Z.Value))
// 	}
// 	return fmt.Sprintf(`<span class="yint">%d</span>`, rowIndex)
// }

// func renderCoordinate(c Coord) string {
// 	axisType := "label"
// 	if c.X < 0 {
// 		axisType = "negative"
// 	} else if c.X > 0 {
// 		axisType = "positive"
// 	}
// 	return fmt.Sprintf(
// 		`<div class="coordinate %s"><div>%s</div><div>%s</div><div>%s</div></div>`,
// 		axisType,
// 		template.HTMLEscapeString(c.Z.Peer),
// 		template.HTMLEscapeString(c.Z.Time),
// 		template.HTMLEscapeString(c.Z.Value),
// 	)
// }
