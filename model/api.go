package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scrapper/utils"
	"scrapper/workers"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (api *api) GetWebcrawler(tag string) (workers.Webcrawler, error) {
	switch tag {
	case api.ComicExtra.GetTag():
		return api.ComicExtra, nil
	case api.ReadComics.GetTag():
		return api.ReadComics, nil
	default:
		return nil, fmt.Errorf("No tag matching %s", tag)
	}

}

func (api *api) GetAllComics() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		source := r.URL.Query().Get("source")
		if source == "" {
			http.Error(w, "No provided source in url", http.StatusBadRequest)
			return
		}

		webcrawler, err := api.GetWebcrawler(source)
		if err != nil {
			msg := fmt.Sprintf("No matching webcrawler source for %s", source)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		url := webcrawler.CreateAllComicsURL()

		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			msg := fmt.Sprintf("Error with GoQuery. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		comics, err := webcrawler.GetAllComics(doc)
		if err != nil {
			msg := fmt.Sprintf("Error with Get All Comics. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(comics)
		if err != nil {
			msg := fmt.Sprintf("Error with Marshaling Comics. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func (api *api) GetPopularComics() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		source := r.URL.Query().Get("source")
		pageNumber := ps.ByName("pageNumber")

		if source == "" {
			http.Error(w, "No provided source in url", http.StatusBadRequest)
			return
		}

		pg, err := strconv.Atoi(pageNumber)
		if err != nil {
			http.Error(w, "Page Number is not a valid in url", http.StatusBadRequest)
			return
		}

		webcrawler, err := api.GetWebcrawler(source)
		if err != nil {
			msg := fmt.Sprintf("No matching webcrawler source for %s", source)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		url := webcrawler.CreatePopularComicsURL(pg)

		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			msg := fmt.Sprintf("Error with GoQuery. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		comics, err := webcrawler.GetPopularComics(doc)
		if err != nil {
			msg := fmt.Sprintf("Error with Get Popular Comics. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(comics)
		if err != nil {
			msg := fmt.Sprintf("Error with Marshaling Comics. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func (api *api) GetComicChapterList() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		source := r.URL.Query().Get("source")
		comicName := ps.ByName("comicName")

		if source == "" {
			http.Error(w, "No provided source in url", http.StatusBadRequest)
			return
		}

		webcrawler, err := api.GetWebcrawler(source)
		if err != nil {
			msg := fmt.Sprintf("No matching webcrawler source for %s", source)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		url := webcrawler.CreateComicChapterListURL(comicName)

		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			msg := fmt.Sprintf("GetComicChapterList error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		numOfChapterPages := webcrawler.GetComicChapterListPageAmount(doc)

		chapters, err := webcrawler.GetComicChapterList(comicName, numOfChapterPages, r)
		if err != nil {
			msg := fmt.Sprintf("GetComicChapterList error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(chapters)
		if err != nil {
			msg := fmt.Sprintf("Error with Marshaling Comics. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func (api *api) GetChapterPages() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		source := r.URL.Query().Get("source")
		comicName := ps.ByName("comicName")
		chapterNumber := ps.ByName("chapterNumber")

		if source == "" {
			http.Error(w, "No provided source in url", http.StatusBadRequest)
			return
		}

		webcrawler, err := api.GetWebcrawler(source)
		if err != nil {
			msg := fmt.Sprintf("No matching webcrawler source for %s", source)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		chapterNum, err := strconv.Atoi(chapterNumber)
		if err != nil {
			msg := fmt.Sprintf("Provided chapter number is not a number.")
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		url := webcrawler.CreateChapterPagesURL(comicName, chapterNum)

		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			msg := fmt.Sprintf("GetChapterPages error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		numOfPages := webcrawler.GetNumberOfPages(doc)

		pages, err := webcrawler.GetChapterPages(comicName, chapterNum, numOfPages, r)
		if err != nil {
			msg := fmt.Sprintf("GetChapterPages error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(pages)
		if err != nil {
			msg := fmt.Sprintf("Error with Marshaling Pages. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func (api *api) GetComicDescription() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		source := r.URL.Query().Get("source")
		comicName := ps.ByName("comicName")

		if source == "" {
			http.Error(w, "No provided source in url", http.StatusBadRequest)
			return
		}

		webcrawler, err := api.GetWebcrawler(source)
		if err != nil {
			msg := fmt.Sprintf("No matching webcrawler source for %s", source)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		url := webcrawler.CreateComicDescriptionURL(comicName)

		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			msg := fmt.Sprintf("GetComicDescription error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		description, err := webcrawler.GetComicDescription(doc)
		if err != nil {
			msg := fmt.Sprintf("GetChapterPages error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		description.FormatedName = comicName

		body, err := json.Marshal(description)
		if err != nil {
			msg := fmt.Sprintf("Error with Marshaling Description. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func (api *api) GetSearchCategories() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		source := r.URL.Query().Get("source")

		if source == "" {
			http.Error(w, "No provided source in url", http.StatusBadRequest)
			return
		}

		webcrawler, err := api.GetWebcrawler(source)
		if err != nil {
			msg := fmt.Sprintf("No matching webcrawler source for %s", source)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		url := webcrawler.CreateSearchURL()

		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			msg := fmt.Sprintf("GetSearchCategories error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		searchOptions, err := webcrawler.GetSearchOptions(doc)
		if err != nil {
			msg := fmt.Sprintf("GetChapterPages error. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(searchOptions)
		if err != nil {
			msg := fmt.Sprintf("Error with Marshaling Description. Error: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
