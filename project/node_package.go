package project

import (
	"encoding/json"
	"fmt"
	"os"
)

type NodePackage map[string]interface{}

func (p NodePackage) getValue(k string) (interface{}, error) {
	if v, ok := p[k]; !ok {
		return nil, fmt.Errorf("no '%s' key found in package", k)
	} else {
		return v, nil
	}
}

func (p NodePackage) Name() (string, error) {
	if v, err := p.getValue("name"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}
func (p NodePackage) Author() (string, error) {
	if v, err := p.getValue("author"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p NodePackage) Description() (string, error) {
	if v, err := p.getValue("description"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p NodePackage) Contributors() (string, error) {
	if v, err := p.getValue("contributors"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p NodePackage) Maintainers() (string, error) {
	if v, err := p.getValue("maintainers"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p NodePackage) Version() (string, error) {
	if v, err := p.getValue("version"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p NodePackage) Scripts() (map[string]string, error) {
	if v, err := p.getValue("scripts"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p NodePackage) Dependencies() (map[string]string, error) {
	if v, err := p.getValue("dependencies"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p NodePackage) DevDependencies() (map[string]string, error) {
	if v, err := p.getValue("devDependencies"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p NodePackage) Read(b []byte) error {
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	return nil
}

func (p NodePackage) ReadFile(fname string) error {
	if _, err := os.Stat(fname); err == nil || os.IsExist(err) {
		if content, err := os.ReadFile(fname); err != nil {
			return err
		} else if err := p.Read(content); err != nil {
			return err
		}
	}
	return nil
}
