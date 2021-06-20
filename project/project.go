package project

import (
	"fmt"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/project/maven"
	"github.com/welschmorgan/go-release-manager/project/node"
)

func init() {
	accessor.Register(&maven.ProjectAccessor{})
	accessor.Register(&node.ProjectAccessor{})
}

func Detect(p string) (accessor.ProjectAccessor, error) {
	for _, a := range accessor.GetAll() {
		log.Infof("%sdetect project: %s - %s\n", strings.Repeat("\t", config.Get().Indent), p, a.AccessorName())
		if ok, err := a.Detect(p); ok {
			if err != nil {
				return nil, err
			}
			return a, nil
		}
	}
	return nil, fmt.Errorf("no accessor found for '%s'", p)
}

func Open(p string) (accessor.ProjectAccessor, error) {
	if a, err := Detect(p); err != nil {
		return nil, err
	} else {
		if err = a.Open(p); err != nil {
			return nil, err
		}
		return a, nil
	}
}
