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
func (ce *ComicExtra) GetComicChapterListPageAmount(doc *goquery.Document) int {
	numOfPages := doc.Find(".general-nav").Children().Length()
	if numOfPages > 0 {
		numOfPages--
	} else {
		numOfPages = 1
	}
	return numOfPages
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

//GetNumberOfPages method to return the number of pages a chapter has
func (ce *ComicExtra) GetNumberOfPages(doc *goquery.Document) int {
	return doc.Find(".full-select").Last().Children().Length()
}

//GetChapterPages method to return the url for all pages
func (ce *ComicExtra) GetChapterPages(comicName string, chapterNumber int, numOfPages int, r *http.Request) ([]string, error) {
	baseURL := ce.BaseURL + comicName + "/chapter-" + strconv.Itoa(chapterNumber)
	pagesChannels := make(chan string, numOfPages)
	var urls []string
	go ce.getComicImageURL(baseURL, numOfPages, r, pagesChannels)

	for url := range pagesChannels {
		urls = append(urls, url)
	}
	return urls, nil
}

//GetComicImageURL go routine to get the urls for the images
func (ce *ComicExtra) getComicImageURL(url string, numOfPages int, r *http.Request, cc chan string) {
	for i := 1; i <= numOfPages; i++ {
		pageURL := url + "/" + strconv.Itoa(i)

		doc, err := utils.GetGoQueryDoc(pageURL, r)
		if err != nil {
			cc <- ""
		}
		link, _ := doc.Find("#main_img").Attr("src")
		cc <- strings.TrimSpace(link)
	}
	close(cc)
}
