package managers

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// ImageFunctions provides standarized functions for image encoding/decoding
type ImageFunctions struct {
	Extensions []string
	Decoder    func(io.Reader) (image.Image, error)
	Encoder    func(w io.Writer, m image.Image) error
}

var (
	mimeMappings = map[string]ImageFunctions{
		"image/png": ImageFunctions{
			Extensions: []string{"png"},
			Decoder:    png.Decode,
			Encoder:    png.Encode,
		},
		"image/jpeg": ImageFunctions{
			Extensions: []string{"jpg,jpeg"},
			Decoder:    jpeg.Decode,
			Encoder: func(w io.Writer, m image.Image) error {
				return jpeg.Encode(w, m, &jpeg.Options{Quality: 90})
			},
		},
		"image/gif": ImageFunctions{
			Extensions: []string{"gif"},
			Decoder:    gif.Decode,
			Encoder: func(w io.Writer, m image.Image) error {
				return gif.Encode(w, m, &gif.Options{NumColors: 256})
			},
		},
	}
)

func isSupported(mimeType string) bool {
	_, ok := mimeMappings[mimeType]
	return ok
}

// GetImageFunctionForExtension returns default ImageFunctions for given extension
func GetImageFunctionForExtension(extension string) ImageFunctions {
	for _, funcs := range mimeMappings {
		for _, ext := range funcs.Extensions {
			if ext == extension {
				return funcs
			}
		}
	}
	return mimeMappings["image/png"]
}

// LoadImageFromFile loads png or jpg images from files
func LoadImageFromFile(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return LoadImage(reader)
}

// LoadImage loads image from io.Reader
// The reader is not closed here, you have to call .Close() from the calling place
func LoadImage(reader io.Reader) (image.Image, error) {
	bts, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	mime := http.DetectContentType(bts)
	if !isSupported(mime) {
		return nil, errors.New("MimeType of image is not supported " + mime)
	}
	copiedReader := bytes.NewReader(bts)
	return GetImageFunctions(mime).Decoder(copiedReader)
}

// GetImageFunctions gets decoding function for mimeType
func GetImageFunctions(mimeType string) ImageFunctions {
	val, ok := mimeMappings[mimeType]
	if !ok {
		return mimeMappings["image/png"]
	}
	return val
}

// SaveImageByMime saves image with encoder based on it's mime type
func SaveImageByMime(mimeType string, file string, img image.Image) error {
	imgFuncs := GetImageFunctions(mimeType)
	return SaveImage(imgFuncs, file, img)
}

// SaveImage saves image with a given encoder
func SaveImage(imageFuncs ImageFunctions, file string, img image.Image) error {
	os.Remove(file)
	fl, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fl.Close()

	b := bufio.NewWriter(fl)
	err = imageFuncs.Encoder(b, img)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}
