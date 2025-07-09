package universe

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MarkdownToHTML converts markdown text to HTML using existing methods
func MarkdownToHTML(markdown string) template.HTML {
	var elements []template.HTML
	lines := strings.Split(markdown, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		switch {
		case strings.HasPrefix(line, "# "):
			elements = append(elements, H1(strings.TrimPrefix(line, "# ")))
		case strings.HasPrefix(line, "## "):
			elements = append(elements, H2(strings.TrimPrefix(line, "## ")))
		case strings.HasPrefix(line, "### "):
			elements = append(elements, H3(strings.TrimPrefix(line, "### ")))
		case strings.HasPrefix(line, "#### "):
			elements = append(elements, H4(strings.TrimPrefix(line, "#### ")))
		case strings.HasPrefix(line, "##### "):
			elements = append(elements, H5(strings.TrimPrefix(line, "##### ")))
		case strings.HasPrefix(line, "###### "):
			elements = append(elements, H6(strings.TrimPrefix(line, "###### ")))
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* "):
			// Handle lists (simplified)
			listItems := []any{}
			for i < len(lines) && (strings.HasPrefix(strings.TrimSpace(lines[i]), "- ") || strings.HasPrefix(strings.TrimSpace(lines[i]), "* ")) {
				item := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(lines[i], "- "), "* "))
				listItems = append(listItems, item)
				i++
			}
			i-- // Back up one since the loop will increment
			elements = append(elements, List(listItems, false))
		case line != "":
			// Regular paragraph with inline formatting
			formatted := processInlineFormatting(line)
			elements = append(elements, template.HTML("<p>"+formatted+"</p>"))
		}
	}

	// Combine all elements
	var result strings.Builder
	for _, element := range elements {
		result.WriteString(string(element))
	}

	return template.HTML(result.String())
}

// Process inline formatting like **bold**, *italic*, `code`, etc.
func processInlineFormatting(text string) string {
	// Bold: **text** or __text__
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*|__(.*?)__`)
	text = boldRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(strings.Trim(match, "*"), "_")
		return string(Strong(content))
	})

	// Italic: *text* or _text_
	italicRegex := regexp.MustCompile(`\*(.*?)\*|_(.*?)_`)
	text = italicRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(strings.Trim(match, "*"), "_")
		return string(Em(content))
	})

	// Inline code: `text`
	codeRegex := regexp.MustCompile("`(.*?)`")
	text = codeRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "`")
		return string(Code(content))
	})

	// Links: [text](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := linkRegex.FindStringSubmatch(match)
		if len(matches) == 3 {
			return string(Link(matches[2], matches[1]))
		}
		return match
	})

	return text
}

// LoadMarkdownFile reads a markdown file and creates a Universe page
func (u *Universe) LoadMarkdownFile(pageName, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	html := MarkdownToHTML(string(content))
	u.CreateFrame(pageName, html)
	return nil
}

// LoadMarkdownDirectory loads all markdown files from a directory
func (u *Universe) LoadMarkdownDirectory(dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			// Use filename without extension as page name
			pageName := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))

			err := u.LoadMarkdownFile(pageName, path)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// LoadMarkdownWithName loads markdown with a custom page name
func (u *Universe) LoadMarkdownWithName(pageName, filePath string) error {
	return u.LoadMarkdownFile(pageName, filePath)
}
