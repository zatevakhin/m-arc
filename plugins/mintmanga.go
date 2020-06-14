package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"../core/types"
	"github.com/PuerkitoBio/goquery"
	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
)

// GetPluginName - Plugin name
func GetPluginName() string {
	return "mintmanga.live"
}

// GetPluginVersion - Version number
func GetPluginVersion() string {
	return "0.0.1"
}

// MangaInfo - Object contains data about manga
type MangaInfo types.MangaInfo

// MangaInfoExport - exported
var MangaInfoExport MangaInfo

func (manga *MangaInfo) SetUrl(url string) {
	manga.URL = url
}

func (manga *MangaInfo) GetMetaData() {
	result := getPageByURL(manga.URL)
	parseMangaMetadata(manga, result)
}

// GetData - Returns all data
func (manga *MangaInfo) GetData() types.MangaInfo {
	return types.MangaInfo(*manga)
}

func getTranslationStatus(status *string) types.TranslationStatus {
	var statusCode types.TranslationStatus
	switch *status {
	case "завершен":
		statusCode = types.TS_FINISHED
	case "продолжается":
		statusCode = types.TS_CONTINUES
	default:
		statusCode = types.TS_UNKNOWN
	}
	return statusCode
}

func getMangaStatus(status *string) types.MangaCompletionStatus {
	var statusCode types.MangaCompletionStatus
	switch *status {
	case "":
		statusCode = types.MS_FINISHED
	case "выпуск продолжается":
		statusCode = types.MS_CONTINUES
	default:
		statusCode = types.MS_UNKNOWN
	}
	return statusCode
}

func getPageByURL(url string) *http.Response {

	result, err := http.Get(string(url))

	if nil != err {
		logrus.Fatal(err)
	}

	if 200 != result.StatusCode {
		logrus.WithFields(logrus.Fields{
			"url":  url,
			"code": result.StatusCode,
		}).Fatal("Can't get page.")
	}

	logrus.Debug(result)

	return result
}

func parseMangaMetadata(manga *MangaInfo, response *http.Response) {
	document, err := goquery.NewDocumentFromResponse(response)

	if nil != err {
		log.Fatal(err)
	}

	leftContent := document.Find("div.leftContent")

	leftContent.Find("meta").Each(func(i int, s *goquery.Selection) {
		itempropValue, itempropExists := s.Attr("itemprop")
		contentValue, contentExists := s.Attr("content")

		if itempropExists && contentExists {
			switch itempropValue {
			case "name":
				manga.Title = contentValue
				break
			case "alternativeHeadline":
				manga.TitleAlternative = contentValue
				break
			case "description":
				manga.Description = contentValue
				break
			}
		}
	})

	covers := leftContent.Find("div.picture-fotorama")

	covers.Find("img").Each(func(i int, s *goquery.Selection) {
		srcValue, srcExists := s.Attr("src")
		// fullValue, fullExists := s.Attr("data-full")

		if srcExists {
			manga.Covers = append(manga.Covers, srcValue)
		}

		// if fullExists {
		// 	manga.CoversBig = append(manga.Covers, fullValue)
		// }
	})

	covers.Find("a").Each(func(i int, s *goquery.Selection) {
		hrefValue, hrefExists := s.Attr("href")
		// fullValue, fullExists := s.Attr("data-full")

		if hrefExists {
			manga.Covers = append(manga.Covers, hrefValue)
		}

		// if fullExists {
		// 	manga.CoversBig = append(manga.Covers, fullValue)
		// }
	})

	subjectMeta := leftContent.Find("div.subject-meta")

	subjectMeta.Find("span.elem_genre").Each(func(i int, s *goquery.Selection) {
		manga.Genres = append(manga.Genres, strings.Trim(s.Text(), " ,"))
	})

	subjectMeta.Find("span.elem_author").Each(func(i int, s *goquery.Selection) {
		manga.Authors = append(manga.Authors, strings.Trim(s.Text(), " ,"))
	})

	subjectMeta.Find("span.elem_translator").Each(func(i int, s *goquery.Selection) {
		manga.Translators = append(manga.Translators, strings.Trim(s.Text(), " ,"))
	})

	// subjectMeta.Find("span.elem_limitation").Each(func(i int, s *goquery.Selection) {
	// 	manga.Rating = append(manga.Rating, strings.Trim(s.Text(), " ,"))
	// })

	yearSelection := subjectMeta.Find("span.elem_year")
	if nil != yearSelection {
		year, err := strconv.Atoi(strings.Trim(yearSelection.Text(), " "))
		if err == nil {
			manga.Year = year
		} else {
			logrus.Warnf("Year value (%s) was not parsed!", yearSelection.Text())
		}
	}

	subjectMeta.Find("p b").Each(func(i int, s *goquery.Selection) {
		parent := s.Parent()
		b := s.Empty()

		bText := strings.ToLower(strings.Trim(b.Text(), "\t\r\n: "))
		pText := strings.ToLower(strings.Trim(parent.Text(), "\t\r\n: "))

		switch bText {
		case "томов":
			x := strings.Split(pText, ",")

			if len(x) == 2 {
				manga.MangaStatus = getMangaStatus(&x[1])
			}

			volumes, err := strconv.Atoi(x[0])
			if err == nil {
				manga.Volumes = volumes
			} else {
				logrus.Warnf("Volumes value (%s) was not parsed!", pText)
			}

			break
		case "перевод":
			manga.Translation = getTranslationStatus(&pText)
			break
		case "сингл":
			manga.IsSinge = true
			break
		default:
			logrus.Warnf("Unknown params (b: %s, p: %s)", bText, pText)
		}

	})

	manga.HasAgeLimit = len(leftContent.Find("div.mature-message").Nodes) > 0 ||
		len(leftContent.Find("div.mtr-message").Nodes) > 0
}

