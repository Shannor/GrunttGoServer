package model

import (
	"scrapper/workers"

	"github.com/julienschmidt/httprouter"
)

//GrunttAPI API for the RESTful calls
type GrunttAPI interface {
	GetWebcrawler(source string) (workers.Webcrawler, error)
	GetAllComics() httprouter.Handle
	GetPopularComics() httprouter.Handle
	GetComicChapterList() httprouter.Handle
	GetChapterPages() httprouter.Handle
	GetComicDescription() httprouter.Handle
	GetSearchCategories() httprouter.Handle
	// Search() httprouter.Handle
}

type api struct {
	ComicExtra workers.Webcrawler
	ReadComics workers.Webcrawler
}

//GetAPIInstance returns the API for the restful calls
func GetAPIInstance(ce workers.Webcrawler, rcw workers.Webcrawler) GrunttAPI {
	return &api{
		ComicExtra: ce,
		ReadComics: rcw,
	}
}
