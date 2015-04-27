package main

import (
	"github.com/julienschmidt/httprouter"
	curveapi "github.com/maciekmm/curveapi/models"
	"github.com/maciekmm/curvesignatures/managers"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var cachedTemplates *template.Template = template.Must(template.New("").Funcs(template.FuncMap{"rankBeautify": managers.ConvertToReadableName, "getRegionIcon": getRegionIconURL, "divAndCeil": divideAndCeil, "getRegionName": getRegionName}).ParseFiles("templates/includes/footer.tmpl", "templates/includes/header.tmpl", "templates/includes/route-chooser.tmpl", "templates/index.tmpl", "templates/player.tmpl", "templates/create/create.tmpl", "templates/create/created.tmpl", "templates/create/search.tmpl", "templates/404.tmpl", "templates/api.tmpl"))

type M map[string]interface{}

func divideAndCeil(a, b int) string {
	return strconv.FormatInt(int64(math.Ceil(float64(a)/float64(b))), 10)
}

func getRegionIconURL(rank string) string {
	reg, err := managers.GetRegionFromRank(rank)
	if err != nil {
		return ""
	}
	return "/assets/" + reg.Suffix + ".png"
}

func getRegionName(rank string) string {
	reg, err := managers.GetRegionFromRank(rank)
	if err != nil {
		return ""
	}
	return reg.FullName
}

func mainPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := cachedTemplates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		log.Println(err)
	}
}

func notFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(404)
	err := cachedTemplates.ExecuteTemplate(rw, "404", nil)
	if err != nil {
		log.Println(err)
	}
}

func createView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if (err != nil || len(r.PostForm["name"]) == 0 && len(r.PostForm["player-id"]) == 0) || r.Method != "POST" {
		cachedTemplates.ExecuteTemplate(w, "create-search", M{"ShowMessage": false})
		return
	}

	var player *curveapi.Profile
	if len(r.PostForm["name"]) != 0 {
		player, err = managers.GetPlayerDataByName(r.PostFormValue("name"))
	} else {
		playerId, err := strconv.Atoi(r.PostFormValue("player-id"))
		if err == nil {
			player, err = managers.GetPlayerData(playerId)
		}
	}

	if err != nil {
		cachedTemplates.ExecuteTemplate(w, "create-search", M{
			"ShowMessage": true,
			"Message":     "Player not found, make sure you spelled it right.",
		})
		return
	}
	if len(r.PostForm["rank"]) != 0 && len(r.PostForm["layout"]) != 0 {
		layout := managers.GetLayoutById(r.PostFormValue("layout"))
		err = cachedTemplates.ExecuteTemplate(w, "created", M{
			"ranks":       r.PostForm["rank"],
			"layout":      layout,
			"player":      player,
			"link":        URL + "/img/" + strconv.Itoa(player.UID) + "/" + layout.Name() + "/" + "+" + strings.Join(r.PostForm["rank"], "+") + ".png",
			"profileLink": URL + "/profile/" + strconv.Itoa(player.UID),
		})
		return
	}

	err = cachedTemplates.ExecuteTemplate(w, "create", M{
		"layouts": managers.GetRegisteredLayoutNames(),
		"profile": player,
	})
	if err != nil {
		log.Println(err)
	}
}

func playerView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	playerId, err := strconv.Atoi(ps.ByName("player"))
	if err != nil {
		notFound(w, r)
		return
	}
	player, err := managers.GetPlayerData(playerId)
	if err != nil {
		notFound(w, r)
		return
	}
	for k, v := range player.Ranks {
		if v.Bonus == 500 {
			delete(player.Ranks, k)
		}
	}
	err = cachedTemplates.ExecuteTemplate(w, "player", player)
	if err != nil {
		log.Println(err)
	}
}

func apiView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := cachedTemplates.ExecuteTemplate(w, "api", nil)
	if err != nil {
		log.Println(err)
	}
}
