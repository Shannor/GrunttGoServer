package comicutils

import (
	"net/http"
	"scrapper/model"
	"scrapper/workers"

	"github.com/PuerkitoBio/goquery"
)

//CreateAllComicsURL function to create the URL for each source
func CreateAllComicsURL(urlParam string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + "comic-list", nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "changeMangaList?type=text", nil
	default:
		return "", &model.RequestError{ProvidedParam: urlParam}
	}
}

func CreatePopularComicsURL(urlParam string, pageNumber string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + "popular-comic/" + pageNumber, nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "filterList?page=" +
			pageNumber + "&sortBy=views&asc=false", nil
	default:
		return "", &model.RequestError{ProvidedParam: urlParam}
	}
}

func CreateChapterURL(urlParam string, comicName string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + "comic/" + comicName, nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "comic/" + comicName, nil
	default:
		return "", &model.RequestError{ProvidedParam: urlParam}

	}
}

func CreateReadComicURL(urlParam string, comicName string, chapterNumber string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + comicName + "/chapter-" + chapterNumber, nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "comic/" + comicName + "/" + chapterNumber, nil
	default:
		return "", &model.RequestError{ProvidedParam: urlParam}

	}
}

//TODO: Add err returns to all Get methods
func GetAllComics(doc *goquery.Document, param string) []model.Comic {
	switch param {
	case workers.ComicExtraURLParam:
		return workers.GetAllComicsComicExtra(doc)
	case workers.ReadComicsURLParam:
		return workers.GetAllComicsReadComics(doc)
	default:
		return nil
	}
}

func GetPopularComics(doc *goquery.Document, param string) []model.PopularComic {
	switch param {
	case workers.ComicExtraURLParam:
		return workers.GetPopularComicsComicExtra(doc)
	case workers.ReadComicsURLParam:
		return workers.GetPopularComicsReadComics(doc)
	default:
		return nil
	}
}

func GetChapters(doc *goquery.Document, r *http.Request, param string, comicName string) []model.Chapter {
	switch param {
	case workers.ComicExtraURLParam:
		return workers.GetChaptersComicExtra(doc, r, comicName)
	case workers.ReadComicsURLParam:
		return workers.GetChaptersReadComics(doc)
	default:
		return nil
	}
}

func GetChapterImages(doc *goquery.Document, r *http.Request, param string, url string) []string {
	switch param {
	case workers.ComicExtraURLParam:
		return nil
	case workers.ReadComicsURLParam:
		return nil
	default:
		return nil
	}
}
