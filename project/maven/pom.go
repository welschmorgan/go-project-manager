package maven

import (
	"encoding/xml"
	"os"
)

type POMDependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

type POMProject struct {
	Xmlns             string `xml:"xmlns,attr"`
	XmlnsXsi          string `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string `xml:"xsi:schemaLocation,attr"`

	ModelVersion string                   `xml:"modelVersion"`
	GroupId      string                   `xml:"groupId"`
	ArtifactId   string                   `xml:"artifactId"`
	Version      string                   `xml:"version"`
	Properties   map[string]string        `xml:"properties"`
	Dependencies map[string]POMDependency `xml:"dependencies"`
}

type POMFile struct {
	Root *POMProject `xml:"project"`
}

func NewPOMFile(modelVersion string) POMFile {
	return POMFile{
		Root: &POMProject{
			Xmlns:             "http://maven.apache.org/POM/" + modelVersion,
			XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
			XsiSchemaLocation: "http://maven.apache.org/POM/" + modelVersion + " http://maven.apache.org/xsd/maven-" + modelVersion + ".xsd",
			ModelVersion:      modelVersion,
			Properties: map[string]string{
				"maven.compiler.source": "1.8",
				"maven.compiler.target": "1.8",
			},
		},
	}
}

func (p POMFile) Write(b []byte) error {
	if data, err := xml.MarshalIndent(&p, "", "  "); err != nil {
		return err
	} else {
		copy(data, b)
	}
	return nil
}

func (p POMFile) Read(b []byte) error {
	return xml.Unmarshal(b, &p)
}

func (p POMFile) WriteFile(fname string) error {
	xml := []byte{}
	if err := p.Write(xml); err != nil {
		return err
	}
	return os.WriteFile(fname, xml, 0755)
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
