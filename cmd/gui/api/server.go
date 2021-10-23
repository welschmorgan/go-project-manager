package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/release"
	"github.com/welschmorgan/go-release-manager/version"
	"gopkg.in/yaml.v2"
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
	var releaseUndoActions = []release.UndoAction{}
	dir := config.Get().Workspace.Path.Join(".grlm", "undos").Expand()
	fmt.Println("Scanning ", dir)
	undoActions := map[string][]release.UndoAction{}
	undoFiles := []Undo{}
	readFile := func(name, path string, fi os.FileInfo) error {
		if content, err := os.ReadFile(path); err != nil {
			return fmt.Errorf("failed to load undo %s, %s", path, err.Error())
		} else {
			if err = yaml.Unmarshal(content, &releaseUndoActions); err != nil {
				return fmt.Errorf("failed to load undo %s, %s", path, err.Error())
			}
			fmt.Println("Adding ", name)
			undoActions[name] = releaseUndoActions
			undoFiles = append(undoFiles, Undo{
				Name: name,
				Date: fi.ModTime().String(),
			})
		}
		return nil
	}
	var readDir func(dir string) (err error)

	readDir = func(dir string) (err error) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read directory: " + err.Error())
		}
		path := ""
		var fi os.FileInfo
		for _, e := range entries {
			fmt.Println("Reading ", e.Name())
			path = filepath.Join(dir, e.Name())
			err = nil
			if fi, err = e.Info(); err != nil {
				log.Errorf("failed to retrieve file infos for %s, %s", path, err.Error())
			} else if fi.IsDir() {
				err = readDir(path)
			} else {
				err = readFile(e.Name(), path, fi)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := readDir(dir); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	sort.Slice(undoFiles, func(i, j int) bool {
		return undoFiles[i].Date < undoFiles[j].Date
	})

	if json, err := json.MarshalIndent(undoFiles, "", "  "); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("failed to marshal undos: " + err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(json))
	}
}

func (s *APIServer) Serve() {
	s.mux.HandleFunc("/home", s.getHome)
	s.mux.HandleFunc("/api/projects", s.getProjects)
	s.mux.HandleFunc("/api/projects/scan", s.getProjects)
	s.mux.HandleFunc("/api/versions", s.getVersions)
	s.mux.HandleFunc("/api/undos", s.getUndos)
	s.mux.HandleFunc("/api/workspace", s.getWorkspace)

	s.provideAssets()

	log.Fatal(s.ListenAndServe())
}

func (s *APIServer) provideAssets() error {
	s.mux.Handle("/", http.FileServer(AssetFile()))
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
