package workers

import (
	"log"
	"net/http"
	"scrapper/model"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/PuerkitoBio/goquery"
)

const ComicExtraURL = "http://www.comicextra.com/"
const ComicExtraURLParam = "ce"

func GetAllComicsComicExtra(doc *goquery.Document) []model.Comic {
	var comics []model.Comic
	doc.Find(".series-col li").Each(func(index int, item *goquery.Selection) {
		comic := model.Comic{}
		aTag := item.Children()
		comic.Title = aTag.Text()
		comic.Link, _ = aTag.Attr("href")
		comic.Category = item.Parent().SiblingsFiltered("div").Text()

		if _, err := strconv.Atoi(comic.Category); err == nil {
			comic.Category = "#"
		}

		comics = append(comics, comic)
	})
	return comics
}

func GetPopularComicsComicExtra(doc *goquery.Document) []model.PopularComic {
	var comics []model.PopularComic
	doc.Find(".cartoon-box").Each(func(index int, item *goquery.Selection) {
		comic := model.PopularComic{}
		//Gets top level information
		comic.Title = item.Find("h3").Children().Text()
		comic.Link, _ = item.Find("h3").Children().Attr("href")
		comic.Img, _ = item.Find("img").Attr("src")
		count := item.Find(".detail").First().Text()
		split := strings.Split(count, " ")
		val, err := strconv.Atoi(split[0])
		if err == nil {
			comic.IssueCount = val
		}

		comics = append(comics, comic)
	})
	return comics
}

func GetChaptersComicExtra(doc *goquery.Document, r *http.Request, comicName string) []model.Chapter {
	var chapters []model.Chapter
	doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
		chapter := model.Chapter{}
		obj := item.Find("a")
		chapter.Link, _ = obj.Attr("href")
		chapter.ChapterName = obj.Text()
		chapter.ReleaseDate = item.Find("td").Last().Text()
		chapters = append(chapters, chapter)
	})

	pageCount := doc.Find(".general-nav").Children().Length() - 1

	if pageCount > 0 {
		chapterChannels := make(chan []model.Chapter, pageCount)

		go getExtraChapters(pageCount, r, comicName, chapterChannels)

		for i := range chapterChannels {
			chapters = append(chapters, i...)
		}
	}
	return chapters
}

func getExtraChapters(extras int, r *http.Request, comicName string, cc chan []model.Chapter) {
	for i := 2; i <= extras; i++ {

		var (
			doc         *goquery.Document
			err, docErr error
			resp        *http.Response
		)
		url := ComicExtraURL + "comic/" + comicName + "/" + strconv.Itoa(i)

		if appengine.IsDevAppServer() {
			c := appengine.NewContext(r)
			client := urlfetch.Client(c)
			resp, err = client.Get(url)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		} else {
			resp, err = http.Get(url)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}

		if err != nil {
			log.Printf(err.Error())
			return
		}

		if resp.StatusCode != 200 {
			log.Printf(resp.Status)
			return
		}

		if docErr != nil {
			log.Printf(docErr.Error())
			return
		}

		defer resp.Body.Close()

		var chapters []model.Chapter
		doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
			chapter := model.Chapter{}
			obj := item.Find("a")
			chapter.Link, _ = obj.Attr("href")
			chapter.ChapterName = obj.Text()
			chapter.ReleaseDate = item.Find("td").Last().Text()
			chapters = append(chapters, chapter)

		})

		cc <- chapters
	}
	close(cc)
}

func GetChapterImagesComicExtra(doc *goquery.Document, r *http.Request, url string) []string {
	numOfPages := doc.Find(".full-select").First().Children().Length()
	pagesChannels := make(chan string, numOfPages)

	go getComicImageURL(url, r, numOfPages, pagesChannels)

	for i := range pagesChannels {
		urls = append(urls, i)
	}
}
