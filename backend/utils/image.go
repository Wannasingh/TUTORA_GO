package utils

import (
	"bytes"
	"image"
	_ "image/gif"  // Register GIF decoder
	"image/jpeg"
	_ "image/png"  // Register PNG decoder
	"io"
)

// OptimizeImage decodes PNG, JPEG, or GIF images and re-encodes them as optimized JPEGs
func OptimizeImage(r io.Reader, contentType string) (io.Reader, string, error) {
	// If it's not a common image format, bypass optimization
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return r, contentType, nil
	}

	// Decode the image
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}

	// Create a buffer to write the optimized JPEG
	buf := new(bytes.Buffer)

	// Encode to JPEG with 75% quality (high compression, minimal quality loss)
	options := &jpeg.Options{Quality: 75}
	err = jpeg.Encode(buf, img, options)
	if err != nil {
		return nil, "", err
	}

	return buf, "image/jpeg", nil
}
