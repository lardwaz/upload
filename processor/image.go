package processor

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
	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/job"
	"go.lsl.digital/lardwaz/upload/option"
	"go.lsl.digital/lardwaz/upload/processor/box"
	"go.lsl.digital/lardwaz/upload/processor/position"
	utypes "go.lsl.digital/lardwaz/upload/types"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
}

// Image implements the processor interface
type Image struct {
	options sdk.OptionsImage
}

// NewImage returns a new Image
func NewImage(opts ...func(sdk.OptionsImage)) *Image {
	options := option.EvaluateImageOptions(opts...)
	processor := &Image{
		options: options,
	}

	return processor
}

// Options returns OptionsImage
func (p Image) Options() sdk.OptionsImage {
	return p.options
}

// Process adds a job to process an image based on specific options
func (p *Image) Process(file sdk.Uploaded, validate bool) (sdk.Job, error) {
	content := file.Content()
	if !utypes.IsValidImage(content) {
		return nil, fmt.Errorf("image type invalid")
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		log.Printf("error decoding image: %v", err)
		return nil, err
	}

	// Check min width and height
	if validate && p.Options().MinWidth() != option.NoLimit && config.Width < p.Options().MinWidth() {
		log.Printf("image %v lower than min width: %v\n", file.DiskPath(), p.Options().MinWidth())
		return nil, fmt.Errorf("image width less than %dpx", p.Options().MinWidth())
	}

	if validate && p.Options().MinHeight() != option.NoLimit && config.Height < p.Options().MinHeight() {
		log.Printf("image %v lower than min height: %v\n", file.DiskPath(), p.Options().MinHeight())
		return nil, fmt.Errorf("image height less than %dpx", p.Options().MinHeight())
	}

	job := job.NewGeneric(file)

	go p.process(job, &config)

	return job, nil
}

func (p *Image) process(job sdk.Job, config *image.Config) {
	var (
		img image.Image
		err error

		isPROD = p.Options().IsPROD()
	)

	p.Options().Formats().Each(func(name string, format sdk.OptionsFormat) {
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
		if format.Backdrop() != nil && format.Backdrop().Path() != "" && !landscape {
			diskPathBackdrop := format.Backdrop().Path()
			// Scale down srcImage to fit the bounding box
			img = imaging.Fit(img, newWidth, newHeight, imaging.Lanczos)

			// Open a new image to use as backdrop layer
			var back image.Image
			if !isPROD {
				back, err = imaging.Open(diskPathBackdrop + "-" + format.Name())
			} else {
				var staticAsset *os.File
				staticAsset, err = box.Asset.Open(diskPathBackdrop + "-" + format.Name())
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

		if format.Watermark() != nil && format.Watermark().Path() != "" {
			diskPathWatermark := format.Watermark().Path()
			var watermark image.Image
			if !isPROD {
				watermark, err = imaging.Open(diskPathWatermark + "-" + format.Name())
			} else {
				var staticAsset *os.File
				staticAsset, err = box.Asset.Open(diskPathWatermark + "-" + format.Name())
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
					format.Watermark().SetHorizontal(position.Left)
					fallthrough
				case position.Left:
					watermarkPos.X += format.Watermark().OffsetX()
				case position.Right:
					RightX := bgBounds.Min.X + bgW - watermarkW
					watermarkPos.X = RightX - format.Watermark().OffsetX()
				case position.Center:
					CenterX := bgBounds.Min.X + bgW/2
					watermarkPos.X = CenterX - watermarkW/2 + format.Watermark().OffsetX()
				}

				switch format.Watermark().Vertical() {
				default:
					format.Watermark().SetVertical(position.Top)
					fallthrough
				case position.Top:
					watermarkPos.Y += format.Watermark().OffsetY()
				case position.Bottom:
					BottomY := bgBounds.Min.Y + bgH - watermarkH
					watermarkPos.Y = BottomY - format.Watermark().OffsetY()
				case position.Center:
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
