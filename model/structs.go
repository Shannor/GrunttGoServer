package model

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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
	ComicExtra Webcrawler
	ReadComics Webcrawler
}

type Webcrawler interface {
	CreateAllComicsURL() string
	CreatePopularComicsURL(int) (string, error)
	CreateComicChapterListURL(string) (string, error)
	CreateChapterPagesURL(string, int) (string, error)
	CreateComicDescriptionURL(string) (string, error)
	CreateSearchURL() (string, error)
	GetAllComics(*goquery.Document) (Comics, error)
	GetPopularComics(*goquery.Document) (PopularComics, error)
	GetComicChapterList(*goquery.Document, string) (Chapters, error)
	GetChapterPages(*goquery.Document, string, int) ([]string, error)
	GetSearchCategories() httprouter.Handle
	GetComicDescription() httprouter.Handle
	Search(*goquery.Document) httprouter.Handle
}

type ComicExtra struct {
	BaseURL string
}

type ReadComics struct {
	BaseURL string
}

type Comic struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Category string `json:"category"`
}
type Comics []Comic

type PopularComic struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	Img        string `json:"img"`
	IssueCount int    `json:"issueCount"`
}
type PopularComics []PopularComic

type Chapter struct {
	ChapterName string `json:"chapterName"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}
type Chapters []Chapter

type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Img   string `json:"img"`
	Genre string `json:"genre"`
}
type SearchResults []SearchResult

// RequestError  error for when the provided param doesn't match
type RequestError struct {
	ProvidedParam string
}

//ResponseError error for when there is a response but not a 200
type ResponseError struct {
	ResponseCode int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Provided Param( '%s' ) does not match any options.",
		e.ProvidedParam)
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Error Response code: %d", e.ResponseCode)
}

//ErrorHandler function to write all errors to the page
func ErrorHandler(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}
