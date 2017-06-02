package workers

import (
	"scrapper/model"
	"strconv"
	"strings"

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
