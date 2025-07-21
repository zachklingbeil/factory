package zero

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type Frame interface {
	Build(elements []One) One
	Final(class string, elements []One) One
	JS(js string) One
	CSS(css string) One
	AddKeybind(containerId string, keyHandlers map[string]string) One
	AddScrollKeybinds() One
}

// --- frame Implementation ---
type frame struct{}

func NewFrame() Frame {
	return &frame{}
}

func (f *frame) Build(elements []One) One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(el))
	}
	return One(template.HTML(b.String()))
}

func (f *frame) Final(class string, elements []One) One {
	content := string(f.Build(elements))
	return One(template.HTML(fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), content)))
}

func (f *frame) JS(js string) One {
	return One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
}

func (f *frame) CSS(css string) One {
	return One(template.HTML(fmt.Sprintf(`<style>%s</style>`, css)))
}

func (f *frame) AddKeybind(containerId string, keyHandlers map[string]string) One {
	var handlers strings.Builder
	for key, handlerCode := range keyHandlers {
		handlers.WriteString(fmt.Sprintf(`
         if (event.key === %q) {
            %s
         }
        `, key, handlerCode))
	}
	js := fmt.Sprintf(`
document.addEventListener('DOMContentLoaded', () => {
   const container = document.getElementById(%q);
   if (!container) return;
   container.tabIndex = 0;
   container.addEventListener('keydown', (event) => {
      %s
   });
});
`, containerId, handlers.String())
	return One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
}

func (f *frame) AddScrollKeybinds() One {
	return f.JS(
		`document.addEventListener('keydown', function(event) {
            const c = document.getElementById('one');
            if (!c) return;
            if (event.key === 'w') {
                c.scrollBy({ top: -100, behavior: 'smooth' });
            }
            if (event.key === 's') {
                c.scrollBy({ top: 100, behavior: 'smooth' });
            }
        });`,
	)
}
