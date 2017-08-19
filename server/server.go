package server

import (
	"fmt"
	"net/http"
	"scrapper/model"
	"scrapper/workers"

	"github.com/julienschmidt/httprouter"
)

func init() {

	api := model.GetAPIInstance(
		workers.GetComicExtraInstance(),
		workers.GetReadComicsInstance(),
	)

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome to Gruntt-Comics Backend\n")
	})
	router.GET("/all-comics", api.GetAllComics())
	router.GET("/popular-comics/:pageNumber", api.GetPopularComics())
	router.GET("/chapter-list/:comicName", api.GetComicChapterList())
	router.GET("/chapter-pages/:comicName/:chapterNumber", api.GetChapterPages())
	router.GET("/description/:comicName", api.GetComicDescription())
	router.GET("/categories", api.GetSearchCategories())
	http.Handle("/", router)
}
