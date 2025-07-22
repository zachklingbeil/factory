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
func (o *One) ServeRoot() {
	o.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(*o.Frames[0]))
	})
}

func (o *One) Pathless(cssPath, domain string, frames int) {
	frame := o.BuildPathless(cssPath, domain, frames)
	o.Frames = append(o.Frames, &frame)
}

// Serve frames at /frame/{index}
func (o *One) ServeFrames() {
	o.HandleFunc("/frame/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-FRAMES", strconv.Itoa(len(o.Frames)))
		idxStr := strings.TrimPrefix(r.URL.Path, "/frame/")
		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx < 0 || idx >= len(o.Frames) {
			http.Error(w, "Frame not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(*o.Frames[idx]))
	})
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
