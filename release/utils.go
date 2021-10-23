package release

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/version"
)

func rollbackRelease(releases *[]*Release) {
	fmt.Fprintf(os.Stderr, "Rolling back all releases...\n")
	fs.DumpDirStack(os.Stderr)
	errs := []error{}
	for _, r := range *releases {
		// fmt.Fprintf(os.Stderr, "  Rolling back %s v%s:\n", r.Project.Name, r.Context.version)
		// for _, u := range r.UndoActions {
		// 	fmt.Fprintf(os.Stderr, "    - Undo '%s' (path: %s, action: %s)\n", u.Title, u.Path, u.Name)
		// }
		if err := r.Undo(); err != nil {
			errs = append(errs, err)
		}
	}
	errStr := ""
	for _, e := range errs {
		if len(errStr) > 0 {
			errStr += "\n"
		}
		errStr += e.Error()
	}
	if len(errStr) > 0 {
		panic(errStr)
	}
}

func DoRelease(mode string) (err error) {
	typ, err := version.ParsePart(mode)
	if err != nil {
		return err
	}
	config.Get().ReleaseType = typ

	var releases = []*Release{}
	if !config.Get().Workspace.Initialized {
		panic("Workspace has not been initialized yet, run `grlm init`")
	}

	// cleanup release on ctrl-c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGQUIT)

	go func() {
		<-sigs
		rollbackRelease(&releases)
		os.Exit(0)
	}()

	errs := make(map[string]string)
	for _, prj := range config.Get().Workspace.Projects {
		if r, err := NewRelease(prj); err != nil {
			return err
		} else {
			if err = r.PrepareContext(); err != nil {
				errs[prj.Name] = err.Error()
			} else {
				releases = append(releases, r)
			}
		}
	}

	if err := preCheckRelease(&releases, errs); err != nil {
		rollbackRelease(&releases)
		return ErrUserAbort
	}

	for _, r := range releases {
		if err = r.Do(); err != nil {
			return err
		}
	}

	println("Check if everything is OK, if it isn't, answering 'no' now will rollback what has been done.")
	if ok, err := ui.AskYN("Is everything ok"); err != nil || !ok {
		rollbackRelease(&releases)
		return ErrUserAbort
	}

	if err := SaveReleases(releases); err != nil {
		rollbackRelease(&releases)
		return ErrUserAbort
	}
	return nil
}

func preCheckRelease(releases *[]*Release, errs map[string]string) error {
	confirmed := true
	if len(errs) > 0 {
		confirmed = false
		log.Errorf("Some project(s) cannot be released:")
		for key, msg := range errs {
			log.Errorf("\t- [%s] %s", key, msg)
		}
		log.Errorf("However, %d projects can be released:", len(*releases))
		for _, r := range *releases {
			log.Errorf("\t- %s", r.Project.Name)
		}
		if confirmed, _ = ui.AskYN("Do you still want to release them"); !confirmed {
			*releases = make([]*Release, 0)
			return ErrUserAbort
		}
	}
	return nil
}

func LoadRelease(version string) ([]*Release, error) {
	releases := []*Release{}
	version = strings.TrimSuffix(version, ".yaml")
	dir := config.Get().Workspace.Path.Join(".grlm", "releases").Expand()
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, version+".yaml")
	if content, err := os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("could not read file %s, %s", path, err.Error())
	} else if err = json.Unmarshal(content, &releases); err != nil {
		return nil, fmt.Errorf("could not unmarshal release train: %s", err.Error())
	}
	return releases, nil
}

func LoadReleaseTrain() (map[string][]*Release, error) {
	dir := config.Get().Workspace.Path.Join(".grlm", "releases").Expand()
	if entries, err := os.ReadDir(dir); err != nil {
		return nil, err
	} else {
		ret := map[string][]*Release{}
		for _, e := range entries {
			version := strings.TrimSuffix(e.Name(), ".yaml")
			if releases, err := LoadRelease(version); err != nil {
				return nil, err
			} else {
				ret[version] = append(ret[version], releases...)
			}
		}
		return ret, nil
	}
}

func SaveReleases(releases []*Release) error {
	if jsonData, err := json.MarshalIndent(releases, "", "  "); err != nil {
		return fmt.Errorf("could not marshal release train to json: %s", err.Error())
	} else {
		dir := config.Get().Workspace.Path.Join(".grlm", "releases").Expand()
		os.MkdirAll(dir, 0755)
		path := filepath.Join(dir, releases[0].Context.Version.String()+".yaml")
		if err = os.WriteFile(path, jsonData, 0755); err != nil {
			return fmt.Errorf("could not write file %s, %s", path, err.Error())
		}
	}
	return nil
}
