package imp0rt

import "github.com/Zenika/marcel/api/db/plugins"

func Plugins(inputFile string) error {
	var pluginsList []plugins.Plugin

	return imp0rt(inputFile, &pluginsList, func() error {
		return plugins.UpsertAll(pluginsList)
	})
}