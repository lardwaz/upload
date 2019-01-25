package upload

import (
	"bytes"
	"fmt"
	"log"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/lsldigital/gocipe-upload/core"
	filetype "gopkg.in/h2non/filetype.v1"
)

const (
	// typeImageJPG denotes image of file type jpg
	typeImageJPG = "jpg"
	// typeImageJPEG denotes image of file type jpeg
	typeImageJPEG = "jpeg"
	// typeImagePNG denotes image of file type png
	typeImagePNG = "png"
)

const chanSize = 10

type Job struct {
	File UploadedFile
	Config       *image.Config
}

type ImageProcessor struct{
	options *optionsImage
	jobs chan Job
	done chan string
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
}

func NewImageProcessor(opts ...OptionImage) *ImageProcessor {
	options := EvaluateImageOptions(opts...)
	processor := &ImageProcessor{
		options: options,
		jobs: make(chan Job, chanSize),
		done: make(chan string, chanSize),
	}

	go processor.startDispatcher()

	return processor
}

// Process adds a job to process an image based on specific options
func (p *ImageProcessor) Process(file UploadedFile, validate bool) error {
	content := file.Content()
	if !filetype.IsImage(content) {
		return fmt.Errorf("image type invalid")
	}

	config, imgType, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		log.Printf("error decoding image: %v", err)
		return err
	}

	switch imgType {
	case typeImageJPG, typeImageJPEG, typeImagePNG:
		//all ok
	default:
		return fmt.Errorf("image type %s invalid", imgType)
	}

	// Check min width and height
	if validate && p.options.minWidth != core.NoLimit && config.Width < p.options.minWidth {
		log.Printf("image %v lower than min width: %v\n", file.DiskPath(), p.options.minWidth)
		return fmt.Errorf("image width less than %dpx", p.options.minWidth)
	}

	if validate && p.options.minHeight != core.NoLimit && config.Height < p.options.minHeight {
		log.Printf("image %v lower than min height: %v\n", file.DiskPath(), p.options.minHeight)
		return fmt.Errorf("image height less than %dpx", p.options.minHeight)
	}

	job := Job{
		File: file,
		Config:       &config,
	}
	p.jobs <- job

	return nil
}


func (p *ImageProcessor) startDispatcher() {
	jobs := make(map[string]struct{})

	for {
		select {
		case done := <-p.done:
			delete(jobs, done)
		case job := <-p.jobs:
			if _, exists := jobs[job.File.DiskPath()]; !exists {
				jobs[job.File.DiskPath()] = struct{}{}
				go p.process(job)
			}
		}
	}
}

func (p *ImageProcessor) process(job Job) error {
	panic("Not implemented")
	return nil
}