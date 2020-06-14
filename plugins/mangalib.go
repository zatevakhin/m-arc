package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"../core/types"
	"../core/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

// GetPluginName - Plugin name
func GetPluginName() string {
	return "mangalib.me"
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

func (manga *MangaInfo) GetChapters() {

	for i := range manga.Chapters {
		result := getPageByURL(manga.Chapters[i].URL)
		parseMangaChapter(&manga.Chapters[i], result)
		break
	}
}

// GetData - Returns all data
func (manga *MangaInfo) GetData() types.MangaInfo {
	return types.MangaInfo(*manga)
}

func getPageByURL(url string) *http.Response {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))

	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	client := &http.Client{Timeout: time.Second * 10}

	result, err := client.Do(req)

	if nil != err {
		logrus.Fatal(err)
	}

	if 200 != result.StatusCode {
		logrus.WithFields(logrus.Fields{
			"url":  url,
			"code": result.StatusCode,
		}).Fatal("Can't get page.")
	}

	return result
}

func parseMangaMetadata(manga *MangaInfo, response *http.Response) {
	document, err := goquery.NewDocumentFromResponse(response)

	if nil != err {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse(manga.URL)

	if err != nil {
		logrus.Fatalf("The URL wasn't parsed!")
	}

	manga.Hostname = parsedURL.Hostname()

	page := document.Find("div.page")

	manga.Title = page.Find("div.manga-title h1").Text()
	manga.TitleAlternative = page.Find("div.manga-title small").Text()
	manga.MangaStatus = types.MS_UNKNOWN

	manga.Description = page.Find("div.info-desc__content").Text()

	mangaInfo := page.Find("div.info-list.manga-info")

	mangaInfo.Find("div.info-list__row").Each(func(i int, s *goquery.Selection) {
		key := strings.ToLower(s.Find("strong").Text())

		switch key {
		case "тип":
			parseMangaType(manga, s)
			break
		case "автор":
			parseMangaAuthor(manga, s)
			break
		case "художник":
			parseMangaArtist(manga, s)
			break
		case "перевод":
			parseMangaTranslationStatus(manga, s)
			break
		case "рейтинг":
			parseMangaRating(manga, s)
			break
		case "просмотров":
			break
		case "дата релиза":
			parseMangaReleaseDate(manga, s)
			break
		case "формат выпуска":
			parseMangaFormat(manga, s)
			break
		case "жанры":
			parseMangaGenres(manga, s)
			break
		case "переводчики":
			parseMangaTranslators(manga, s)
			break
		case "издатель":
			parseMangaPublisher(manga, s)
			break
		}
	})

	chapters := document.Find("div.chapters-list")

	chapterItems := chapters.Find("div.chapter-item")
	chapterItems.Each(func(i int, s *goquery.Selection) {
		link := s.Find("div.chapter-item__name a")

		// dataID, _ := s.Attr("data-id")
		dataIndex, _ := s.Attr("data-index")
		dataVolume, _ := s.Attr("data-volume")
		dataNumber, _ := s.Attr("data-number")
		chapterTitle, _ := link.Attr("title")
		chapterURL, _ := link.Attr("href")

		var chapter types.MangaChapter

		index, err := strconv.Atoi(dataIndex)

		if err == nil {
			chapter.Index = index
		} else {
			logrus.Warnf("chapter.Index value (%s) was not parsed!", dataIndex)
		}

		chapter.VolumeNumber = dataVolume
		chapter.ChapterNumber = dataNumber
		chapter.IsSpecial = strings.Contains(dataNumber, ".")
		chapter.Title = chapterTitle

		chapter.URL = chapterURL

		manga.Chapters = append(manga.Chapters, chapter)

		sort.Slice(manga.Chapters, func(i, j int) bool {
			return manga.Chapters[i].Index < manga.Chapters[j].Index
		})
	})
}

func parseMangaType(manga *MangaInfo, s *goquery.Selection) {
	t := strings.ToLower(s.Find("span").Text())

	switch t {
	case "манга":
		manga.Type = types.MT_MANGA_JP
		break
	case "руманга":
		manga.Type = types.MT_MANGA_RU
		break
	case "манхва":
		manga.Type = types.MT_MANGA_KOR
		break
	case "маньхуа":
		manga.Type = types.MT_MANGA_CN
		break
	case "oel-манга":
		manga.Type = types.MT_OEL_MANGA
		break
	case "комикс западный":
		manga.Type = types.MT_COMIX
		break
	default:
		manga.Type = types.MT_UNKNOWN
		break
	}
}

func parseMangaAuthor(manga *MangaInfo, s *goquery.Selection) {
	name := strings.ToLower(s.Find("a").Text())
	manga.Authors = append(manga.Authors, name)
}

func parseMangaArtist(manga *MangaInfo, s *goquery.Selection) {
	name := strings.ToLower(s.Find("a").Text())
	manga.Artists = append(manga.Artists, name)
}

func parseMangaTranslationStatus(manga *MangaInfo, s *goquery.Selection) {
	status := strings.ToLower(s.Find("span").Text())

	switch status {
	case "продолжается":
		manga.Translation = types.TS_CONTINUES
		break
	case "завершен":
		manga.Translation = types.TS_FINISHED
		break
	case "заморожен":
		manga.Translation = types.TS_FREEZED
		break
	default:
		manga.Translation = types.TS_UNKNOWN
		break
	}
}

func parseMangaRating(manga *MangaInfo, s *goquery.Selection) {
	rating := strings.ToLower(s.Find("span").Text())

	switch rating {
	case "16+":
	case "18+":
		manga.Rating = types.MR_NC17_R
		break
	default:
		manga.Rating = types.MR_UNKNOWN
		break
	}
}

func parseMangaReleaseDate(manga *MangaInfo, s *goquery.Selection) {
	releaseDate := strings.Trim(strings.ToLower(s.Find("span").Text()), " ")
	year, err := strconv.Atoi(releaseDate)

	if err == nil {
		manga.Year = year
	} else {
		logrus.Warnf("Year value (%s) was not parsed!", releaseDate)
	}
}

func parseMangaFormat(manga *MangaInfo, s *goquery.Selection) {
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		format := strings.ToLower(s.Find("span").Text())
		switch format {
		case "4-кома (ёнкома)":
			manga.Format = append(manga.Format, types.MF_4_COMA)
			break
		case "сборник":
			manga.Format = append(manga.Format, types.MF_COMPILATION)
			break
		case "додзинси":
			manga.Format = append(manga.Format, types.MF_DOJINSHI)
			break
		case "в цвете":
			manga.Format = append(manga.Format, types.MF_IN_COLOR)
			break
		case "сингл":
			manga.Format = append(manga.Format, types.MF_SINGLE)
			break
		case "веб":
			manga.Format = append(manga.Format, types.MF_WEB)
			break
		default:
			manga.Format = append(manga.Format, types.MF_UNKNOWN)
			break
		}
	})

}

