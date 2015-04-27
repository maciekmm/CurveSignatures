package layouts

import (
	"code.google.com/p/draw2d/draw2d"
	curveapi "github.com/maciekmm/curveapi/models"
	"github.com/maciekmm/curvesignatures/managers"
	"github.com/maciekmm/curvesignatures/models"
	"image"
	"image/color"
	"strconv"
)

var userbarLayout *UserbarLayout

type UserbarLayout struct {
	background image.Image
}

func (dl UserbarLayout) Name() string {
	return "userbar"
}

func (dl UserbarLayout) Height() int {
	return 19
}

func (dl UserbarLayout) Width() int {
	return 350
}

func (dl UserbarLayout) MaxRanks() int {
	return 2
}

func GetUserbarLayout() UserbarLayout {
	if userbarLayout == nil {
		bg, _ := managers.LoadImageFromFile("./assets/userbar-background.png")
		userbarLayout = &UserbarLayout{bg}
	}
	return *userbarLayout
}

func (dl UserbarLayout) Render(conf *models.Configuration, profile *curveapi.Profile) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, dl.Width(), dl.Height()))
	gc := draw2d.NewGraphicContext(img)
	defer gc.Close()
	gc.SetFontData(draw2d.FontData{"visitor2", draw2d.FontFamilyMono, draw2d.FontStyleNormal})
	gc.Translate(0, 0)
	gc.SetFillColor(color.RGBA{15, 15, 15, 255})
	gc.SetStrokeColor(image.Black)

	//Draw background
	draw2d.Rect(gc, 0, 0, float64(dl.Width()), float64(dl.Height()))
	gc.DrawImage(dl.background)
	gc.Stroke()

	//Draw crown and star
	gc.Translate(4, 2)
	gc.SetFontSize(10)
	if profile.Premium {
		gc.DrawImage(managers.Premium)
	}
	if profile.Champion {
		gc.Translate(20, 0)
		gc.DrawImage(managers.ChampionCrown)
	}

	//Draw username
	gc.SetFillColor(image.White)
	gc.Translate(20, 13)
	gc.SetLineWidth(1)
	gc.StrokeString(profile.Name)
	width := gc.FillString(profile.Name)
	gc.Translate(width+20, -2)
	//Draw ranks
	for k, v := range conf.Ranks {
		if k >= dl.MaxRanks() {
			break
		}
		if len(v) == 0 {
			continue
		}
		val, ok := profile.Ranks[v]
		if !ok {
			continue
		}
		rank, err := managers.GetRegionFromRank(v)
		if err != nil {
			continue
		}
		readable, _ := managers.ConvertToReadableName(v)
		gc.StrokeString(rank.Suffix + " " + readable + ": " + strconv.Itoa(val.Rank))
		stringWidth := gc.FillString(rank.Suffix + " " + readable + ": " + strconv.Itoa(val.Rank))
		gc.Translate(stringWidth+10, 0)
	}

	//Draw curvefever logo
	gc.SetMatrixTransform(draw2d.NewIdentityMatrix())
	gc.Translate(293, 1)
	gc.Scale(0.8, 0.8)
	gc.DrawImage(managers.GameLogo)
	return img, nil
}
