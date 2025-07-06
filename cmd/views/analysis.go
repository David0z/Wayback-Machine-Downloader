package views

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"
	"waybackdownloader/cmd/repository/db"
	"waybackdownloader/cmd/web/api"

	"github.com/gdamore/tcell/v2"
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
		mimetypeQuantity, err := config.DB.MimetypeQuantity(websiteURL)
		if err != nil {
			panic(err)
		}

		form := tview.NewForm()

		maxLength := 0
		maxRowCountLength := 0
		maxDownloadLength := 0

		for _, mimetype := range mimetypeQuantity {
			if len(mimetype.Mimetype) > maxLength {
				maxLength = len(mimetype.Mimetype)
			}
			if len(fmt.Sprintf("%d", mimetype.RowCount)) > maxRowCountLength {
				maxRowCountLength = len(fmt.Sprintf("%d", mimetype.RowCount))
			}
			if len(fmt.Sprintf("%d", mimetype.Downloaded)) > maxDownloadLength {
				maxDownloadLength = len(fmt.Sprintf("%d", mimetype.Downloaded))
			}
		}

		mimetypeBoolMap := map[string]bool{}

		for _, mimetype := range mimetypeQuantity {
			if mimetype.RowCount != mimetype.Downloaded {
				mimetypeBoolMap[mimetype.Mimetype] = false
				formatted := fmt.Sprintf(
					"%-*s: %-*d items | Downloaded: %-*d",
					maxLength, mimetype.Mimetype,
					maxRowCountLength, mimetype.RowCount,
					maxDownloadLength, mimetype.Downloaded)

				mimetypeStr := mimetype.Mimetype

				form.AddCheckbox(formatted, false, func(checked bool) {
					mimetypeBoolMap[mimetypeStr] = checked
				})
			}
		}

		form.AddButton("Download", func() {
			mimetypeArr := []string{}

			for mimetype, isTrue := range mimetypeBoolMap {
				if isTrue {
					mimetypeArr = append(mimetypeArr, mimetype)
				}
			}

			mimetypeFiltered := map[string]db.MimetypeQuantity{}

			for _, mimetype := range mimetypeQuantity {
				if slices.Contains(mimetypeArr, mimetype.Mimetype) {
					mimetypeFiltered[mimetype.Mimetype] = mimetype
				}
			}

			textView := tview.NewTextView().
				SetDynamicColors(true).
				SetWrap(true).
				SetWordWrap(true).
				SetTextAlign(tview.AlignCenter).
				SetText(generateDownloadText(mimetypeFiltered, mimetypeArr, "initiated"))

			go downloadLoop(config, textView, mimetypeFiltered, mimetypeArr, websiteURL)

			config.App.SetRoot(textView, true)
		})
		form.SetBorder(true).SetTitle(fmt.Sprintf(`Available mimetypes for "%s"`, websiteURL)).SetTitleAlign(tview.AlignLeft)

		config.App.QueueUpdateDraw(func() {
			config.App.SetRoot(form, true)
		})
	}()
}

func generateDownloadText(mimetypeFiltered map[string]db.MimetypeQuantity, mimetypeArr []string, nowDownloading string) string {
	str := fmt.Sprintf("NOW DOWNLOADING: \"%s\"\n\n", nowDownloading)
	str += strings.Repeat("=", len(str))
	str += "\n\n"

	maxLength := 0
	maxRowCountLength := 0
	maxDownloadLength := 0

	for _, mimetype := range mimetypeFiltered {
		if len(mimetype.Mimetype) > maxLength {
			maxLength = len(mimetype.Mimetype)
		}
		if len(fmt.Sprintf("%d", mimetype.RowCount)) > maxRowCountLength {
			maxRowCountLength = len(fmt.Sprintf("%d", mimetype.RowCount))
		}
		if len(fmt.Sprintf("%d", mimetype.Downloaded)) > maxDownloadLength {
			maxDownloadLength = len(fmt.Sprintf("%d", mimetype.Downloaded))
		}
	}

	for _, mimeStr := range mimetypeArr {
		mimetype := mimetypeFiltered[mimeStr]
		str += fmt.Sprintf(
			"%-*s: %-*d items | Downloaded: %-*d\n",
			maxLength, mimetype.Mimetype,
			maxRowCountLength, mimetype.RowCount,
			maxDownloadLength, mimetype.Downloaded)
	}

	return str
}

func downloadLoop(config *config.Config, textView *tview.TextView, mimetypeFiltered map[string]db.MimetypeQuantity, mimetypeArr []string, websiteURL string) {
	for {
		link, err := config.DB.GetOne(websiteURL, mimetypeArr)
		if err != nil {
			panic(err)
		}

		if link == nil {
			break
		}

		textView.Clear()
		fmt.Fprintf(textView, generateDownloadText(mimetypeFiltered, mimetypeArr, link.Original))
		config.App.Draw()

		api.WaybackDownloadFile(config, *link)
		// @TODO save errors in db and skip
		err = config.DB.UpdateURL(*link)
		if err != nil {
			panic(err)
		}

		currMime := mimetypeFiltered[link.Mimetype]
		currMime.Downloaded++
		mimetypeFiltered[link.Mimetype] = currMime
	}

	modal := tview.NewModal().
		SetText("Downloading finished successfully").
		AddButtons([]string{"OK"}).
		SetBackgroundColor(tcell.ColorGreen).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			config.App.SetRoot(MainMenuView(config), true)
		})

	config.App.QueueUpdateDraw(func() {
		config.App.SetRoot(modal, true)
	})
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
				existingLink, err := config.DB.HasAny(entry.Name())
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
			list.AddItem(folderName, "", rune(index+'0'), nil)
		}
		// Add Back button
		list.AddItem("Back", "Return to main menu", 'b', func() {
			config.App.SetRoot(MainMenuView(config), true)
		})

		list.SetBorder(true).SetTitle("Select website URL").SetTitleAlign(tview.AlignLeft)

		list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			if mainText == "Back" {
				config.App.SetRoot(MainMenuView(config), true)
				return
			}
			Analysis_View(config, mainText)
		})

		config.App.QueueUpdateDraw(func() {
			config.App.SetRoot(list, true)
		})
	}()
}
