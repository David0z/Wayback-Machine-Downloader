package views

import (
	"waybackdownloader/cmd/data/config"
	"waybackdownloader/cmd/repository/db"

	"github.com/rivo/tview"
)

func Options_View(config *config.Config) {
	form := tview.NewForm()

	form.AddCheckbox("Copy full file path", (*config.Options.Options)[db.OptionsMap[db.OPTION_COPY_FULL_PATH]], func(checked bool) {
		config.SetOption(db.OptionsMap[db.OPTION_COPY_FULL_PATH], checked)
	})

	form.AddButton("Back", func() {
		config.App.SetRoot(MainMenuView(config), true)
	})

	form.SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	config.App.SetRoot(form, true)
}
