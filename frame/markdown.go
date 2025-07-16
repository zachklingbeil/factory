package frame

import (
	"html/template"
	"os"
	"regexp"
	"strings"
)

type Markdown struct {
	H1         template.CSS
	H2         template.CSS
	H3         template.CSS
	Spacing    template.CSS
	PSpacing   template.CSS
	Code       template.CSS
	InlineCode template.CSS
	Padding    template.CSS
	Img        template.CSS
	Img2       template.CSS
	Img3       template.CSS
	Img4       template.CSS
	Tables     template.CSS
	Links      template.CSS
}

// Compiled regex patterns for better performance
var (
	bold   = regexp.MustCompile(`\*\*(.*?)\*\*|__(.*?)__`)
	italic = regexp.MustCompile(`(?:\*([^*]+)\*|_([^_]+)_)`)
	code   = regexp.MustCompile("`([^`]+)`")
	link   = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	img    = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)(?:\s+"([^"]*)")?\)`)
)

func (f *Frame) MarkdownToHTML(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}

	markdown := string(content)
	if markdown == "" {
		return template.HTML("")
	}

	lines := strings.Split(markdown, "\n")
	mdElements := make([]template.HTML, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "######"):
			mdElements = append(mdElements, f.H6(strings.TrimSpace(line[6:])))
		case strings.HasPrefix(line, "#####"):
			mdElements = append(mdElements, f.H5(strings.TrimSpace(line[5:])))
		case strings.HasPrefix(line, "####"):
			mdElements = append(mdElements, f.H4(strings.TrimSpace(line[4:])))
		case strings.HasPrefix(line, "###"):
			mdElements = append(mdElements, f.H3(strings.TrimSpace(line[3:])))
		case strings.HasPrefix(line, "##"):
			mdElements = append(mdElements, f.H2(strings.TrimSpace(line[2:])))
		case strings.HasPrefix(line, "#"):
			mdElements = append(mdElements, f.H1(strings.TrimSpace(line[1:])))
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* "):
			listItems, newIndex := f.parseList(lines, i)
			mdElements = append(mdElements, f.List(listItems, false))
			i = newIndex
		default:
			formatted := f.processInlineFormatting(line)
			mdElements = append(mdElements, f.Paragraph(formatted))
		}
	}

	// Combine markdown elements with any additional elements passed in
	allElements := append(mdElements, elements...)
	return f.CreateFrame(allElements...)
}

// parseList extracts list items starting from the given index
func (f *Frame) parseList(lines []string, startIndex int) ([]any, int) {
	listItems := make([]any, 0, 10)
	i := startIndex

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "- ") && !strings.HasPrefix(line, "* ") {
			break
		}

		// Remove list marker and process inline formatting
		item := strings.TrimSpace(line[2:])
		listItems = append(listItems, f.processInlineFormatting(item))
		i++
	}

	return listItems, i - 1
}

// processInlineFormatting applies inline markdown formatting
func (f *Frame) processInlineFormatting(text string) string {
	// Images
	text = img.ReplaceAllStringFunc(text, func(match string) string {
		matches := img.FindStringSubmatch(match)
		if len(matches) >= 3 {
			return string(f.Img(matches[2], matches[1], "50vw", "auto"))
		}
		return match
	})

	// Links
	text = link.ReplaceAllStringFunc(text, func(match string) string {
		matches := link.FindStringSubmatch(match)
		if len(matches) == 3 {
			return string(f.Link(matches[2], matches[1]))
		}
		return match
	})

	// Bold
	text = bold.ReplaceAllStringFunc(text, func(match string) string {
		content := match[2 : len(match)-2]
		return string(f.Strong(content))
	})

	// Italic
	text = italic.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		return string(f.Em(content))
	})

	// Code
	text = code.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		return string(f.Code(content))
	})

	return text
}
