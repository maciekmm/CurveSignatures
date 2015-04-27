package managers

import (
	"bytes"
	"code.google.com/p/draw2d/draw2d"
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"strconv"
	"time"
	//Decided to use this
	curveapi "github.com/maciekmm/curveapi/models"
	"github.com/maciekmm/curvesignatures/models"
	"log"
	"os"
)

const (
	//API_URL = "http://localhost:2000/"
	API_URL = "http://curveapi.cf/"
)

var ChampionCrown, Premium, GameLogo image.Image
var signatureFiles *models.Guardian = models.New()

func init() {
	err := os.Mkdir("./players", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatalln("Error occured while creating players folder.")
		os.Exit(1)
	}
	checkErrors := func(ex error) {
		if ex != nil {
			log.Fatalln("Could not load file.", ex)
			os.Exit(1)
		}
	}
	ChampionCrown, err = LoadImageFromFile("./assets/crown.png")
	checkErrors(err)
	Premium, err = LoadImageFromFile("./assets/premium.png")
	checkErrors(err)
	GameLogo, err = LoadImageFromFile("./assets/game-logo.png")
	checkErrors(err)
	draw2d.SetFontFolder("./public")
}

// Gets player data from api site
func GetPlayerData(playerId int) (*curveapi.Profile, error) {
	req, err := httpClient.Get(API_URL + "user/" + strconv.Itoa(playerId))
	if err != nil || req.StatusCode != 200 {
		return nil, errors.New("Could not load player's profile.")
	}
	return loadPlayerData(req)
}

func GetPlayerDataByName(playerName string) (*curveapi.Profile, error) {
	req, err := httpClient.Get(API_URL + "username/" + playerName)
	if err != nil || req.StatusCode != 200 {
		return nil, errors.New("Could not load player's profile.")
	}
	return loadPlayerData(req)
}

func loadPlayerData(req *http.Response) (*curveapi.Profile, error) {
	var profile curveapi.Profile
	defer req.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	err := json.Unmarshal(buf.Bytes(), &profile)
	if err != nil {
		return nil, errors.New("Could not load player's profile.")
	}
	return &profile, nil
}

// Gets folder player is associated with
func GetPlayerFolder(profileId int) (dir string, avatarDir string) {
	dir = "./players/" + strconv.Itoa(profileId) + "/"
	avatarDir = dir + "avatars/"
	os.Mkdir(dir, 0750)
	os.Mkdir(avatarDir, 0750)
	return
}

// Gets signatures from cache of renders it
// May be blocking for some period of time
func GetSignature(params models.RequestParameters) (image.Image, error) {
	dir, _ := GetPlayerFolder(params.PlayerID)
	signatureFileName := dir + params.LayoutName + params.Ranks.CombineRanks()
	fil := signatureFiles.GetOrCreate(signatureFileName)
	rw := fil.Lock
	rw.RLock()
	signature, err := LoadImageFromFile(signatureFileName)
	rw.RUnlock()

	if err != nil {
		defer fil.Done()
		rw.Lock()
		profile, err := GetPlayerData(params.PlayerID)
		if err != nil {
			rw.Unlock()
			return nil, err
		}
		signature, err := params.Layout.Render(params.Ranks, profile)

		if err != nil {
			rw.Unlock()
			return nil, err
		}

		go func() {
			defer rw.Unlock()
			SaveImage(signatureFileName, signature)
		}()
		return signature, nil
	}

	if !fil.OnceCaller.Do(func() {
		go func() {
			defer fil.Done()
			profile, err := GetPlayerData(params.PlayerID)
			if err != nil {
				log.Println(err)
				return
			}

			rw.RLock()
			file, err := os.Open(signatureFileName)

			//Check if there's a pending update for signature
			//TODO Somehow check if there was a change in data (curveapi)
			if err == nil {
				stat, err := file.Stat()
				if err == nil {
					lastUpdate := time.Unix(profile.LastUpdate, 0).UTC()
					if stat.ModTime().UTC().After(lastUpdate) {
						rw.RUnlock()
						return
					}
				}
			}
			rw.RUnlock()

			//Render signature
			signature, err := params.Layout.Render(params.Ranks, profile)
			if err != nil {
				log.Println("Error occured while generating signature ", err)
				return
			}
			//Save rendered file
			rw.Lock()
			defer rw.Unlock()
			SaveImage(signatureFileName, signature)
		}()
	}) {
		fil.Done()
	}
	return signature, nil
}
