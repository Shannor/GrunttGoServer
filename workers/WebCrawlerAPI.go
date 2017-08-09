package workers

import "github.com/PuerkitoBio/goquery"

//Webcrawler Interface for all services that must be implemented
type Webcrawler interface {
	GetTag() string
	CreateAllComicsURL() string
	CreatePopularComicsURL(int) string
	CreateComicChapterListURL(string) string
	CreateChapterPagesURL(string, int) string
	CreateComicDescriptionURL(string) string
	// CreateSearchURL() (string, error)
	GetAllComics(*goquery.Document) (Comics, error)
	GetPopularComics(*goquery.Document) (PopularComics, error)
	// GetComicChapterList(*goquery.Document, string) (Chapters, error)
	// GetChapterPages(*goquery.Document, string, int) ([]string, error)
	// GetSearchCategories() ([]string, error)
	// GetComicDescription() error
	// Search(*goquery.Document) (SearchResults, error)
}

//ComicExtra Struct for ComicExtra site webcrawler
type ComicExtra struct {
	BaseURL string
	Tag     string
}

//GetComicExtraInstance Returns the webcrawler interaface for ComicExtra
func GetComicExtraInstance() Webcrawler {
	return &ComicExtra{BaseURL: "http://www.comicextra.com/", Tag: "ce"}
}

//ReadComics Struct for ReadComics.website webcrawler
type ReadComics struct {
	BaseURL string
	Tag     string
}

//GetReadComicsInstance Returns a webcrawler instance for ReadComics.Website
func GetReadComicsInstance() Webcrawler {
	return &ReadComics{BaseURL: "http://readcomics.website/", Tag: "rcw"}
}

//Comic Base struct for basic comics
type Comic struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Category string `json:"category"`
}

//Comics Struct for a slice of Comic
type Comics []Comic

//PopularComic Struct for popular comics
type PopularComic struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	Img        string `json:"img"`
	IssueCount int    `json:"issueCount"`
}

//PopularComics Struct for a slice of Popularcomics
type PopularComics []PopularComic

//Chapter Struct for chapter information
type Chapter struct {
	ChapterName string `json:"chapterName"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}

//Chapters Type for slice of Chapter
type Chapters []Chapter

//SearchResult Struct for response from a search
type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Img   string `json:"img"`
	Genre string `json:"genre"`
}

//SearchResults Type for slice of SearchResult
type SearchResults []SearchResult

//Description Struct for information on a particular comic
type Description struct {
	Description   string `json:"description"`
	LargeImg      string `json:"largeImg"`
	FormatedName  string `json:"formatedName"`
	Name          string `json:"name"`
	AlternateName string `json:"alternateName"`
	ReleaseYear   string `json:"releaseYear"`
	Status        string `json:"status"`
	Author        string `json:"author"`
	Genre         string `json:"genre"`
}
