package core

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"plugin"

	"./types"
	"github.com/sirupsen/logrus"
)

type PluginManager struct {
	PluginsFolder string
}

func getPluginNameFromURL(urlString string) string {
	u, err := url.Parse(urlString)

	if err != nil {
		logrus.Fatalf("The URL wasn't parsed!")
	}

	return u.Hostname()
}

// GetPluginForSite - sss
func (pm *PluginManager) GetPluginForSite(urlString string) *types.AbstractPlugin {
	siteName := getPluginNameFromURL(urlString)
	pluginItem := pm.loadPlugin(siteName)

	if nil != pluginItem {
		(*pluginItem).SetUrl(urlString)
	}

	return pluginItem
}

func (pm *PluginManager) loadPlugin(siteName string) *types.AbstractPlugin {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.so", pm.PluginsFolder))

	if err != nil {
		log.Fatal(err)
	}

	var sitePlugin types.AbstractPlugin

	for _, file := range files {
		plug, err := plugin.Open(file)

		if err != nil {
			log.Println(err)
			continue
		}

		name, err := plug.Lookup("GetPluginName")
		if err != nil {
			log.Println(name, err)
			continue
		}

		version, err := plug.Lookup("GetPluginVersion")
		if err != nil {
			log.Println(version, err)
			continue
		}

		GetPluginName, ok := name.(func() string)
		if !ok {
			log.Println(file, "have no function GetPluginName")
			continue
		}

		GetPluginVersion, ok := version.(func() string)
		if !ok {
			log.Println(file, "have no function GetPluginVersion")
			continue
		}

		mangaInfoPlugin, err := plug.Lookup("MangaInfoExport")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if GetPluginName() != siteName {
			// fmt.Println(GetPluginName(), "!=", siteName)
			continue
		}

		if !checkVersion(GetPluginVersion()) {
			fmt.Println("checkVersion(GetPluginVersion())")
			continue
		}

		sitePlugin, ok = mangaInfoPlugin.(types.AbstractPlugin)

		if !ok {
			fmt.Println("unexpected type from module symbol")
			continue
		}

		break
	}

	if sitePlugin == nil {
		fmt.Println("Plugin not found for", siteName)
		return nil
	}

	return &sitePlugin
}

func checkVersion(version string) bool {
	return true
}
