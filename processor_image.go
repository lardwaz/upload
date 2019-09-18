package upload

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"go.lsl.digital/lardwaz/upload/core"
)

const (
	// TypeImageJPG denotes image of file type jpg
	TypeImageJPG = "jpg"
	// TypeImageJPEG denotes image of file type jpeg
	TypeImageJPEG = "jpeg"
	// TypeImagePNG denotes image of file type png
	TypeImagePNG = "png"
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
type ImageProcessor struct {
	options OptionsImage
}

// NewImageProcessor returns a new ImageProcessor
func NewImageProcessor(opts ...func(OptionsImage)) ImageProcessor {
	options := evaluateImageOptions(opts...)
	processor := ImageProcessor{
		options: options,
	}

	return processor
}

// Options returns OptionsImage
func (p ImageProcessor) Options() OptionsImage {
	return p.options
}

// Process adds a job to process an image based on specific options
func (p *ImageProcessor) Process(file Uploaded, validate bool) (Job, error) {
	content := file.Content()
	if !isValidImage(content) {
		return nil, fmt.Errorf("image type invalid")
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		log.Printf("error decoding image: %v", err)
		return nil, err
	}

	// Check min width and height
	if validate && p.options.MinWidth() != core.NoLimit && config.Width < p.options.MinWidth() {
		log.Printf("image %v lower than min width: %v\n", file.DiskPath(), p.options.MinWidth())
		return nil, fmt.Errorf("image width less than %dpx", p.options.MinWidth())
	}

	if validate && p.options.MinHeight() != core.NoLimit && config.Height < p.options.MinHeight() {
		log.Printf("image %v lower than min height: %v\n", file.DiskPath(), p.options.MinHeight())
		return nil, fmt.Errorf("image height less than %dpx", p.options.MinHeight())
	}

	job := NewGenericJob(file)

	go p.process(job, &config)

	return job, nil
}

func (p *ImageProcessor) process(job Job, config *image.Config) {
	var (
		img image.Image
		err error
	)

	formats := p.options.Formats()

	formats.Each(func(name string, format OptionsFormat) {
		if format.Name() == "" {
			return
		}

		imgDiskPath := job.File().DiskPath()

		img, err = imaging.Open(imgDiskPath)
		if err != nil {
			log.Printf("Image error: %v\n", err)
			return
		}

		// Prepare metra for processing
		newWidth := format.Width()
		newHeight := format.Height()

		// Do not upscale
		if format.Width() > config.Width {
			newWidth = config.Width
		}
		if format.Height() > config.Height {
			newHeight = config.Height
		}

		// -1 pixel size does not exist
		if format.Width() < 0 {
			newWidth = 0
		}
		if format.Height() < 0 {
			newHeight = 0
		}

		landscape := config.Height < config.Width
		preserveAspect := newWidth <= 0 || newHeight <= 0

		// Do not crop and resize when using backdrop but downscale
		if _diskPathBackdrop != "" && format.Backdrop() && !landscape {
			// Scale down srcImage to fit the bounding box
			img = imaging.Fit(img, newWidth, newHeight, imaging.Lanczos)

			// Open a new image to use as backdrop layer
			var back image.Image
			if core.Env == core.EnvironmentDEV {
				back, err = imaging.Open(_diskPathBackdrop + "-" + format.Name())
			} else {
				var staticAsset *os.File
				staticAsset, err = _assetBox.Open(_diskPathBackdrop + "-" + format.Name())
				if err != nil {
					// if err, fall back to a blue background backdrop
					back = imaging.New(format.Width(), format.Height(), color.NRGBA{0, 29, 56, 0})
				}
				defer staticAsset.Close()
				back, _, err = image.Decode(staticAsset)
			}

			if err != nil {
				// if err, fall back to a blue background backdrop
				back = imaging.New(format.Width(), format.Height(), color.NRGBA{0, 29, 56, 0})
			} else {
				// Resize and crop backdrop accordingly
				back = imaging.Fill(back, format.Width(), format.Height(), imaging.Center, imaging.Lanczos)
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

		if _diskPathWatermark != "" && format.Watermark() != nil {
			var watermark image.Image
			if core.Env == core.EnvironmentDEV {
				watermark, err = imaging.Open(_diskPathWatermark + "-" + format.Name())
			} else {
				var staticAsset *os.File
				staticAsset, err = _assetBox.Open(_diskPathWatermark + "-" + format.Name())
				if err != nil {
					log.Printf("Watermark not found: %v", err)
					return
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

				switch format.Watermark().Horizontal() {
				default:
					format.Watermark().SetHorizontal(Left)
					fallthrough
				case Left:
					watermarkPos.X += format.Watermark().OffsetX()
				case Right:
					RightX := bgBounds.Min.X + bgW - watermarkW
					watermarkPos.X = RightX - format.Watermark().OffsetX()
				case Center:
					CenterX := bgBounds.Min.X + bgW/2
					watermarkPos.X = CenterX - watermarkW/2 + format.Watermark().OffsetX()
				}

				switch format.Watermark().Vertical() {
				default:
					format.Watermark().SetVertical(Top)
					fallthrough
				case Top:
					watermarkPos.Y += format.Watermark().OffsetY()
				case Bottom:
					BottomY := bgBounds.Min.Y + bgH - watermarkH
					watermarkPos.Y = BottomY - format.Watermark().OffsetY()
				case Center:
					CenterY := bgBounds.Min.Y + bgH/2
					watermarkPos.Y = CenterY - watermarkH/2 + format.Watermark().OffsetY()
				}

				img = imaging.Overlay(img, watermark, watermarkPos, 1.0)
			}
		}

		imagingFormat, err := imaging.FormatFromFilename(imgDiskPath)
		if err != nil {
			log.Printf("Image get format error: %v", err)
			return
		}

		outputFile, err := os.Create(imgDiskPath + "-" + format.Name())
		if err != nil {
			log.Printf("Image get format error: %v", err)
			return
		}
		defer outputFile.Close()

		if err := imaging.Encode(outputFile, img, imagingFormat); err != nil {
			log.Printf("Image encode format error: %v", err)
		}
	})

	job.SetDone()
}
