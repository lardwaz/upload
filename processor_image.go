package upload

import (
	"bytes"
	"fmt"
	"log"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
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

	chanSize = 10
)

// Anchor points for X,Y
const (
	Left = iota
	Right
	Top
	Bottom
	Center
)

var (
	// Disk paths to static assets
	_diskPathWatermark string
	_diskPathBackdrop  string

	// _assetBox satisfies the AssetBoxer interface
	_assetBox assetBoxer

	// TopLeft is the top-left position for watermark
	TopLeft = EvaluateWatermarkOptions(WatermarkHorizontal(Left), WatermarkVertical(Top))
	// TopCenter is the top-center position for watermark
	TopCenter = EvaluateWatermarkOptions(WatermarkHorizontal(Center), WatermarkVertical(Top))
	// TopRight is the top-right position for watermark
	TopRight = EvaluateWatermarkOptions(WatermarkHorizontal(Right), WatermarkVertical(Top))
	// CenterRight is the center-right position for watermark
	CenterRight = EvaluateWatermarkOptions(WatermarkHorizontal(Right), WatermarkVertical(Center))
	// BottomRight is the bottom-right position for watermark
	BottomRight = EvaluateWatermarkOptions(WatermarkHorizontal(Right), WatermarkVertical(Bottom))
	// BottomCenter is the bottom-center position for watermark
	BottomCenter = EvaluateWatermarkOptions(WatermarkHorizontal(Center), WatermarkVertical(Bottom))
	// BottomLeft is the bottom-left position for watermark
	BottomLeft = EvaluateWatermarkOptions(WatermarkHorizontal(Left), WatermarkVertical(Bottom))
	// CenterLeft is the center-left position for watermark
	CenterLeft = EvaluateWatermarkOptions(WatermarkHorizontal(Left), WatermarkVertical(Center))
)
type Job struct {
	File UploadedFile
	Config       *image.Config
}

type assetBoxer interface {
	Open(string) (*os.File, error)
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

func (p *ImageProcessor) process(job Job) {
	var (
		img image.Image
		err error
	)

	for _, format := range p.options.formats {
		if format.name == "" || format.width <= 0 || format.height <= 0 {
			continue
		}

		// imageProcess(job, format)

		imgDiskPath := job.File.DiskPath()

		img, err = imaging.Open(imgDiskPath)
		if err != nil {
			log.Printf("Image error: %v\n", err)
			continue
		}

		// Prepare metra for processing
		newWidth := format.width
		newHeight := format.height

		// Do not upscale
		if job.Config.Width < format.width {
			newWidth = job.Config.Width
		}
		if job.Config.Height < job.Config.Height {
			newHeight = job.Config.Height
		}

		landscape := job.Config.Height < job.Config.Width

		// Do not crop and resize when using backdrop but downscale
		if _diskPathBackdrop != "" && format.backdrop && !landscape {
			// Scale down srcImage to fit the bounding box
			img = imaging.Fit(img, newWidth, newHeight, imaging.Lanczos)

			// Open a new image to use as backdrop layer
			var back image.Image
			if core.Env == core.EnvironmentDEV {
				back, err = imaging.Open("../assets/" + _diskPathBackdrop)
			} else {
				var staticAsset *os.File
				staticAsset, err = _assetBox.Open(_diskPathBackdrop)
				if err != nil {
					// if err, fall back to a blue background backdrop
					back = imaging.New(format.width, format.height, color.NRGBA{0, 29, 56, 0})
				}
				defer staticAsset.Close()
				back, _, err = image.Decode(staticAsset)
			}

			if err != nil {
				// if err, fall back to a blue background backdrop
				back = imaging.New(format.width, format.height, color.NRGBA{0, 29, 56, 0})
			} else {
				// Resize and crop backdrop accordingly
				back = imaging.Fill(back, format.width, format.height, imaging.Center, imaging.Lanczos)
			}

			// Overlay image in center on backdrop layer
			img = imaging.OverlayCenter(back, img, 1.0)
		} else if _diskPathBackdrop != "" {
			// Resize and crop the image to fill the [newWidth x newHeight] area
			img = imaging.Fill(img, newWidth, newHeight, imaging.Center, imaging.Lanczos)
		}

		if _diskPathWatermark != "" && format.watermark != nil {
			var watermark image.Image
			if core.Env == core.EnvironmentDEV {
				watermark, err = imaging.Open("../assets/" + _diskPathWatermark + ":" + format.name)
			} else {
				var staticAsset *os.File
				staticAsset, err = _assetBox.Open(_diskPathWatermark + ":" + format.name)
				if err != nil {
					log.Printf("Watermark not found: %v", err)
					continue
				}
				defer staticAsset.Close()
				watermark, _, err = image.Decode(staticAsset)
			}
			if err == nil {
				bgBounds := img.Bounds()
				bgW := bgBounds.Dx()
				bgH := bgBounds.Dy()

				watermarkBounds := watermark.Bounds()
				watermarkW := watermarkBounds.Dx()
				watermarkH := watermarkBounds.Dy()

				var watermarkPos image.Point

				switch format.watermark.horizontal {
				default:
					format.watermark.horizontal = Left
					fallthrough
				case Left:
					watermarkPos.X += format.watermark.offsetX
				case Right:
					RightX := bgBounds.Min.X + bgW - watermarkW
					watermarkPos.X = RightX - format.watermark.offsetX
				case Center:
					CenterX := bgBounds.Min.X + bgW/2
					watermarkPos.X = CenterX - watermarkW/2 + format.watermark.offsetX
				}

				switch format.watermark.vertical {
				default:
					format.watermark.vertical = Top
					fallthrough
				case Top:
					watermarkPos.Y += format.watermark.offsetY
				case Bottom:
					BottomY := bgBounds.Min.Y + bgH - watermarkH
					watermarkPos.Y = BottomY - format.watermark.offsetY
				case Center:
					CenterY := bgBounds.Min.Y + bgH/2
					watermarkPos.Y = CenterY - watermarkH/2 + format.watermark.offsetY
				}

				img = imaging.Overlay(img, watermark, watermarkPos, 1.0)
			}
		}

		imagingFormat, err := imaging.FormatFromFilename(imgDiskPath)
		if err != nil {
			log.Printf("Image get format error: %v", err)
			continue
		}

		outputFile, err := os.Create(imgDiskPath + ":" + format.name)
		if err != nil {
			log.Printf("Image get format error: %v", err)
			continue
		}
		defer outputFile.Close()

		if err := imaging.Encode(outputFile, img, imagingFormat); err != nil {
			log.Printf("Image encode format error: %v", err)
		}
	}

	p.done <- job.File.DiskPath()
}