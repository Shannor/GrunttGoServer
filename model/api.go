package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"scrapper/utils"
	"scrapper/workers"

	"github.com/julienschmidt/httprouter"
)

func (api *api) GetWebcrawler(tag string) (workers.Webcrawler, error) {
	switch tag {
	case api.ComicExtra.GetTag():
		return api.ComicExtra, nil
	case api.ReadComics.GetTag():
		return api.ReadComics, nil
	default:
		return nil, errors.New("No webcrawler that matches that tag")
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
