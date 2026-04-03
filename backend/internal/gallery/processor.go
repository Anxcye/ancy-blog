// File: processor.go
// Purpose: Extract EXIF metadata, generate display/placeholder image assets, and strip GPS data.
// Module: backend/internal/gallery, image processing layer.
// Related: gallery handler upload flow, storage uploader, and gallery service.
package gallery

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"strings"
	"time"

	"github.com/buckket/go-blurhash"
	"github.com/disintegration/imaging"
	"github.com/jdeng/goheif"
	"github.com/rwcarlsen/goexif/exif"
)

// ExifData holds whitelisted EXIF fields extracted from an image.
type ExifData struct {
	CameraMake   string
	CameraModel  string
	LensModel    string
	FocalLength  string
	Aperture     string
	ShutterSpeed string
	ISO          string
	TakenAt      time.Time
	Latitude     float64
	Longitude    float64
	HasGPS       bool
	Width        int
	Height       int
}

// ProcessedAssets holds the generated image derivatives.
type ProcessedAssets struct {
	Display         []byte
	DisplayWidth    int
	DisplayHeight   int
	Large           []byte
	LargeWidth      int
	LargeHeight     int
	PlaceholderHash string
}

const (
	displayMaxWidth = 800
	largeMaxWidth   = 2400
	jpegQuality     = 85
	blurHashXComp   = 4
	blurHashYComp   = 3
)

// ExtractExif reads EXIF data from the image bytes.
func ExtractExif(data []byte) ExifData {
	result := ExifData{}

	x, err := decodeExifMetadata(data)
	if err != nil {
		// No EXIF data, try to get dimensions from image decoding.
		if img, err2 := decodeImageConfig(data); err2 == nil {
			result.Width = img.Width
			result.Height = img.Height
		}
		return result
	}

	// Camera make
	if tag, err := x.Get(exif.Make); err == nil {
		result.CameraMake = cleanExifString(tag.String())
	}
	// Camera model
	if tag, err := x.Get(exif.Model); err == nil {
		result.CameraModel = cleanExifString(tag.String())
	}
	// Lens model
	if tag, err := x.Get(exif.LensModel); err == nil {
		result.LensModel = cleanExifString(tag.String())
	}
	// Focal length
	if tag, err := x.Get(exif.FocalLength); err == nil {
		if num, denom, err2 := tag.Rat2(0); err2 == nil && denom != 0 {
			fl := float64(num) / float64(denom)
			result.FocalLength = fmt.Sprintf("%.0fmm", fl)
		}
	}
	// Aperture (FNumber)
	if tag, err := x.Get(exif.FNumber); err == nil {
		if num, denom, err2 := tag.Rat2(0); err2 == nil && denom != 0 {
			f := float64(num) / float64(denom)
			result.Aperture = fmt.Sprintf("f/%.1f", f)
		}
	}
	// Shutter speed (ExposureTime)
	if tag, err := x.Get(exif.ExposureTime); err == nil {
		if num, denom, err2 := tag.Rat2(0); err2 == nil && denom != 0 {
			if num < denom {
				result.ShutterSpeed = fmt.Sprintf("%d/%d", num, denom)
			} else {
				result.ShutterSpeed = fmt.Sprintf("%.1fs", float64(num)/float64(denom))
			}
		}
	}
	// ISO
	if tag, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if v, err2 := tag.Int(0); err2 == nil {
			result.ISO = fmt.Sprintf("%d", v)
		}
	}
	// Shooting time
	if t, err := x.DateTime(); err == nil {
		result.TakenAt = t
	}
	// GPS
	if lat, lon, err := x.LatLong(); err == nil {
		result.Latitude = lat
		result.Longitude = lon
		result.HasGPS = true
	}
	// Dimensions from EXIF
	if tag, err := x.Get(exif.PixelXDimension); err == nil {
		if v, err2 := tag.Int(0); err2 == nil {
			result.Width = v
		}
	}
	if tag, err := x.Get(exif.PixelYDimension); err == nil {
		if v, err2 := tag.Int(0); err2 == nil {
			result.Height = v
		}
	}

	return result
}

// ProcessImage generates display, large, and placeholder assets from the source image.
func ProcessImage(data []byte) (*ProcessedAssets, error) {
	src, err := decodeSourceImage(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := src.Bounds()
	origW := bounds.Dx()

	// Generate large version (cleaned, quality-optimized)
	large := src
	if origW > largeMaxWidth {
		large = imaging.Resize(src, largeMaxWidth, 0, imaging.Lanczos)
	}
	largeBuf := &bytes.Buffer{}
	if err := jpeg.Encode(largeBuf, large, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, fmt.Errorf("failed to encode large: %w", err)
	}
	largeBounds := large.Bounds()

	// Generate display version
	display := src
	if origW > displayMaxWidth {
		display = imaging.Resize(src, displayMaxWidth, 0, imaging.Lanczos)
	}
	displayBuf := &bytes.Buffer{}
	if err := jpeg.Encode(displayBuf, display, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, fmt.Errorf("failed to encode display: %w", err)
	}
	displayBounds := display.Bounds()

	// Generate BlurHash placeholder
	thumb := imaging.Resize(src, 32, 0, imaging.Box)
	hash, err := blurhash.Encode(blurHashXComp, blurHashYComp, thumb)
	if err != nil {
		hash = ""
	}

	return &ProcessedAssets{
		Display:         displayBuf.Bytes(),
		DisplayWidth:    displayBounds.Dx(),
		DisplayHeight:   displayBounds.Dy(),
		Large:           largeBuf.Bytes(),
		LargeWidth:      largeBounds.Dx(),
		LargeHeight:     largeBounds.Dy(),
		PlaceholderHash: hash,
	}, nil
}

// ReadAll reads the multipart file fully into memory.
func ReadAll(r io.Reader) ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func cleanExifString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"")
	s = strings.TrimSpace(s)
	return s
}

func decodeSourceImage(data []byte) (image.Image, error) {
	src, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err == nil {
		return src, nil
	}
	heifImage, heifErr := goheif.Decode(bytes.NewReader(data))
	if heifErr == nil {
		return heifImage, nil
	}
	return nil, err
}

func decodeImageConfig(data []byte) (image.Config, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err == nil {
		return cfg, nil
	}
	heifCfg, heifErr := goheif.DecodeConfig(bytes.NewReader(data))
	if heifErr == nil {
		return heifCfg, nil
	}
	return image.Config{}, err
}

func decodeExifMetadata(data []byte) (*exif.Exif, error) {
	if x, err := exif.Decode(bytes.NewReader(data)); err == nil {
		return x, nil
	}
	exifBytes, err := goheif.ExtractExif(bytes.NewReader(data))
	if err != nil || len(exifBytes) == 0 {
		return nil, err
	}
	return exif.Decode(bytes.NewReader(exifBytes))
}
