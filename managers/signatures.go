package managers

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"strconv"
	"time"

	"code.google.com/p/draw2d/draw2d"
	//Decided to use this
	"log"
	"os"

	curveapi "github.com/maciekmm/curveapi/models"
	"github.com/maciekmm/curvesignatures/models"
)

const (
	// API_URL is a http address of curveapi
	//API_URL = "http://localhost:2000/"
	API_URL = "http://curveapi.cf/"
)

// ChampionCrown, Premium, GameLogo are images used frequently by submodules
var ChampionCrown, Premium, GameLogo image.Image
var signatureFiles = models.New()

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

// GetPlayerData gets player data from api site by its id
func GetPlayerData(playerID int) (*curveapi.Profile, error) {
	req, err := httpClient.Get(API_URL + "user/" + strconv.Itoa(playerID))
	if err != nil || req.StatusCode != 200 {
		return nil, errors.New("Could not load player's profile.")
	}
	return loadPlayerData(req)
}

// GetPlayerDataByName gets player data from api site by its name
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

// GetPlayerFolder gets folder player is associated with
func GetPlayerFolder(profileID int) (dir string, avatarDir string) {
	dir = "./players/" + strconv.Itoa(profileID) + "/"
	avatarDir = dir + "avatars/"
	os.Mkdir(dir, 0750)
	os.Mkdir(avatarDir, 0750)
	return
}

// GetSignature gets signatures from cache of renders it
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
			SaveImageByMime("image/png", signatureFileName, signature)
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
			SaveImageByMime("image/png", signatureFileName, signature)
		}()
	}) {
		fil.Done()
	}
	return signature, nil
}
