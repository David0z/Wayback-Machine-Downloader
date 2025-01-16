package analysis

import "github.com/rivo/tview"

func AnalysisView(websiteURL string) tview.Primitive {
	modal := tview.NewModal().
		SetText("Analysis View Text").
		SetTitle("Analysis View").
		SetBorder(true)

	return modal
}
