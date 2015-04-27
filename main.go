package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/maciekmm/curvesignatures/managers"
	"github.com/maciekmm/curvesignatures/models/layouts"
	"log"
	"net/http"
	"os"
)

const (
	URL = "https://signatures.cf"
)

func init() {
	f, err := os.OpenFile("./curvatures.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	log.SetOutput(f)

	managers.RegisterLayout("default", layouts.GetDefaultLayout())
	managers.RegisterLayout("userbar", layouts.GetUserbarLayout())
}

func main() {
	router := httprouter.New()

	router.GET("/img/:user/:layout/:ranks", requestSignature)
	router.GET("/avatar/:user", requestAvatar)
	router.GET("/profile/:player", playerView)
	router.GET("/api", apiView)
	router.GET("/create", createView)
	router.POST("/create", createView)
	router.GET("/", mainPage)

	router.HandlerFunc("GET", "/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/favicon.ico")
	})

	router.HandlerFunc("GET", "/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/robots.txt")
	})

	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.ServeFiles("/assets/*filepath", http.Dir("./assets"))

	router.NotFound = notFound
	log.Fatalln(http.ListenAndServe(":3000", router))
}
