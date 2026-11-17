package main

import (
	"embed"

	"github.com/egoist/mygo/ui"
)

//go:embed icons/*.svg
var iconFiles embed.FS

func mustIcon(name string) *ui.SVG {
	data, err := iconFiles.ReadFile("icons/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	return ui.MustParseSVG(data)
}

var (
	homeIcon      = mustIcon("house")
	chartIcon     = mustIcon("chart-no-axes-column")
	cartIcon      = mustIcon("shopping-cart")
	usersIcon     = mustIcon("users")
	settingsIcon  = mustIcon("settings")
	databaseIcon  = mustIcon("database")
	percentIcon   = mustIcon("percent")
	bellIcon      = mustIcon("bell")
	sunIcon       = mustIcon("sun")
	moonIcon      = mustIcon("moon")
	chevronIcon   = mustIcon("chevron-down")
	arrowUpIcon   = mustIcon("arrow-up")
	arrowDownIcon = mustIcon("arrow-down")
	arrowIcon     = mustIcon("arrow-right")
	searchIcon    = mustIcon("search")
)
