package data

import (
	"waybackdownloader/cmd/util"

	"github.com/rivo/tview"
)

type Config struct {
	App *tview.Application
}

func (c *Config) Init() {
	util.CreatePathIfNotExists(MAIN_PATH)
}
