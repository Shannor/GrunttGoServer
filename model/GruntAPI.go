package model

import (
	"scrapper/workers"

	"github.com/julienschmidt/httprouter"
)

type GrunttAPI interface {
	GetAllComics() httprouter.Handle
	GetPopularComics() httprouter.Handle
	GetComicDescription() httprouter.Handle
	GetComicChapterList() httprouter.Handle
	GetChapterPages() httprouter.Handle
	GetSearchCategories() httprouter.Handle
	Search() httprouter.Handle
}

type api struct {
	ComicExtra workers.Webcrawler
	ReadComics workers.Webcrawler
}
