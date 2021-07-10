package node

import (
	"encoding/json"
	"fmt"
	"os"
)

type Package map[string]interface{}

func (p Package) getValue(k string) (interface{}, error) {
	if v, ok := p[k]; !ok {
		return nil, fmt.Errorf("no '%s' key found in package", k)
	} else {
		return v, nil
	}
}

func (p Package) Name() (string, error) {
	if v, err := p.getValue("name"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}
func (p Package) Author() (string, error) {
	if v, err := p.getValue("author"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p Package) Description() (string, error) {
	if v, err := p.getValue("description"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p Package) Contributors() (string, error) {
	if v, err := p.getValue("contributors"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p Package) Maintainers() (string, error) {
	if v, err := p.getValue("maintainers"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p Package) Version() (string, error) {
	if v, err := p.getValue("version"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p Package) Scripts() (map[string]string, error) {
	if v, err := p.getValue("scripts"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p Package) Dependencies() (map[string]string, error) {
	if v, err := p.getValue("dependencies"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p Package) DevDependencies() (map[string]string, error) {
	if v, err := p.getValue("devDependencies"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p Package) Read(b []byte) error {
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	return nil
}

func (p Package) ReadFile(fname string) (err error) {
	if _, err = os.Stat(fname); err == nil || os.IsExist(err) {
		var content []byte
		if content, err = os.ReadFile(fname); err != nil {
			return err
		}
		err = p.Read(content)
	}
	return err
}
