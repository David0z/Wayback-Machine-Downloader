package views

import (
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"

	"github.com/rivo/tview"
)

func Options_View(config *config.Config) {
	form := tview.NewForm()

	form.AddCheckbox("Copy full file path", (*config.Options.Options)[data.OptionsMap[data.OPTION_COPY_FULL_PATH]], func(checked bool) {
		config.SetOption(data.OptionsMap[data.OPTION_COPY_FULL_PATH], checked)
	})

	form.AddButton("Back", func() {
		config.App.SetRoot(MainMenuView(config), true)
	})

	form.SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	config.App.SetRoot(form, true)
}
