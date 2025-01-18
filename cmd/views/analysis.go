package views

import (
	"fmt"
	"os"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"

	"github.com/rivo/tview"
)

func Analysis_View(config *config.Config, websiteURL string) {
	textView := tview.NewTextView().
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Loading website URLs...")

	config.App.SetRoot(textView, true)

	go func() {
		// get data from links.json

		form := tview.NewForm().
			AddCheckbox("Age 18+", false, nil).
			AddButton("Save", nil).
			AddButton("Download", func() {
				config.App.Stop()
			})
		form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

		flex := tview.NewFlex().
			AddItem(form, 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 0, 1, false)

		config.App.QueueUpdateDraw(func() {
			config.App.SetRoot(flex, true)
		})
	}()
}

func AnalysisList_View(config *config.Config) {
	textView := tview.NewTextView().
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Loading user's URLs...")

	config.App.SetRoot(textView, true)

	go func() {
		baseFolder := data.MAIN_PATH

		entries, err := os.ReadDir(baseFolder)
		if err != nil {
			panic(fmt.Sprintf("Error reading directory: %v\n", err))
		}

		folderNameList := []string{}

		for _, entry := range entries {
			if entry.IsDir() {
				existingLink, err := config.DB.SelectOne(entry.Name())
				if err != nil {
					panic(err)
				}

				if existingLink {
					folderNameList = append(folderNameList, entry.Name())
				}
			}
		}

		if len(folderNameList) == 0 {
			modal := tview.NewModal().
				SetText("No website URLs found.").
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					config.App.SetRoot(MainMenuView(config), true)
				})

			config.App.QueueUpdateDraw(func() {
				config.App.SetRoot(modal, true)
			})
			return
		}

		list := tview.NewList()

		for index, folderName := range folderNameList {
			list.AddItem(folderName, "", rune(index), nil)
		}

		list.SetBorder(true).SetTitle("Select website URL").SetTitleAlign(tview.AlignLeft)

		list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			Analysis_View(config, mainText)
		})

		config.App.QueueUpdateDraw(func() {
			config.App.SetRoot(list, true)
		})
	}()
}
