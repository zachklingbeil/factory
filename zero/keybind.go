package zero

import (
	"fmt"
	"html/template"
	"strings"
)

type Keybind interface {
	AddKeybind(containerId string, keyHandlers map[string]string) One
}
type keybind struct{}

func (k *keybind) AddKeybind(containerId string, keyHandlers map[string]string) One {
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
