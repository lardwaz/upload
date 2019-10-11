package handler

import (
	"net/http"
	"path"
	"strings"

	"go.lsl.digital/lardwaz/upload"
)

// HTTPImageDir is an http.Handler that serves a directory.
// If a generated file is missing, it yields a temporary redirect to the original file.
type HTTPImageDir struct {
	root   http.FileSystem
	prefix string
	opts   upload.OptionsImage
}

func (h HTTPImageDir) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	var suffix string

	formats := h.opts.Formats()

	formats.Each(func(name string, format upload.OptionsFormat) {
		formatSuffix := "-" + format.Name()
		if strings.HasSuffix(p, formatSuffix) {
			suffix = formatSuffix
		}
	})

	if suffix == "" {
		//a previous attempt to lookup the file resulted into a call to this function
		//do not attempt to look up again
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	noSuffix := strings.TrimSuffix(p, suffix)
	p = path.Join(h.prefix, noSuffix)

	http.Redirect(w, r, p, http.StatusTemporaryRedirect)
}
