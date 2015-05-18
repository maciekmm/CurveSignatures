package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maciekmm/curvesignatures/managers"
	"github.com/maciekmm/curvesignatures/models"
)

var ranksRegex = regexp.MustCompile("^([\\+]{1}([A-Za-z0-9_]+))+\\.png$")
var groupsRegex = regexp.MustCompile("(?:\\+([a-z0-9_]+))")

func requestAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		w.WriteHeader(400)
		serveJSON(w, message{"Invalid UserID"})
		return
	}
	profile, err := managers.GetPlayerData(userID)
	avatarChan := managers.GetPlayerAvatar(*profile)
	img := <-avatarChan
	if img != nil {
		w.Header().Set("Content-Type", "image/png")
		buffer := new(bytes.Buffer)
		if err := png.Encode(buffer, img); err != nil {
			panic("Error occured while encoding image")
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("Error occured while writing image" + err.Error())
		}
	} else {
		http.Error(w, "Not Found", 404)
	}
}

func requestSignature(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		w.WriteHeader(400)
		serveJSON(w, message{"Invalid UserID"})
		return
	}

	layout := ps.ByName("layout")

	layoutFunc := managers.GetLayoutByID(layout)
	if layoutFunc == nil {
		w.WriteHeader(404)
		serveJSON(w, message{"Layout not found"})
		return
	}

	if !ranksRegex.MatchString(ps.ByName("ranks")) {
		w.WriteHeader(400)
		serveJSON(w, message{"Wrong syntax"})
		return
	}

	reqRanks := groupsRegex.FindAllStringSubmatch(ps.ByName("ranks"), -1)

	var pureRanks = make([]string, len(reqRanks))
	for i := 0; i < len(reqRanks); i++ {
		pureRanks[i] = reqRanks[i][1]
	}

	config := &models.Configuration{Ranks: pureRanks}
	reqs := models.RequestParameters{LayoutName: layout, PlayerID: userID, Layout: layoutFunc, Ranks: config}

	src, err := managers.GetSignature(reqs)

	if err != nil {
		w.WriteHeader(404)
		serveJSON(w, &message{err.Error()})
		return
	}

	img, castOk := src.(image.Image)
	if !castOk {
		panic("Could not assert image")
	}

	w.Header().Set("Content-Type", "image/png")
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		panic("Error occured while encoding image")
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("Error occured while writing image" + err.Error())
	}
}

type message struct {
	Message string `json:"message"`
}

func serveJSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "text/json")
	result, _ := json.Marshal(data)
	rw.Write(result)
}
