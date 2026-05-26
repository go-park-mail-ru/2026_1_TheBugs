package photo

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"golang.org/x/image/draw"
)

func GetKeyFromPath(path string) string {
	return strings.TrimPrefix(path, "/")
}

func MakeUrlFromPath(path, publicHost, bucket string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return fmt.Sprintf("%s/%s%s", publicHost, bucket, path)
}
func ResizeAndCropJPEG(input []byte, width, height int, quality int) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		log.Printf("image.Decode: %s", err)
		return nil, fmt.Errorf("image.Decode: %w", err)
	}

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	targetRatio := float64(width) / float64(height)
	srcRatio := float64(srcW) / float64(srcH)

	var cropX, cropY, cropW, cropH int
	if srcRatio > targetRatio {
		cropH = srcH
		cropW = int(float64(srcH) * targetRatio)
		cropX = (srcW - cropW) / 2
		cropY = 0
	} else {
		cropW = srcW
		cropH = int(float64(srcW) / targetRatio)
		cropY = (srcH - cropH) / 2
		cropX = 0
	}

	cropped := image.NewRGBA(image.Rect(0, 0, cropW, cropH))
	draw.Draw(cropped, cropped.Bounds(), src, image.Point{cropX, cropY}, draw.Src)

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), cropped, cropped.Bounds(), draw.Over, nil)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, dst, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("jpeg.Encode: %w", err)
	}

	return buf.Bytes(), nil
}
