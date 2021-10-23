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
	webview.WebView
	config WebAppViewConfig
	url    string
}

func NewWebAppView(cfg WebAppViewConfig) *WebAppView {
	return &WebAppView{
		WebView: nil,
		config:  cfg,
		url:     fmt.Sprintf("http://%s/%s", cfg.hostPort, cfg.homePage),
	}
}

func (a *WebAppView) Start() {
	a.WebView = webview.New(a.config.debug)
	defer a.Destroy()
	a.Bind("setTitle", func(args string) error {
		a.SetTitle(args)
		return nil
	})
	a.SetTitle(a.config.title)
	a.SetSize(a.config.width, a.config.height, webview.HintNone)
	a.Navigate(a.url)
	a.Run()
}
