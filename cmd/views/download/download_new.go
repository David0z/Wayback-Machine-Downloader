package view_download

import (
	"fmt"
	"path"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/util"
	"waybackdownloader/cmd/views/analysis"
	"waybackdownloader/cmd/web/api"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DownloadNew(app *tview.Application) {
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

				drawDownloadInfo(app, websiteURL)

				go func() {
					api.WaybackCollectionSave(websiteURL, folderPath)

					app.QueueUpdateDraw(func() {
						drawDownloadFinishModal(app, websiteURL, analysis.AnalysisView(websiteURL))
					})
				}()
			}
		})

	app.SetRoot(inputField, true)
}

func drawDownloadInfo(app *tview.Application, websiteURL string) {
	textView := tview.NewTextView().
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf(`Downloading website URLs of "%s", this may take a while...`, websiteURL))

	app.SetRoot(textView, true)
}

func drawDownloadFinishModal(app *tview.Application, websiteURL string, nextView tview.Primitive) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf(`Finished downloading URLs of "%s"`, websiteURL)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(nextView, true)
		})

	app.SetRoot(modal, true)
}
