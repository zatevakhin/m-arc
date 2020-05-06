package main

import (
	"fmt"

	"./core"
	"./core/database"
)

// http://fanfox.net/
// https://mangalib.me
// https://mintmanga.live
// https://readmanga.me/
// https://selfmanga.ru/
// https://mangapark.net/
// https://h-chan.me/
// https://manga-chan.me/

func main() {
	var pluginManager core.PluginManager

	pluginManager.SetPluginsFolder("plugins")

	// pluginInstance := pluginManager.GetPluginForSite("https://readmanga.me/angel_beats__heaven_s_door__A1b916d")
	// pluginInstance := pluginManager.GetPluginForSite("https://mintmanga.live/gate___thus_the_jsdf_fought_there")
	pluginInstance := pluginManager.GetPluginForSite("http://fanfox.net/manga/citrus_saburouta/")
	// pluginInstance := pluginManager.GetPluginForSite("https://mangalib.me/octave")

	if nil != pluginInstance {
		(*pluginInstance).GetMetaData()
		(*pluginInstance).GetChapters()
		fmt.Print(*pluginInstance)
	}

	db := database.GetInitializedDatabase()
	defer db.Close()

	// m := macaron.Classic()

	// db := getInitializedDatabase()
	// defer db.Close()

	// m.Get("/", handlers.IndexPage)

	// log.Println("Server is running...")
	// log.Println(http.ListenAndServe("0.0.0.0:4000", m))
}

// package main

// import (
// 	"fmt"

// 	"./core"
// )

// func main() {

// 	var pluginManager core.PluginManager

// 	pluginManager.SetPluginsFolder("plugins")

// 	pluginInstance := *pluginManager.GetPlugin("http://mintmanga.com/chto_vy_zdes_delaete__sensei__")
// 	pluginInstance.SetUrl("http://mintmanga.com/chto_vy_zdes_delaete__sensei__")
// 	pluginInstance.GetMetaData()
// 	pluginInstance.GetChapters()
// 	fmt.Print(pluginInstance)

// 	return

// 	// plug, err := plugin.Open("plugins/mintmanga.so")

// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	os.Exit(1)
// 	// }

// 	// symbolMangaInfoExport, err := plug.Lookup("MangaInfoExport")

// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	os.Exit(1)
// 	// }

// 	// var manga types.AbstractPlugin
// 	// manga, ok := symbolMangaInfoExport.(types.AbstractPlugin)
// 	// if !ok {
// 	// 	fmt.Println("unexpected type from module symbol")
// 	// 	os.Exit(1)
// 	// }

// 	// manga.SetUrl("http://readmanga.me/hito")
// 	// manga.SetUrl("http://fanfox.net/")
// 	// manga.SetUrl("https://mangapark.net/")
// 	// manga.SetUrl("http://readmanga.me/battle_angel_alita")
// 	// manga.SetUrl("http://mintmanga.com/chto_vy_zdes_delaete__sensei__")
// 	// manga.GetMetaData()
// 	// manga.GetChapters()
// 	// fmt.Print(manga)

// }

// // var info plugins.MangaInfo

// // info.URL = "http://readmanga.me/battle_angel_alita"

// // info.GetChapters()
// // https://github.com/vladimirvivien/go-plugin-example
// // https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9

// // func (url *AbstractPluginInfo) getChapters() {

// // }

// // func (url *AbstractPluginInfo) getMetaData() {

// // }

// // func main() {
// // var manga MangaInfo

// // manga.url = "http://readmanga.me/battle_angel_alita"

// // var murl AbstractPluginUrl

// // murl = manga

// // }

// // func parseMangaPage(response *http.Response) {
// // 	document, err := goquery.NewDocumentFromResponse(response)

// // 	if nil != err {
// // 		log.Fatal(err)
// // 	}

// // 	var info MangaInfo

// // 	meta := document.Find("div.subject-meta")

// // 	meta.Find("span.elem_genre").Each(func(i int, s *goquery.Selection) {
// // 		info.genres = append(info.genres, strings.Trim(s.Text(), " ,"))
// // 	})

// // 	meta.Find("span.elem_author").Each(func(i int, s *goquery.Selection) {
// // 		info.authors = append(info.authors, strings.Trim(s.Text(), " ,"))
// // 	})

// // 	meta.Find("span.elem_translator").Each(func(i int, s *goquery.Selection) {
// // 		info.translators = append(info.translators, strings.Trim(s.Text(), " ,"))
// // 	})

// // 	meta.Find("span.elem_tag").Each(func(i int, s *goquery.Selection) {
// // 		info.categories = append(info.categories, strings.Trim(s.Text(), " ,"))
// // 	})

// // 	meta.Find("span.elem_limitation").Each(func(i int, s *goquery.Selection) {
// // 		info.rating = append(info.rating, strings.Trim(s.Text(), " ,"))
// // 	})

// // 	yearSelection := meta.Find("span.elem_year")
// // 	if nil != yearSelection {
// // 		info.year = yearSelection.Text()
// // 	}

// // 	fmt.Print(info)
// // }

// // func downloadByURL(url string) {

// // 	result, err := http.Get(url)

// // 	if nil != err {
// // 		logrus.Fatal(err)
// // 	}

// // 	if 200 != result.StatusCode {
// // 		logrus.WithFields(logrus.Fields{
// // 			"url":  url,
// // 			"code": result.StatusCode,
// // 		}).Fatal("Can't get page.")
// // 	}

// // 	logrus.Debug(result)

// // 	parseMangaPage(result)
// // }
// // downloadByURL("http://readmanga.me/battle_angel_alita")

// // package main

// // import (
// //     // "net/http"
// //     // "github.com/go-macaron/session"
// //     // "gopkg.in/macaron.v1"
// // )

// // import (
// //     log "github.com/sirupsen/logrus"
// // )

// // func init() {
// //     // Only log the warning severity or above.
// //     log.SetLevel(log.DebugLevel)
// // }

// // func main() {

// //     // m := macaron.New()
// //     // m.Use(session.Sessioner(session.Options{
// //     //     CookieName: "dfsid",
// //     //     IDLength:   32,
// //     // }))

// //     // m.Use(macaron.Static("www/static"))

// //     // //m.Get("/", indexGetHandle)
// //     // m.Get("/login", loginGetHandle)

// //     // m.Get("/", func(sess session.Store) string {
// //     //     sess.Set("session", "session middleware")
// //     //     return sess.Get("session").(string)
// //     // })

// //     // log.Debug("Server is running...")
// //     // http.ListenAndServe("0.0.0.0:4000", m)
// // }

// // // func loginGetHandle(ctx *macaron.Context) string {
// // //     return "the request path is: " + ctx.Req.RequestURI
// // // }

// // // func indexGetHandle(ctx *macaron.Context) string {
// // //     return "the request path is: " + ctx.Req.RequestURI
// // // }
