package apps

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func registerPlayerUI() {
	rootDir := playerWebDir()

	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/player/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		relativePath := strings.TrimPrefix(r.URL.Path, "/player/")
		relativePath = filepath.Clean(strings.TrimPrefix(relativePath, "/"))

		switch {
		case relativePath == "." || relativePath == "":
			http.ServeFile(w, r, filepath.Join(rootDir, "index.html"))
		case strings.HasPrefix(relativePath, ".."):
			http.NotFound(w, r)
		default:
			targetPath := filepath.Join(rootDir, relativePath)
			if fileExists(targetPath) {
				http.ServeFile(w, r, targetPath)
				return
			}

			// Fallback to SPA entry for client-side routes.
			http.ServeFile(w, r, filepath.Join(rootDir, "index.html"))
		}
	})
}

func playerWebDir() string {
	if configuredDir := strings.TrimSpace(os.Getenv("VMUSIC_WEB_DIR")); configuredDir != "" {
		return configuredDir
	}

	workDir, err := os.Getwd()
	if err != nil {
		return filepath.Join("web", "player")
	}

	sourceDir := filepath.Join(workDir, "web", "player")
	distDir := filepath.Join(sourceDir, "dist")
	if fileExists(filepath.Join(distDir, "index.html")) {
		return distDir
	}
	return sourceDir
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
