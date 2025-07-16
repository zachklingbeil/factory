package frame

import (
	"bytes"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

type Frame struct {
	Md *goldmark.Markdown
}

func NewFrame() *Frame {
	return &Frame{
		Md: initGoldmark(),
	}
}

// Compiled regex patterns for better performance
var (
	bold   = regexp.MustCompile(`\*\*(.*?)\*\*|__(.*?)__`)
	italic = regexp.MustCompile(`(?:\*([^*]+)\*|_([^_]+)_)`)
	code   = regexp.MustCompile("`([^`]+)`")
	link   = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	img    = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)(?:\s+"([^"]*)")?\)`)
)

func (f *Frame) CreateFrame(elements ...template.HTML) template.HTML {
	if len(elements) == 0 {
		return template.HTML("")
	}

	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}

	return template.HTML(builder.String())
}

func (f *Frame) AddCSS(styles map[string]string) template.HTML {
	var builder strings.Builder
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	return template.HTML(builder.String())
}

func (f *Frame) AddJS(js string) template.HTML {
	var builder strings.Builder
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	return template.HTML(builder.String())
}

func (f *Frame) FromMarkdown(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		return template.HTML("")
	}
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, template.HTML(buf.String()))
	allElements = append(allElements, elements...)
	return f.CreateFrame(allElements...)
}

func (f *Frame) MarkdownToHTML(file string) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}

	markdown := string(content)
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
			elements = append(elements, f.H6(strings.TrimSpace(line[6:])))
		case strings.HasPrefix(line, "#####"):
			elements = append(elements, f.H5(strings.TrimSpace(line[5:])))
		case strings.HasPrefix(line, "####"):
			elements = append(elements, f.H4(strings.TrimSpace(line[4:])))
		case strings.HasPrefix(line, "###"):
			elements = append(elements, f.H3(strings.TrimSpace(line[3:])))
		case strings.HasPrefix(line, "##"):
			elements = append(elements, f.H2(strings.TrimSpace(line[2:])))
		case strings.HasPrefix(line, "#"):
			elements = append(elements, f.H1(strings.TrimSpace(line[1:])))
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* "):
			listItems, newIndex := f.parseList(lines, i)
			elements = append(elements, f.List(listItems, false))
			i = newIndex
		default:
			formatted := f.processInlineFormatting(line)
			elements = append(elements, f.Paragraph(formatted))
		}
	}

	return f.CreateFrame(elements...)
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
			return string(f.Img(matches[2], matches[1], "", ""))
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

// combineElements efficiently combines HTML elements
func (f *Frame) combineElements(elements []template.HTML) template.HTML {
	if len(elements) == 0 {
		return template.HTML("")
	}

	var result strings.Builder
	for _, element := range elements {
		result.WriteString(string(element))
	}

	return template.HTML(result.String())
}
