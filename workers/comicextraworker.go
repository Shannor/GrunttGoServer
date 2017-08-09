package workers

import (
	"net/http"
	"scrapper/utils"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//GetTag bool check to see if provided tag matches this struct
func (ce *ComicExtra) GetTag() string {
	return ce.Tag
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

//GetComicChapterListPageAmount returns the number of pages the chapter list spans
func (ce *ComicExtra) GetComicChapterListPageAmount(doc *goquery.Document) (int, error) {
	pageCount := doc.Find(".general-nav").Children().Length() - 1
	return pageCount, nil
}

type chaptersResult struct {
	Chapters Chapters
	err      error
}

//GetComicChapterList return the list of all the chapters of a comic
func (ce *ComicExtra) GetComicChapterList(comicName string, numOfPages int, r *http.Request) (Chapters, error) {
	var chapters Chapters
	chapterChannels := make(chan Chapters, numOfPages)
	baseURL := ce.BaseURL + "comic/" + comicName
	go getExtraChapters(numOfPages, r, baseURL, chapterChannels)

	for i := range chapterChannels {
		chapters = append(chapters, i...)
	}
	return chapters, nil
}

func getExtraChapters(pageAmount int, r *http.Request, baseURL string, cc chan Chapters) {
	for i := 1; i <= pageAmount; i++ {

		url := baseURL + "/" + strconv.Itoa(i)
		doc, err := utils.GetGoQueryDoc(url, r)
		if err != nil {
			return
		}
		var chapters Chapters
		doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
			chapter := Chapter{}
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