func parseMangaGenres(manga *MangaInfo, s *goquery.Selection) {
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		manga.Genres = append(manga.Genres, strings.Trim(s.Text(), " "))
	})
}

func parseMangaTranslators(manga *MangaInfo, s *goquery.Selection) {
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		manga.Translators = append(manga.Translators, strings.Trim(s.Text(), " "))
	})
}

func parseMangaPublisher(manga *MangaInfo, s *goquery.Selection) {
	manga.Publisher = strings.Trim(strings.ToLower(s.Find("span").Text()), " ")
}

func parseMangaChapter(chapter *types.MangaChapter, response *http.Response) {
	document, err := goquery.NewDocumentFromResponse(response)

	if nil != err {
		log.Fatal(err, document)
	}

	pagesRegex := regexp.MustCompile(`window.__pg = (?P<pages>[^;]+)`)
	infoRegex := regexp.MustCompile(`window.__info = (?P<info>[^;]+)`)

	pagesScript := document.Find("script#pg").Text()

	pages := utils.FindNamedMatches(pagesRegex, pagesScript)["pages"]

	var info string

	document.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		infoScript := s.Text()

		if strings.Contains(infoScript, "window.__info") {
			info = utils.FindNamedMatches(infoRegex, infoScript)["info"]
			return false
		}
		return true
	})

	var objInfo map[string]interface{}
	var objPages []interface{}

	err = json.Unmarshal([]byte(info), &objInfo)

	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(pages), &objPages)

	if err != nil {
		fmt.Println(err)
	}

	img := objInfo["img"]

	if item, ok := img.(map[string]interface{}); ok {
		chapterURL := item["url"].(string)
		serverID := item["server"].(string)

		serversObj := objInfo["servers"]
		if servers, ok := serversObj.(map[string]interface{}); ok {
			serversURL := servers[serverID]

			for _, v := range objPages {
				if item, ok := v.(map[string]interface{}); ok {
					imageURL := item["u"].(string)
					ImageIndex := int(item["p"].(float64))

					URL := fmt.Sprintf("%s%s%s", serversURL, chapterURL, imageURL)
					chapter.Pages = append(chapter.Pages, types.MangaPage{ImageURL: URL, ImageIndex: ImageIndex})
				}
			}
		}
	}
}
