package main

import (
	"log"
	"waybackdownloader/cmd/data/config"
	"waybackdownloader/cmd/views"

	_ "github.com/glebarez/go-sqlite"
	"github.com/rivo/tview"
)

func main() {
	var config = config.Config{}
	config.Init()
	config.App = tview.NewApplication()

	list := views.MainMenuView(&config)

	if err := config.App.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}
