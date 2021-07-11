package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webview/webview"
	"github.com/welschmorgan/go-release-manager/config"
)

func expandForLoops(content string, data map[string]interface{}) (string, error) {
	forPattern := regexp.MustCompile(`(?mis)\{\{\s*for\s+([\w]+)\s+in\s+([\w]+)\s*\}\}([\w\W.]*){{endfor\s+([\w]+)\s*}}`)
	forMatches := forPattern.FindAllStringSubmatch(content, -1)

	for _, match := range forMatches {
		varName := strings.TrimSpace(match[1])
		listName := strings.TrimSpace(match[2])
		body := strings.TrimSpace(match[3])
		endListName := strings.TrimSpace(match[4])
		println("varName = " + varName)
		println("listName = " + listName)
		println("body = " + body)
		println("endListName = " + endListName)
		if !strings.EqualFold(listName, endListName) {
			return "", fmt.Errorf("syntax error: invalid expression: %s", match[0])
		}
		rv := reflect.ValueOf(data[listName])
		for i := 0; i < rv.Len(); i++ {
			varPattern := regexp.MustCompile(`\{\{\s*` + varName + `([\w\.]+)\s*\}\}`)
			varMatches := varPattern.FindAllStringSubmatch(body, -1)
			rfItem := reflect.Indirect(rv.Index(i))
			for _, varMatch := range varMatches {
				varMatch[1] = strings.TrimPrefix(varMatch[1], ".")
				body = strings.ReplaceAll(body, varMatch[0], rfItem.FieldByName(varMatch[1]).String())
			}
		}
		content = strings.Replace(content, match[0], body, -1)
	}
	return content, nil
}

func provideAssets() error {
	// dir, err := AssetDir("")
	// if err != nil {
	// 	return err
	// }
	// mimeTypes := map[string]string {
	// 	".js": "text/javascript",
	// 	".css": "text/css",
	// 	".html": "text/html",
	// }
	// for _, d := range dir {
	// 	mimeType := "text/plain"
	// 	for k, v := range mimeTypes {
	// 		if d
	// 	}
	// 	provideAsset(d, "")

	provideAsset("app/main.js", "text/javascript")
	provideAsset("app/style.css", "text/css")
	provideAsset("app/pages/projects.html", "text/html")
	provideAsset("app/pages/projects.js", "text/javascript")
	provideAsset("app/pages/home.html", "text/html")
	provideAsset("app/pages/home.js", "text/javascript")
	return nil
}

func provideAsset(name, contentType string) {
	http.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		// Getting the headers so we can set the correct mime type
		println("provide asset: " + name)
		headers := w.Header()
		headers["Content-Type"] = []string{contentType}
		fmt.Fprint(w, string(MustAsset(name)))
	})
}

var Command = &cobra.Command{
	Use:   "gui",
	Short: "Interface to show workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		debug := true
		w := webview.New(debug)
		indexPage := string(MustAsset("index.html"))
		indexPage = regexp.MustCompile(`[\n\s]+`).ReplaceAllString(indexPage, " ")
		defer w.Destroy()
		go func() {
			http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
				data := map[string]interface{}{
					"projects": config.Get().Workspace.Projects,
				}
				defer func() {
					if data := recover(); data != nil {
						w.WriteHeader(500)
						w.Write([]byte(fmt.Sprintf("%v", data)))
					}
				}()

				if b, err := expandForLoops(string(indexPage), data); err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
				} else {
					indexPage = b
				}

				w.WriteHeader(200)
				w.Write([]byte(indexPage))
			})
			http.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
				if json, err := json.MarshalIndent(config.Get().Workspace.Projects, "", "  "); err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				} else {
					w.WriteHeader(200)
					w.Write(json)
				}
			})
			provideAssets()
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()
		w.SetTitle("GRLM:UI")
		w.SetSize(800, 600, webview.HintNone)
		w.Navigate("http://localhost:8080/home")
		w.Run()
		return nil
	},
}
