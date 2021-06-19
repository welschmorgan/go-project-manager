package project

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/project/maven"
	"github.com/welschmorgan/go-release-manager/project/node"
)

var All []accessor.ProjectAccessor = []accessor.ProjectAccessor{
	&node.ProjectAccessor{},
	&maven.ProjectAccessor{},
}
var AllNames = []string{}

func init() {
	for _, a := range All {
		AllNames = append(AllNames, a.Name())
	}
}

func instanciate(a accessor.ProjectAccessor) accessor.ProjectAccessor {
	inst := reflect.New(reflect.TypeOf(a).Elem())
	return inst.Interface().(accessor.ProjectAccessor)
}

func Get(n string) accessor.ProjectAccessor {
	loName := strings.ToLower(n)
	for _, a := range All {
		if strings.ToLower(a.Name()) == loName {
			return instanciate(a)
		}
	}
	return nil
}

func Detect(p string) (accessor.ProjectAccessor, error) {
	for _, a := range All {
		if config.Get().Verbose {
			fmt.Printf("%sdetect project: %s - %s\n", strings.Repeat("\t", config.Get().Indent), p, a.Name())
		}
		if ok, err := a.Detect(p); ok {
			if err != nil {
				return nil, err
			}
			return instanciate(a), nil
		}
	}
	return nil, fmt.Errorf("no accessor found for '%s'", p)
}

func Open(p string) (accessor.ProjectAccessor, error) {
	if a, err := Detect(p); err != nil {
		return nil, err
	} else {
		ret := instanciate(a)
		if err = ret.Open(p); err != nil {
			return nil, err
		}
		return ret, nil
	}
}
