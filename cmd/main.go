package main

import (
	"log"
	"waybackdownloader/cmd/data"
	view_download "waybackdownloader/cmd/views/download"

	"github.com/rivo/tview"
)

var config = data.Config{}

func main() {
	config.Init()
	config.App = tview.NewApplication()

	list := tview.NewList().
		AddItem(data.DOWNLOAD_NEW_VIEW_TEXT, "Provide a link for a website to download its content", '1', nil).
		AddItem(data.RESUME_DOWNLOAD_VIEW_TEXT, "Continue downloading resources from initiated process", '2', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			config.App.Stop()
		})

	list.SetBorder(true).SetTitle("Menu").SetTitleAlign(tview.AlignLeft)

	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		switch mainText {
		case data.DOWNLOAD_NEW_VIEW_TEXT:
			view_download.DownloadNew(config.App)
		case data.RESUME_DOWNLOAD_VIEW_TEXT:
		}
	})

	if err := config.App.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}
