package managers

import (
	curveapi "github.com/maciekmm/curveapi/models"
	"github.com/maciekmm/curvesignatures/models"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	defaultAvatar image.Image
	avatarFiles   *models.Guardian = models.New()
	httpClient    *http.Client     = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
)

func init() {
	var err error
	defaultAvatar, err = LoadImageFromFile("./assets/default-avatar.png")
	if err != nil {
		log.Fatalln("Could not load file.", err)
		os.Exit(1)
	}
}

// Gets player avatar
func GetPlayerAvatar(profile curveapi.Profile) <-chan image.Image {
	img := make(chan image.Image)
	if len(profile.Picture) == 0 {
		go func() {
			img <- defaultAvatar
		}()
		return img
	}

	go func() {
		_, avDir := GetPlayerFolder(profile.UID)
		urlParts := strings.Split(profile.Picture, "/")
		fileName := urlParts[len(urlParts)-1]

		fle := avatarFiles.GetOrCreate(avDir + fileName)
		defer fle.Done()
		rw := fle.Lock
		rw.RLock()
		avatarFile, err := os.Open(avDir + fileName)

		if err == nil {
			output, err := GetImageDecodingFunction(fileName)(avatarFile)
			if err != nil {
				rw.RUnlock()
				rw.Lock()
				os.Remove(avDir + fileName)
			} else {
				rw.RUnlock()
				img <- output
				return
			}
		}
		rw.RUnlock()
		if err != nil {
			req, err := httpClient.Get(profile.Picture)

			if err == nil {
				defer req.Body.Close()
				output, err := GetImageDecodingFunction(fileName)(req.Body)
				if err == nil {
					img <- output
					rw.Lock()
					defer rw.Unlock()
					fle.OnceCaller.Do(func() {
						err := SaveImage(avDir+fileName, output)
						if err != nil {
							log.Println(err)
						}
					})
				} else {
					img <- defaultAvatar
				}
			} else {
				log.Println("Couldn't fetch avatar from curvefever servers.", err)
				output, err := getRandomAvatar(avDir)
				if err != nil {
					img <- defaultAvatar
					return
				}
				img <- output
				return
			}
			return
		}
	}()
	return img
}

func getRandomAvatar(avDir string) (image.Image, error) {
	mtchsm, err := filepath.Glob(avDir + "picture-*")
	if err != nil || len(mtchsm) < 1 {
		return nil, err
	}

	ra := avatarFiles.GetOrCreate(mtchsm[0])
	defer ra.Done()
	ra.Lock.RLock()
	defer ra.Lock.RUnlock()
	output, err := LoadImageFromFile(mtchsm[0])
	if err != nil {
		return nil, err
	}
	return output, nil
}
