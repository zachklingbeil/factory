package universe

type Universe struct {
	Element *Element
	Index   string
}

func New() *Universe {
	return &Universe{
		Element: NewElements(),
	}
}

// // ServeIndex serves the embedded index.html file.
// func (u *Universe) ServeIndex(w http.ResponseWriter, r *http.Request) {
// 	data, err := content.ReadFile("index.html")
// 	if err != nil {
// 		http.Error(w, "index.html not found", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.Write(data)
// }

// // ServeMainFragment serves HTML content for the <main> element.
// func (u *Universe) ServeMainFragment(w http.ResponseWriter, r *http.Request) {
// 	fragment := `<h2>Welcome to the Home Page!</h2><p>This is server-rendered content.</p>`
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.Write([]byte(fragment))
// }
