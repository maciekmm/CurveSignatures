package models

import (
	curveapi "github.com/maciekmm/curveapi/models"
	"image"
)

type Layout interface {
	Name() string
	Width() int
	Height() int
	MaxRanks() int
	Render(conf *Configuration, profile *curveapi.Profile) (image.Image, error)
}
