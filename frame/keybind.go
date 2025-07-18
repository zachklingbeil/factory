package frame

import "html/template"

func (f *Frame) AddPathKeybind(key, path string) *template.HTML {
	return f.AddJS(
		`document.addEventListener('keydown', function(event) {
            if (event.key === '` + key + `') {
                fetch('` + path + `')
                    .then(r => r.text())
                    .then(html => {
                        const c = document.getElementById('frame');
                        if (c) c.innerHTML = html;
                    });
            }
        });`,
	)
}

func (f *Frame) AddScrollKeybinds() *template.HTML {
	return f.AddJS(
		`document.addEventListener('keydown', function(event) {
            const c = document.getElementById('frame');
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

func (f *Frame) AddNavKeybinds() *template.HTML {
	prev := f.AddPathKeybind("q", "/frame/prev")
	next := f.AddPathKeybind("e", "/frame/next")
	html := *prev + *next
	return &html
}
