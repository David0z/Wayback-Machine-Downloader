package views

import (
	"waybackdownloader/cmd/data/config"

	"github.com/rivo/tview"
)

func Options_View(config *config.Config) {
	form := tview.NewForm().
		AddButton("Back", func() {
			config.App.SetRoot(MainMenuView(config), true)
		}).
		AddButton("Quit", func() {
			config.App.Stop()
		})

	form.SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	config.App.SetRoot(form, true)
}
