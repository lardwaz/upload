package upload

import (
	"net/http"
	"path"
	"strings"
)

// httpImageDirHandler is an http.Handler that serves a directory.
// If a generated file is missing, it yields a temporary redirect to the original file.
type httpImageDirHandler struct {
	root   http.FileSystem
	prefix string
	opts   *OptionsImage
}

func (s httpImageDirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	var suffix string
	for _, format := range s.opts.Formats() {
		formatSuffix := ":" + format.name
		if strings.HasSuffix(p, formatSuffix) {
			suffix = formatSuffix
		}
	}

	if suffix == "" {
		//a previous attempt to lookup the file resulted into a call to this function
		//do not attempt to look up again
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	noSuffix := strings.TrimSuffix(p, suffix)
	p = path.Join(s.prefix, noSuffix)

	// go func() {
	// 	var buf []byte

	// 	dirPath := path.Join(s.opts.dir, s.opts.destination)
	// 	fileName := strings.TrimPrefix(noSuffix, "/"+s.opts.destination)
	// 	fileDiskPath := filepath.Join(dirPath, fileName)

	// 	buf, err := ioutil.ReadFile(fileDiskPath)
	// 	if err != nil {
	// 		log.Printf("error opening %v: %v\n", fileDiskPath, err)
	// 	}

	// 	_imagist.Add(buf, fileDiskPath, s.opts.ImgDimensions, false)

	// }()
	http.Redirect(w, r, p, http.StatusTemporaryRedirect)
}

//HTTPImageDirHandler serves images from a directory with imagist fallback
// func HTTPImageDirHandler(router *mux.Router, root http.FileSystem, prefix string, paths map[string]*Options) {
// 	for path, opts := range paths {
// 		h := web.FileServerWithNotFoundHandler(root, httpImageDirHandler{root: root, prefix: prefix, opts: opts})
// 		router.PathPrefix(prefix + "/" + path).Handler(http.StripPrefix(prefix, h))
// 	}
// }
