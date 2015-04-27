package layouts

import (
	"code.google.com/p/draw2d/draw2d"
	curveapi "github.com/maciekmm/curveapi/models"
	"github.com/maciekmm/curvesignatures/managers"
	"github.com/maciekmm/curvesignatures/models"
	"image"
	"image/color"
	"image/draw"
	"strconv"
)

var defLayout *DefaultLayout

type DefaultLayout struct {
	background image.Image
}

func (dl DefaultLayout) Name() string {
	return "default"
}

func (dl DefaultLayout) Height() int {
	return 95
}

func (dl DefaultLayout) Width() int {
	return 300
}

func (dl DefaultLayout) MaxRanks() int {
	return 6
}

func GetDefaultLayout() DefaultLayout {
	if defLayout == nil {
		bg, _ := managers.LoadImageFromFile("./assets/default-background.png")
		defLayout = &DefaultLayout{bg}
	}
	return *defLayout
}

func (dl DefaultLayout) Render(conf *models.Configuration, profile *curveapi.Profile) (image.Image, error) {
	//Get player avatar
	avatarChan := managers.GetPlayerAvatar(*profile)

	//Essentials
	img := image.NewRGBA(image.Rect(0, 0, dl.Width(), dl.Height()))
	gc := draw2d.NewGraphicContext(img)
	defer gc.Close()
	gc.SetFontData(draw2d.FontData{"roboto", draw2d.FontFamilyMono, draw2d.FontStyleNormal})
	gc.Translate(0, 0)
	gc.SetFillColor(color.RGBA{15, 15, 15, 255})
	gc.SetStrokeColor(image.Black)

	//Draw background
	draw2d.Rect(gc, 0, 0, float64(dl.Width()), float64(dl.Height()))
	gc.DrawImage(dl.background)
	gc.Stroke()

	//Draw star and crown
	gc.Translate(95, 5)
	if profile.Premium {
		gc.DrawImage(managers.Premium)
	}
	if profile.Champion {
		gc.Translate(20, 0)
		gc.DrawImage(managers.ChampionCrown)
	}

	//Draw nickname
	gc.SetFillColor(image.White)
	gc.SetFontSize(13)
	gc.Translate(17, 13)
	wdth := gc.FillString(profile.Name)
	x, _ := gc.LastPoint()
	height := 3

	maxRanks := dl.MaxRanks()
	if dl.Width()-int(x)-int(wdth) < managers.GameLogo.Bounds().Dx() {
		height = dl.Height() - 27
		maxRanks -= 1
	}
	//Draw game logo
	gc.SetMatrixTransform(draw2d.NewIdentityMatrix())
	gc.Translate(float64(dl.Width())-73, float64(height))
	gc.DrawImage(managers.GameLogo)

	//Draw ranks
	gc.SetMatrixTransform(draw2d.NewIdentityMatrix())
	gc.Translate(100, 30)
	gc.SetFontSize(9)

	longest := float64(0)
	for k, v := range conf.Ranks {
		if k >= maxRanks {
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
		gc.DrawImage(rank.Icon)
		gc.Translate(18, 12)
		readable, _ := managers.ConvertToReadableName(v)
		stringWidth := gc.FillString(readable + ": " + strconv.Itoa(val.Rank))
		if stringWidth > longest {
			longest = stringWidth
		}
		gc.Translate(-18, 8)
		if (k+1)%3 == 0 {
			gc.Translate(longest+25, -60)
		}
	}

	//Draw avatar
	avatar := <-avatarChan
	gc.SetMatrixTransform(draw2d.NewIdentityMatrix())
	//log.Println((float64(height) - float64(avatar.Bounds().Size().Y)*scale) / 2.0)
	avatarWidth := float64(avatar.Bounds().Size().X)
	avatarHeight := float64(avatar.Bounds().Size().Y)
	gc.Translate(5+(float64(85)-avatarWidth)/2, 5+(float64(85)-avatarHeight)/2)
	gc.SetFillColor(color.RGBA{51, 178, 60, 255})
	gc.DrawImage(avatar)
	draw2d.DrawImage(avatar, img, gc.GetMatrixTransform(), draw.Over, draw2d.LinearFilter)
	gc.SetMatrixTransform(draw2d.NewIdentityMatrix())
	return img, nil
}
