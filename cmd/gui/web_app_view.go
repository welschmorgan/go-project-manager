package gui

import (
	"fmt"

	"github.com/webview/webview"
)

type WebAppViewConfig struct {
	width    int
	height   int
	homePage string
	title    string
	hostPort string
	debug    bool
}

type WebAppView struct {
	config WebAppViewConfig
	url    string
}

func NewWebAppView(cfg WebAppViewConfig) *WebAppView {
	return &WebAppView{
		config: cfg,
		url:    fmt.Sprintf("http://%s/%s", cfg.hostPort, cfg.homePage),
	}
}

func (a *WebAppView) Start() {
	view := webview.New(a.config.debug)
	defer view.Destroy()
	view.SetTitle(a.config.title)
	view.SetSize(a.config.width, a.config.height, webview.HintNone)
	view.Navigate(a.url)
	view.Run()
}
