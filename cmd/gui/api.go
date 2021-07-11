package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
)

type APIServer struct {
	*http.Server
	mux       *http.ServeMux
	indexPage string
}

func NewAPIServer(listenAddr string) *APIServer {
	indexPage := string(MustAsset("index.html"))
	indexPage = regexp.MustCompile(`[\n\s]+`).ReplaceAllString(indexPage, " ")
	mux := http.NewServeMux()
	return &APIServer{
		Server: &http.Server{
			Addr:    listenAddr,
			Handler: mux,
		},
		mux:       mux,
		indexPage: indexPage,
	}
}

func (s *APIServer) recover(w http.ResponseWriter, r *http.Request) func() {
	return func() {
		if data := recover(); data != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", data)))
		}
	}
}

func (s *APIServer) getHome(w http.ResponseWriter, r *http.Request) {
	defer s.recover(w, r)()
	w.WriteHeader(200)
	w.Write([]byte(s.indexPage))
}

func (s *APIServer) getProjects(w http.ResponseWriter, r *http.Request) {
	defer s.recover(w, r)()
	if json, err := json.MarshalIndent(config.Get().Workspace.Projects, "", "  "); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(200)
		w.Write(json)
	}
}

func (s *APIServer) getVersions(w http.ResponseWriter, r *http.Request) {
	defer s.recover(w, r)()
	cfg := config.Get()
	if !cfg.Workspace.Initialized {
		panic("Workspace has not been initialized yet, run `grlm init`")
	}
	type Response struct{ Name, Version string }
	var response []Response
	for _, p := range cfg.Workspace.Projects {
		if a, err := accessor.Open(p.Path); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		} else {
			var v version.Version
			if v, err = a.ReadVersion(); err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			} else {
				response = append(response, Response{Name: p.Name, Version: v.String()})
			}
		}
	}
	if json, err := json.MarshalIndent(response, "", "  "); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(200)
		w.Write(json)
	}
}

func (s *APIServer) getWorkspace(w http.ResponseWriter, r *http.Request) {
	defer s.recover(w, r)()
	if json, err := json.MarshalIndent(config.Get().Workspace, "", "  "); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(200)
		w.Write(json)
	}
}

func (s *APIServer) Serve() {
	s.mux.HandleFunc("/home", s.getHome)
	s.mux.HandleFunc("/api/projects", s.getProjects)
	s.mux.HandleFunc("/api/versions", s.getVersions)
	s.mux.HandleFunc("/api/workspace", s.getWorkspace)

	s.provideAssets()

	log.Fatal(s.ListenAndServe())
}

func (s *APIServer) provideAssets() error {
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

	s.mux.Handle("/", http.FileServer(AssetFile()))
	// s.provideAsset("app/main.js", "text/javascript")
	// s.provideAsset("app/style.css", "text/css")
	// s.provideAsset("app/pages/projects/projects.html", "text/html")
	// s.provideAsset("app/pages/projects/projects.js", "text/javascript")
	// s.provideAsset("app/pages/home/home.html", "text/html")
	// s.provideAsset("app/pages/home/home.js", "text/javascript")
	// s.provideAsset("app/pages/versions/versions.html", "text/html")
	// s.provideAsset("app/pages/versions/versions.js", "text/javascript")
	return nil
}

func (s *APIServer) provideAsset(name, contentType string) {
	http.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		// Getting the headers so we can set the correct mime type
		println("provide asset: " + name)
		headers := w.Header()
		headers["Content-Type"] = []string{contentType}
		fmt.Fprint(w, string(MustAsset(name)))
	})
}
