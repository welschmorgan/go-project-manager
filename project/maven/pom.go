package maven

import (
	"encoding/xml"
	"fmt"
	"os"
)

type POMFile map[string]interface{}

func (p POMFile) getValue(k string) (interface{}, error) {
	if v, ok := p[k]; !ok {
		return nil, fmt.Errorf("no '%s' key found in package", k)
	} else {
		return v, nil
	}
}

func (p POMFile) Name() (string, error) {
	if v, err := p.getValue("name"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}
func (p POMFile) Author() (string, error) {
	if v, err := p.getValue("author"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p POMFile) Description() (string, error) {
	if v, err := p.getValue("description"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p POMFile) Contributors() (string, error) {
	if v, err := p.getValue("contributors"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p POMFile) Maintainers() (string, error) {
	if v, err := p.getValue("maintainers"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p POMFile) Version() (string, error) {
	if v, err := p.getValue("version"); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

func (p POMFile) Scripts() (map[string]string, error) {
	if v, err := p.getValue("scripts"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p POMFile) Dependencies() (map[string]string, error) {
	if v, err := p.getValue("dependencies"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p POMFile) DevDependencies() (map[string]string, error) {
	if v, err := p.getValue("devDependencies"); err != nil {
		return nil, err
	} else {
		return v.(map[string]string), nil
	}
}

func (p POMFile) Read(b []byte) error {
	if err := xml.Unmarshal(b, &p); err != nil {
		return err
	}
	return nil
}

func (p POMFile) ReadFile(fname string) error {
	if _, err := os.Stat(fname); err == nil || os.IsExist(err) {
		if content, err := os.ReadFile(fname); err != nil {
			return err
		} else if err := p.Read(content); err != nil {
			return err
		}
	}
	return nil
}
