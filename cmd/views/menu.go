package views

import (
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"

	"github.com/rivo/tview"
)

func MainMenuView(config *config.Config) *tview.List {
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
			DownloadWebsiteURLs_View(config)
		case data.RESUME_DOWNLOAD_VIEW_TEXT:
			AnalysisList_View(config)
		}
	})

	return list
}
