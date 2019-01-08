package upload

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/lsldigital/gocipe-upload/imagist"
	"github.com/lsldigital/gocipe-upload/util"

	"github.com/gosimple/slug"
)

// Constants for files package
const (
	_fileMode = os.FileMode(0644)
)

// Options holds file upload options
type Options struct {
	Dir            string
	Destination    string
	MediaPrefixURL string
	FileType       int
	MaxSize        int
	ConvertTo      string
	ImgDimensions  *imagist.ImageDimensions
}

// supported file types by fileupload
const (
	TypeInvalid = iota
	TypeImage
	TypeVideo
	TypeAudio
	TypeDocument
	TypeSheet
	TypeCSV
	TypePDF
)

var (
	_imagist *imagist.Imagist

	_env = util.EnvironmentDEV
)

func init() {
	_imagist = imagist.New()
}

// SetEnv sets the environment gocipe-upload operates in
func SetEnv(env string) {
	switch env {
	case util.EnvironmentDEV, util.EnvironmentPROD:
		// We are good :)
	default:
		// Invalid environment
		return
	}
	_env = env
}

// Upload validates and saves create file
func Upload(fileName string, fileContent []byte, options *Options) (string, string, error) {
	dirPath := path.Join(options.Dir, options.Destination)
	fileName = buildFileName(fileName)
	filePath := path.Join(options.MediaPrefixURL, options.Destination, fileName)
	fileDiskPath := filepath.Join(dirPath, fileName)

	// Create full directory structure to store image
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Printf("error creating directories: %v\n", err)
		return fileDiskPath, filePath, err
	}

	// Create file on disk if it does not exist
	if _, err := os.Stat(fileDiskPath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(fileDiskPath, fileContent, _fileMode); err != nil {
			log.Printf("error writing %v: %v\n", fileDiskPath, err)
			return fileDiskPath, filePath, err
		}
	}

	file, err := os.Open(fileDiskPath)
	if err != nil {
		log.Printf("error opening %v: %v\n", fileDiskPath, err)
		return fileDiskPath, filePath, err
	}
	defer file.Close()

	buf, err := ioutil.ReadFile(fileDiskPath)
	if err != nil {
		log.Printf("error reading %v: %v\n", fileDiskPath, err)
		return fileDiskPath, filePath, err
	}

	fileSize := len(buf)
	if options.MaxSize != util.NoLimit && fileSize > options.MaxSize {
		log.Printf("file %v greater than max file size: %v\n", fileName, options.MaxSize)
		return fileDiskPath, filePath, fmt.Errorf("file max size error")
	}

	if options.ConvertTo != "" {
		fileDiskPath, filePath, err = changeExt(fileDiskPath, filePath, options.ConvertTo)
		if err != nil {
			return fileDiskPath, filePath, err
		}
	}

	switch options.FileType {
	case TypeImage:
		err := _imagist.Add(buf, fileDiskPath, options.ImgDimensions, true, _env)
		return fileDiskPath, filePath, err
	case TypeVideo:
		// TODO: Not yet implemented
	case TypeAudio:
		// TODO: Not yet implemented
	case TypeDocument:
		// TODO: Not yet implemented
	case TypeSheet:
		// TODO: Not yet implemented
	case TypeCSV:
		// TODO: Not yet implemented
	case TypePDF:
		// TODO: Not yet implemented
	default:
		// Invalid file type in config
		// Do nothing
	}

	return fileDiskPath, filePath, nil
}

// Delete deletes one file
func Delete(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func buildFileName(oldFilename string) string {
	oldExt := filepath.Ext(oldFilename)
	newFilename := strings.TrimSuffix(oldFilename, oldExt)
	return slug.Make(newFilename) + "_" + time.Now().Format("20060102150405") + oldExt
}

func changeExt(fileDiskPath, filepath string, newExt string) (string, string, error) {
	oldExt := path.Ext(fileDiskPath)
	newExt = "." + newExt
	newfileDiskPath := strings.TrimSuffix(fileDiskPath, oldExt) + newExt
	newfilePath := strings.TrimSuffix(filepath, oldExt) + newExt

	if err := os.Rename(fileDiskPath, newfileDiskPath); err != nil {
		return fileDiskPath, filepath, fmt.Errorf("image ext change to %v failed", newExt)
	}

	return newfileDiskPath, newfilePath, nil
}

// httpImageDirHandler is an http.Handler that serves a directory.
// If a generated file is missing, it yields a temporary redirect to the original file.
type httpImageDirHandler struct {
	root   http.FileSystem
	prefix string
	opts   *Options
}

func (s httpImageDirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	var suffix string
	for _, format := range s.opts.ImgDimensions.Formats {
		formatSuffix := ":" + format.Name
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

	go func() {
		var buf []byte

		dirPath := path.Join(s.opts.Dir, s.opts.Destination)
		fileName := strings.TrimPrefix(noSuffix, "/"+s.opts.Destination)
		fileDiskPath := filepath.Join(dirPath, fileName)

		buf, err := ioutil.ReadFile(fileDiskPath)
		if err != nil {
			log.Printf("error opening %v: %v\n", fileDiskPath, err)
		}

		_imagist.Add(buf, fileDiskPath, s.opts.ImgDimensions, false, _env)

	}()
	http.Redirect(w, r, p, http.StatusTemporaryRedirect)
}

//HTTPImageDirHandler serves images from a directory with imagist fallback
// func HTTPImageDirHandler(router *mux.Router, root http.FileSystem, prefix string, paths map[string]*Options) {
// 	for path, opts := range paths {
// 		h := web.FileServerWithNotFoundHandler(root, httpImageDirHandler{root: root, prefix: prefix, opts: opts})
// 		router.PathPrefix(prefix + "/" + path).Handler(http.StripPrefix(prefix, h))
// 	}
// }
