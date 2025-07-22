package one

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	*zero.Zero
	*fx.Fx
	Frames []*zero.One
}

func NewOne() *One {
	one := &One{
		Zero: zero.NewZero(),
		Fx:   fx.InitFx(),
	}

	return one
}

func (o *One) Pathless(cssPath, domain string, frames int) {
	pathless := o.BuildPathless(cssPath, domain, frames)
	o.Frames = append(o.Frames, &pathless)
}

func (o *One) AddPath(dir, prefix string) error {
	files, err := o.SourcePath(dir)
	if err != nil {
		return err
	}

	for name, val := range files {
		routePath := "/" + strings.Trim(prefix, "/") + "/" + name
		fileVal := val
		o.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", fileVal.Type)
			w.Write(fileVal.Data)
		})
	}
	return nil
}

// frame[0] is served at the root path ("/"), frame[1] and onwards are served at "/frame/{index}"
func (o *One) Index() {
	// Serve index 0 at "/"
	o.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-FRAMES", strconv.Itoa(len(o.Frames)))
		if len(o.Frames) == 0 {
			http.Error(w, "No frame available", http.StatusNotFound)
			return
		}
		w.Write([]byte(*o.Frames[0]))
	})

	// Serve other frames at /frame/{index}
	o.HandleFunc("/frame/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-FRAMES", strconv.Itoa(len(o.Frames)))
		idxStr := strings.TrimPrefix(r.URL.Path, "/frame/")
		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx < 1 || idx >= len(o.Frames) {
			http.Error(w, "Frame not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(*o.Frames[idx]))
	})
}
