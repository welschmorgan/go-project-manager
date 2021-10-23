package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/welschmorgan/go-release-manager/api/controllers"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/release"
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

func (s *APIServer) getUndos(w http.ResponseWriter, r *http.Request) {
	undos, err := release.ListUndos()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "failed to list undo actions: %s", err.Error())
		return
	}

	if json, err := json.MarshalIndent(undos, "", "  "); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("failed to marshal undos: " + err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(json))
	}
}

func (s *APIServer) Controller(prefix string, e http.Handler) *APIServer {
	prefix = strings.TrimSuffix(prefix, "/")
	s.mux.Handle(prefix+"/", http.StripPrefix(prefix, e))
	return s
}

func (s *APIServer) Serve() {
	s.mux.HandleFunc("/home", s.getHome)
	s.mux.HandleFunc("/api/projects", s.getProjects)
	s.mux.HandleFunc("/api/projects/scan", s.getProjects)
	s.mux.HandleFunc("/api/versions", s.getVersions)
	s.mux.HandleFunc("/api/workspace", s.getWorkspace)
	s.Controller("/api/release", controllers.NewReleaseEndpoint())
	// s.mux.Handle("/api/release", http.StripPrefix("/api/release", controllers.NewReleaseEndpoint()))

	rv := reflect.ValueOf(s.mux).Elem()
	routes := rv.FieldByName("m")
	for _, k := range routes.MapKeys() {
		log.Info(k.String())
	}

	s.provideAssets()

	log.Infof("Starting api server on '%s'", s.Addr)
	log.Fatal(s.ListenAndServe())
	log.Infof("Stopped api server")
}

func (s *APIServer) provideAssets() error {
	for _, asset := range AssetNames() {
		s.provideAsset(asset)
	}
	// AssetFile()
	// s.mux.Handle("/", http.FileServer(assetFs()))
	return nil
}

var contentTypes = map[string]string{
	".css":  "text/css;charset=UTF-8",
	".js":   "text/javascript;charset=UTF-8",
	".html": "text/html;charset=UTF-8",
}

func (s *APIServer) provideAsset(asset string) {
	s.mux.HandleFunc("/"+asset, func(w http.ResponseWriter, r *http.Request) {
		defer s.recover(w, r)()
		// Getting the headers so we can set the correct mime type
		uri := r.URL.String()
		name := strings.TrimPrefix(uri, "/")
		ctype := contentTypes[filepath.Ext(name)]
		if content, err := Asset(name); err != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Asset not found: %s", name)
		} else {
			w.Header().Set("Content-Type", ctype)
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(200)
			w.Write(content)
		}
	})
}
