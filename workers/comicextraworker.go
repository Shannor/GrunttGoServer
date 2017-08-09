package workers

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//MatchesTag bool check to see if provided tag matches this struct
func (ce *ComicExtra) MatchesTag(tag string) bool {
	if tag == ce.Tag {
		return true
	}
	return false
}

//CreateAllComicsURL makes the url for ComicExtra's All comics page
func (ce *ComicExtra) CreateAllComicsURL() string {
	return ce.BaseURL + "comic-list"
}

//CreatePopularComicsURL makes the url for ComicExtra's Popular comics page
func (ce *ComicExtra) CreatePopularComicsURL(pageNumber int) string {
	return ce.BaseURL + "popular-comic/" + strconv.Itoa(pageNumber)
}

//CreateComicChapterListURL makes the url for ComicExtra's Chapter list page
func (ce *ComicExtra) CreateComicChapterListURL(comicName string) string {
	return ce.BaseURL + "comic/" + comicName
}

//CreateChapterPagesURL makes the url for ComicExtra's chapter pages
func (ce *ComicExtra) CreateChapterPagesURL(comicName string, chapterNumber int) string {
	return ce.BaseURL + comicName + "/chapter-" + strconv.Itoa(chapterNumber)
}

//CreateComicDescriptionURL makes the url for ComicExtra's Description comic page
func (ce *ComicExtra) CreateComicDescriptionURL(comicName string) string {
	return ce.BaseURL + "comic/" + comicName
}

//GetAllComics performs the webcrawling to return all comics on the site
func (ce *ComicExtra) GetAllComics(doc *goquery.Document) (Comics, error) {
	var comics Comics
	doc.Find(".series-col li").Each(func(index int, item *goquery.Selection) {
		comic := Comic{}
		aTag := item.Children()
		comic.Title = aTag.Text()
		comic.Link, _ = aTag.Attr("href")
		comic.Category = item.Parent().SiblingsFiltered("div").Text()

		if _, err := strconv.Atoi(comic.Category); err == nil {
			comic.Category = "#"
		}

		comics = append(comics, comic)
	})
	return comics, nil
}

//GetPopularComics returns all the popular comics from ComicExtra website
func (ce *ComicExtra) GetPopularComics(doc *goquery.Document) (PopularComics, error) {
	var comics PopularComics
	doc.Find(".cartoon-box").Each(func(index int, item *goquery.Selection) {
		comic := PopularComic{}
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
	return comics, nil
}

//
// func GetChaptersComicExtra(doc *goquery.Document, r *http.Request, comicName string) []Chapter {
// 	var chapters []Chapter
// 	doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
// 		chapter := Chapter{}
// 		obj := item.Find("a")
// 		chapter.Link, _ = obj.Attr("href")
// 		chapter.ChapterName = obj.Text()
// 		chapter.ReleaseDate = item.Find("td").Last().Text()
// 		chapters = append(chapters, chapter)
// 	})
//
// 	pageCount := doc.Find(".general-nav").Children().Length() - 1
//
// 	if pageCount > 0 {
// 		chapterChannels := make(chan []Chapter, pageCount)
//
// 		go getExtraChapters(pageCount, r, comicName, chapterChannels)
//
// 		for i := range chapterChannels {
// 			chapters = append(chapters, i...)
// 		}
// 	}
// 	return chapters
// }
//
// func getExtraChapters(extras int, r *http.Request, comicName string, cc chan []Chapter) {
// 	for i := 2; i <= extras; i++ {
//
// 		var (
// 			doc         *goquery.Document
// 			err, docErr error
// 			resp        *http.Response
// 		)
// 		url := ComicExtraURL + "comic/" + comicName + "/" + strconv.Itoa(i)
//
// 		if appengine.IsDevAppServer() {
// 			c := appengine.NewContext(r)
// 			client := urlfetch.Client(c)
// 			resp, err = client.Get(url)
// 			doc, docErr = goquery.NewDocumentFromResponse(resp)
// 		} else {
// 			resp, err = http.Get(url)
// 			doc, docErr = goquery.NewDocumentFromResponse(resp)
// 		}
//
// 		if err != nil {
// 			log.Printf(err.Error())
// 			return
// 		}
//
// 		if resp.StatusCode != 200 {
// 			log.Printf(resp.Status)
// 			return
// 		}
//
// 		if docErr != nil {
// 			log.Printf(docErr.Error())
// 			return
// 		}
//
// 		defer resp.Body.Close()
//
// 		var chapters []Chapter
// 		doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
// 			chapter := Chapter{}
// 			obj := item.Find("a")
// 			chapter.Link, _ = obj.Attr("href")
// 			chapter.ChapterName = obj.Text()
// 			chapter.ReleaseDate = item.Find("td").Last().Text()
// 			chapters = append(chapters, chapter)
//
// 		})
//
// 		cc <- chapters
// 	}
// 	close(cc)
// }
//
// func GetChapterImagesComicExtra(doc *goquery.Document, r *http.Request, url string) []string {
// 	numOfPages := doc.Find(".full-select").First().Children().Length()
// 	pagesChannels := make(chan string, numOfPages)
// 	urls := make([]string, 0)
//
// 	go getComicImageURL(url, r, numOfPages, pagesChannels)
//
// 	for i := range pagesChannels {
// 		urls = append(urls, i)
// 	}
// 	return urls
// }
//
// func getComicImageURL(url string, r *http.Request, numOfPages int, cc chan string) {
// 	for i := 1; i <= numOfPages; i++ {
// 		pageURL := url + "/" + strconv.Itoa(i)
//
// 		doc, err := utils.GetGoQueryDoc(pageURL, r)
// 		if err != nil {
// 			return
// 		}
// 		link, _ := doc.Find("#main_img").Attr("src")
// 		cc <- link
// 	}
// 	close(cc)
// }
