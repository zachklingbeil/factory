package pathless

import (
	"html/template"
	"regexp"
	"strings"
)

// Compiled regex patterns for better performance
var (
	bold   = regexp.MustCompile(`\*\*(.*?)\*\*|__(.*?)__`)
	italic = regexp.MustCompile(`(?:\*([^*]+)\*|_([^_]+)_)`)
	code   = regexp.MustCompile("`([^`]+)`")
	link   = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	img    = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)(?:\s+"([^"]*)")?\)`)
)

// MarkdownToHTML converts markdown text to HTML using existing methods
func (p *Pathless) MarkdownToHTML(markdown string) template.HTML {
	if markdown == "" {
		return template.HTML("")
	}

	lines := strings.Split(markdown, "\n")
	elements := make([]template.HTML, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "######"):
			elements = append(elements, p.H6(strings.TrimSpace(line[6:])))
		case strings.HasPrefix(line, "#####"):
			elements = append(elements, p.H5(strings.TrimSpace(line[5:])))
		case strings.HasPrefix(line, "####"):
			elements = append(elements, p.H4(strings.TrimSpace(line[4:])))
		case strings.HasPrefix(line, "###"):
			elements = append(elements, p.H3(strings.TrimSpace(line[3:])))
		case strings.HasPrefix(line, "##"):
			elements = append(elements, p.H2(strings.TrimSpace(line[2:])))
		case strings.HasPrefix(line, "#"):
			elements = append(elements, p.H1(strings.TrimSpace(line[1:])))
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* "):
			listItems, newIndex := p.parseList(lines, i)
			elements = append(elements, p.List(listItems, false))
			i = newIndex
		default:
			formatted := p.processInlineFormatting(line)
			elements = append(elements, p.Paragraph(formatted))
		}
	}

	return p.combineElements(elements)
}

// parseList extracts list items starting from the given index
func (p *Pathless) parseList(lines []string, startIndex int) ([]any, int) {
	listItems := make([]any, 0, 10)
	i := startIndex

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "- ") && !strings.HasPrefix(line, "* ") {
			break
		}

		// Remove list marker and process inline formatting
		item := strings.TrimSpace(line[2:])
		listItems = append(listItems, p.processInlineFormatting(item))
		i++
	}

	return listItems, i - 1
}

// processInlineFormatting applies inline markdown formatting
func (p *Pathless) processInlineFormatting(text string) string {
	// Images
	text = img.ReplaceAllStringFunc(text, func(match string) string {
		matches := img.FindStringSubmatch(match)
		if len(matches) >= 3 {
			return string(p.Img(matches[2], matches[1], "", ""))
		}
		return match
	})

	// Links
	text = link.ReplaceAllStringFunc(text, func(match string) string {
		matches := link.FindStringSubmatch(match)
		if len(matches) == 3 {
			return string(p.Link(matches[2], matches[1]))
		}
		return match
	})

	// Bold
	text = bold.ReplaceAllStringFunc(text, func(match string) string {
		content := match[2 : len(match)-2]
		return string(p.Strong(content))
	})

	// Italic
	text = italic.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		return string(p.Em(content))
	})

	// Code
	text = code.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		return string(p.Code(content))
	})

	return text
}

// combineElements efficiently combines HTML elements
func (p *Pathless) combineElements(elements []template.HTML) template.HTML {
	if len(elements) == 0 {
		return template.HTML("")
	}

	var result strings.Builder
	for _, element := range elements {
		result.WriteString(string(element))
	}

	return template.HTML(result.String())
}