func (manga *MangaInfo) GetChapters() {
	u, err := url.Parse(manga.URL)
	if err != nil {
		logrus.Fatalf("The URL wasn't parsed!")
	}

	path := fmt.Sprintf("%s", regexp.MustCompile(`__[0-9a-fA-F]{7}`).ReplaceAll([]byte(strings.Trim(u.Path, "/")), []byte("")))

	rssURL := fmt.Sprintf("http://%s/rss/manga?name=%s", u.Hostname(), path)

	result := getPageByURL(rssURL)

	doc := etree.NewDocument()
	n, err := doc.ReadFrom(result.Body)

	if err != nil {
		logrus.Fatal(n, err)
	}

	rss := doc.SelectElement("rss")
	channel := rss.SelectElement("channel")
	items := channel.SelectElements("item")

	for i, index := len(items)-1, 0; i >= 0; i-- {
		item := items[i]
		index++

		title := item.SelectElement("title").Text()
		link := item.SelectElement("link").Text()
		appendChapter(manga, title, link, index)
	}
}

func appendChapter(manga *MangaInfo, chapter string, link string, index int) {
	// https://regex-golang.appspot.com/assets/html/index.html

	// defaultRegex := regexp.MustCompile(`^(?:.+): (?P<volume>\d+) [-] (?P<chapter>\d+)(?P<name>.+)$`)
	// extraRegex := regexp.MustCompile(`^(?:.+): (?P<volume>\d+) [^-](?P<name>.+)$`)
	// singleRegex := regexp.MustCompile(`^(?:.+): (?P<name>[^-\d]+)$`)

	// result := defaultRegex.FindStringSubmatch(chapter)

	// if len(result) == 0 {
	// 	result = extraRegex.FindStringSubmatch(chapter)
	// 	if len(result) == 0 {
	// 		result = singleRegex.FindStringSubmatch(chapter)
	// 	}
	// }

	// var chapter types.MangaChapter

	// chapter.Index = index
	// chapter.VolumeNumber = "volumeNumber"
	// chapter.ChapterNumber = "chapterNumber"
	// chapter.IsSpecial = isSpecial
	// chapter.Name = strings.Trim(chapterName, " ")
	// chapter.URL = link

	// manga.Chapters = append(manga.Chapters, chapter)

}
