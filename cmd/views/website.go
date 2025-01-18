package views

import (
	"fmt"
	"path"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"
	"waybackdownloader/cmd/util"
	"waybackdownloader/cmd/web/api"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DownloadWebsiteURLs_View(config *config.Config) {
	var inputField *tview.InputField

	inputField = tview.NewInputField().
		SetLabel("Enter website url: ").
		SetPlaceholder("E.g. www.youtube.com").
		SetFieldWidth(100).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				websiteURL := inputField.GetText()
				folderPath := path.Join(data.MAIN_PATH, websiteURL)
				util.CreatePathIfNotExists(folderPath)

				drawDownloadInfo(config, websiteURL)

				go func() {
					api.WaybackLinksCollectionSave(config, websiteURL, folderPath)

					config.App.QueueUpdateDraw(func() {
						drawDownloadFinishModal(config, websiteURL)
					})
				}()
			}
		})

	config.App.SetRoot(inputField, true)
}

func drawDownloadInfo(config *config.Config, websiteURL string) {
	textView := tview.NewTextView().
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf(`Downloading website URLs of "%s", this may take a while...`, websiteURL))

	config.App.SetRoot(textView, true)
}

func drawDownloadFinishModal(config *config.Config, websiteURL string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf(`Finished downloading URLs of "%s"`, websiteURL)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			go func() {
				Analysis_View(config, websiteURL)
			}()
		})

	config.App.SetRoot(modal, true)
}
