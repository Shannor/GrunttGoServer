package workers

import (
	"scrapper/model"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	ReadComicsURL      = "http://readcomics.website/"
	ReadComicsURLParam = "rcw"
)

func GetAllComicsReadComics(doc *goquery.Document) []model.Comic {
	var comics []model.Comic
	doc.Find(".type-content li").Each(func(index int, item *goquery.Selection) {
		comic := model.Comic{}
		aTag := item.ChildrenFiltered("a")
		comic.Title = strings.TrimSpace(aTag.Text())
		comic.Link, _ = aTag.Attr("href")

		if _, err := strconv.Atoi(string(comic.Title[0])); err == nil {
			comic.Category = "#"
		} else {
			comic.Category = string(comic.Title[0])
		}
		comics = append(comics, comic)
	})
	return comics
}

func GetPopularComicsReadComics(doc *goquery.Document) []model.PopularComic {
	var comics []model.PopularComic
	doc.Find(".media").Each(func(index int, item *goquery.Selection) {
		comic := model.PopularComic{}
		//Gets top level information
		obj := item.Find(".chart-title")
		comic.Title = obj.Text()
		comic.Link, _ = obj.Attr("href")
		comic.Img, _ = item.Find("img").Attr("src")
		count := item.Find("i").Parent().Text()
		count = strings.TrimSpace(count)
		val, err := strconv.Atoi(count)
		if err == nil {
			comic.IssueCount = val
		}

		comics = append(comics, comic)
	})
	return comics
}

func GetChaptersReadComics(doc *goquery.Document) []model.Chapter {
	var chapters []model.Chapter
	doc.Find(".chapters").Children().Each(func(index int, item *goquery.Selection) {
		chapter := model.Chapter{}
		obj := item.Find("a")
		chapter.Link, _ = obj.Attr("href")
		chapter.ChapterName = obj.Text()
		chapter.ReleaseDate = strings.TrimSpace(item.Find(".date-chapter-title-rtl").Text())
		chapters = append(chapters, chapter)
	})

	return chapters
}
