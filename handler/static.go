package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	staticPath = "./dist/"
)

// HandleStaticFile handles requests to the front end. It checks to see if the path leads to an actual file and, if so, serves it. Otherwise it serves the index page, assuming it's a route in the front end. This was copied with minimal changes from the mux documentation.
func (h *Handler) HandleStaticFile(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(staticPath, "/"))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		log.Printf("error serving static file '%s': %s", r.URL.Path, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
}
