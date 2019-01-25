package upload

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/lsldigital/gocipe-upload/core"
)

// Uploaded represents the uploaded file
type Uploaded interface {
	URLPath() string
	DiskPath() string
	Content() []byte
	Save(content []byte, overwrite bool) error
	Delete() error
	ChangeExt(newExt string) error
}

// UploadedFile implements File interface
type UploadedFile struct {
	url      string
	diskPath string
	content  []byte
	options  options
}

// NewUploadedFile returns a new UploadedFile struct
func NewUploadedFile(name string, opts options) UploadedFile {
	dirPath := path.Join(opts.dir, opts.destination)
	name = AddTimestamp(name)
	urlPath := path.Join(opts.mediaPrefixURL, opts.destination, name)
	diskPath := filepath.Join(dirPath, name)

	return UploadedFile{
		url:      urlPath,
		diskPath: diskPath,
		options:  opts,
	}
}

// URLPath returns the url path of file
func (u *UploadedFile) URLPath() string {
	return u.url
}

// DiskPath returns the path of file on disk
func (u *UploadedFile) DiskPath() string {
	return u.diskPath
}

// Content returns the path of file on disk
func (u *UploadedFile) Content() []byte {
	return u.content
}

// Save saves file on disk if it does not exist
func (u *UploadedFile) Save(content []byte, overwrite bool) error {
	if !overwrite {
		return nil
	}

	// Verify size
	size := len(content)
	if u.options.maxSize != core.NoLimit && size > u.options.maxSize {
		log.Printf("file %v greater than max file size: %v\n", u.diskPath, u.options.maxSize)
		return fmt.Errorf("file max size error")
	}

	// Creates full directory structure to store image
	dir := path.Dir(u.diskPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("error creating directories %v : %v\n", dir, err)
		return err
	}

	if err := ioutil.WriteFile(u.DiskPath(), content, os.FileMode(0644)); err != nil {
		log.Printf("error writing %v: %v\n", u.DiskPath(), err)
		return err
	}

	u.content = content

	return nil
}

// Delete deletes one file on disk
func (u *UploadedFile) Delete() error {
	if err := os.Remove(u.DiskPath()); err != nil {
		return err
	}
	return nil
}

// ChangeExt changes the extension of file on disk
func (u *UploadedFile) ChangeExt(newExt string) error {
	if newExt == "" {
		return nil
	}

	oldExt := path.Ext(u.DiskPath())
	newFileDiskPath := strings.TrimSuffix(u.DiskPath(), oldExt) + "." + newExt
	newFileURLPath := strings.TrimSuffix(u.URLPath(), oldExt) + "." + newExt

	if err := os.Rename(u.DiskPath(), newFileDiskPath); err != nil {
		return fmt.Errorf("image ext change to %v failed", newExt)
	}

	// if everything ok, update paths
	u.diskPath = newFileDiskPath
	u.url = newFileURLPath

	return nil
}

// AddTimestamp add timestamp information to a filename
func AddTimestamp(oldFilename string) string {
	oldExt := filepath.Ext(oldFilename)
	newFilename := strings.TrimSuffix(oldFilename, oldExt)
	return slug.Make(newFilename) + "_" + time.Now().Format("20060102150405") + oldExt
}
