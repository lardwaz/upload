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
)

// Job represents current image file being processed
type Job struct {
	File	Uploaded
	Config	*image.Config
	Done 	chan struct{}
}

type assetBoxer interface {
	Open(string) (*os.File, error)
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
}

// BackdropImage sets the disk path for backdrop images
func BackdropImage(path string) {
	_diskPathBackdrop = path
}

// WatermarkImage sets the disk path for watermark images
func WatermarkImage(path string) {
	_diskPathWatermark = path
}

// AssetBox sets the asset box to retrieve static assets
func AssetBox(assetBox assetBoxer) {
	_assetBox = assetBox
}

// ImageProcessor implements the processor interface
type ImageProcessor struct{
	options *optionsImage
}

// NewImageProcessor returns a new ImageProcessor
func NewImageProcessor(opts ...OptionImage) *ImageProcessor {
	options := EvaluateImageOptions(opts...)
	processor := &ImageProcessor{
		options: options,
	}

	return processor
}

// Process adds a job to process an image based on specific options
func (p *ImageProcessor) Process(file Uploaded, validate bool) (*Job, error) {
	content := file.Content()
	if !filetype.IsImage(content) {
		return nil, fmt.Errorf("image type invalid")
	}

	config, imgType, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		log.Printf("error decoding image: %v", err)
		return nil, err
	}

	switch imgType {
	case typeImageJPG, typeImageJPEG, typeImagePNG:
		//all ok
	default:
		return nil, fmt.Errorf("image type %s invalid", imgType)
	}

	// Check min width and height
	if validate && p.options.minWidth != core.NoLimit && config.Width < p.options.minWidth {
		log.Printf("image %v lower than min width: %v\n", file.DiskPath(), p.options.minWidth)
		return nil, fmt.Errorf("image width less than %dpx", p.options.minWidth)
	}

	if validate && p.options.minHeight != core.NoLimit && config.Height < p.options.minHeight {
		log.Printf("image %v lower than min height: %v\n", file.DiskPath(), p.options.minHeight)
		return nil, fmt.Errorf("image height less than %dpx", p.options.minHeight)
	}

	job := &Job{
		File:	file,
		Config:	&config,
		Done: 	make(chan struct{}),
	}
	
	go p.process(job)

	return job, nil
}

func (p *ImageProcessor) process(job *Job) {
	var (
		img image.Image
		err error
	)

	for _, format := range p.options.formats {
		if format.name == "" || (format.width <= 0 && format.height <= 0) {
			continue
		}

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
		if format.width > job.Config.Width {
			newWidth = job.Config.Width
		}
		if format.height > job.Config.Height {
			newHeight = job.Config.Height
		}

		// -1 pixel size does not exist
		if format.width < 0 {
			newWidth = 0
		}
		if format.height < 0 {
			newHeight = 0
		}

		landscape := job.Config.Height < job.Config.Width
		preserveAspect := newWidth <= 0 || newHeight <= 0

		// Do not crop and resize when using backdrop but downscale
		if _diskPathBackdrop != "" && format.backdrop && !landscape {
			// Scale down srcImage to fit the bounding box
			img = imaging.Fit(img, newWidth, newHeight, imaging.Lanczos)

			// Open a new image to use as backdrop layer
			var back image.Image
			if core.Env == core.EnvironmentDEV {
				back, err = imaging.Open(_diskPathBackdrop)
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
		} else if preserveAspect {
			// Resize srcImage to proper width or height preserving the aspect ratio.
			img = imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
		} else {
			// Resize and crop the image to fill the [newWidth x newHeight] area
			img = imaging.Fill(img, newWidth, newHeight, imaging.Center, imaging.Lanczos)
		}

		if _diskPathWatermark != "" && format.watermark != nil {
			var watermark image.Image
			if core.Env == core.EnvironmentDEV {
				watermark, err = imaging.Open(_diskPathWatermark + ":" + format.name)
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

	job.Done <- struct{}{}
}