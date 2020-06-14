package main

import (
	"log"
	"net/http"

	"./core"
	"./core/database"
	"./handlers"
	"gopkg.in/macaron.v1"
)

// [-] http://fanfox.net/

// https://mangalib.me
// https://mintmanga.live
// https://readmanga.me/
// https://selfmanga.ru/
// https://mangapark.net/
// https://h-chan.me/
// https://manga-chan.me/

func main() {
	db := database.GetInitializedDatabase()
	defer db.Close()

	pluginManager := core.PluginManager{
		PluginsFolder: "plugins/",
	}

	downloadManager := core.DownloadManager{
		MangaFolder:   "manga/",
		DataBase:      db,
		PluginManager: &pluginManager,
	}

	m := macaron.Classic()

	m.Map(&db)
	m.Map(&downloadManager)

	m.Use(macaron.Renderers(macaron.RenderOptions{
		Directory: "www/templates"}))

	m.Get("/", handlers.IndexPage)
	m.Get("/upload", handlers.UploadPage)

	log.Println("Server is running...")
	log.Println(http.ListenAndServe("0.0.0.0:4000", m))
}
