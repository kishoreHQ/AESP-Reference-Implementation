package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func uiDistPath() string {
	if p := os.Getenv("AESP_UI_DIST"); p != "" {
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			return p
		}
	}
	// Walk from cwd: ui/dist
	candidates := []string{
		"ui/dist",
		filepath.Join("..", "ui", "dist"),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "..", "ui", "dist"))
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return ""
}

// spaFileServer serves static files with SPA fallback to index.html.
func spaFileServer(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never shadow API routes (should not reach here for /api)
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}
		path := filepath.Join(root, filepath.Clean("/"+r.URL.Path))
		if st, err := os.Stat(path); err == nil && !st.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		// SPA fallback
		http.ServeFile(w, r, filepath.Join(root, "index.html"))
	})
}
