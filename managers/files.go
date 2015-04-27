package managers

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

// Loads png or jpg images from files
func LoadImageFromFile(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return GetImageDecodingFunction(path)(reader)
}

// Gets decoding function for extension
func GetImageDecodingFunction(fileName string) func(io.Reader) (image.Image, error) {
	decode := jpeg.Decode
	if strings.HasSuffix(fileName, "png") {
		decode = png.Decode
	}
	return decode
}

// Gets encoding function for extension
func GetImageEncodingFunction(fileName string) func(w io.Writer, m image.Image) error {
	encode := encodeJpeg
	if strings.HasSuffix(fileName, "png") {
		encode = png.Encode
	}
	return encode
}

func SaveImage(file string, img image.Image) error {
	os.Remove(file)
	fl, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fl.Close()

	b := bufio.NewWriter(fl)
	err = GetImageEncodingFunction(file)(b, img)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func encodeJpeg(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, &jpeg.Options{100})
}
