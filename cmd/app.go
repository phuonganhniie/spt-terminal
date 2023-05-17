package main

import (
	"github.com/aditya-K2/tview"
	"github.com/phuonganhniie/spt-terminal/config"
	"github.com/phuonganhniie/spt-terminal/ui"
)

var (
	ImgX  int
	ImgY  int
	ImgW  int
	ImgH  int
	start = true
)

var (
	App  *tview.Application
	Flex *tview.Flex
	cfg  = config.NewConfig()
)

func onConfigChange() {
	ui.SetStyles()
	ui.SetBorderRunes()
}

func NewApplication() *tview.Application {
	onConfigChange()
	config.NewAppConfig().OnConfigChange = onConfigChange

	App = tview.NewApplication()
	return nil
}
