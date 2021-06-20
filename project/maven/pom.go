package maven

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

const DefaultPOMModel = POMModel4
const DefaultPOMVersion = "0.1.0-SNAPSHOT"
const DefaultPOMJavaVersion = "1.8"

type POMProperties map[string]string

type POMPropertiesXmlEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (m POMProperties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(POMPropertiesXmlEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML unmarshals the XML into a map of string to strings,
// creating a key in the map for each tag and setting it's value to the
// tags contents.
//
// The fact this function is on the pointer of Map is important, so that
// if m is nil it can be initialized, which is often the case if m is
// nested in another xml structurel. This is also why the first thing done
// on the first line is initialize it.
func (m *POMProperties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = POMProperties{}
	for {
		var e POMPropertiesXmlEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

type POMDependency struct {
	XMLName    xml.Name `xml:"dependency"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Scope      string   `xml:"scope"`
}

type POMDependencyList struct {
	XMLName      xml.Name        `xml:"dependencies"`
	Dependencies []POMDependency `xml:"dependency"`
}

type POMPlugin struct {
	XMLName xml.Name `xml:"plugin"`

	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type POMPluginList struct {
	XMLName xml.Name `xml:"plugins"`

	Plugins []POMPlugin `xml:"plugin"`
}

type POMPluginManagement struct {
	XMLName xml.Name `xml:"pluginManagement"`

	Plugins POMPluginList `xml:"plugins"`
}

type POMBuild struct {
	XMLName xml.Name `xml:"build"`

	PluginManagement POMPluginManagement `xml:"pluginManagement"`
}

type POMProject struct {
	XMLName xml.Name `xml:"project"`

	Xmlns             string `xml:"xmlns,attr"`
	XmlnsXsi          string `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string `xml:"xsi:schemaLocation,attr"`

	ModelVersion POMModelVersion   `xml:"modelVersion"`
	GroupId      string            `xml:"groupId"`
	ArtifactId   string            `xml:"artifactId"`
	Version      string            `xml:"version"`
	Properties   POMProperties     `xml:"properties"`
	Dependencies POMDependencyList `xml:"dependencies"`
	Build        POMBuild          `xml:"build"`
}

func (p *POMProject) SetModelVersion(v POMModelVersion) {
	p.ModelVersion = v
	p.Xmlns = "http://maven.apache.org/POM/" + v.Version()
	p.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	p.XsiSchemaLocation = "http://maven.apache.org/POM/" + v.Version() + " http://maven.apache.org/xsd/maven-" + v.Version() + ".xsd"
}

type POMModelVersion uint8

const (
	POMModelUnknown POMModelVersion = iota
	POMModel1       POMModelVersion = iota
	POMModel2       POMModelVersion = iota
	POMModel3       POMModelVersion = iota
	POMModel4       POMModelVersion = iota
)

func ParseModelVersion(s string) POMModelVersion {
	if s == POMModel1.Version() || s == fmt.Sprint(POMModel1.MajorVersion()) {
		return POMModel1
	}
	if s == POMModel2.Version() || s == fmt.Sprint(POMModel2.MajorVersion()) {
		return POMModel2
	}
	if s == POMModel3.Version() || s == fmt.Sprint(POMModel3.MajorVersion()) {
		return POMModel3
	}
	if s == POMModel4.Version() || s == fmt.Sprint(POMModel4.MajorVersion()) {
		return POMModel4
	}
	return POMModelUnknown
}

func (v POMModelVersion) MajorVersion() uint8 {
	switch v {
	case POMModelUnknown:
		return 0
	case POMModel1:
		return 1
	case POMModel2:
		return 2
	case POMModel3:
		return 3
	case POMModel4:
		return 4
	default:
		panic(fmt.Sprintf("Unknown POMModelVersion: %d", v))
	}
}
func (v POMModelVersion) Version() string {
	return fmt.Sprintf("%d.0.0", v.MajorVersion())
}

func (v POMModelVersion) String() string {
	return v.Version()
}

type POMModelVersionXmlEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (m POMModelVersion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	e.Encode(POMModelVersionXmlEntry{XMLName: xml.Name{Local: "modelVersion"}, Value: m.String()})

	return e.EncodeToken(start.End())
}

// UnmarshalXML unmarshals the XML into a map of string to strings,
// creating a key in the map for each tag and setting it's value to the
// tags contents.
//
// The fact this function is on the pointer of Map is important, so that
// if m is nil it can be initialized, which is often the case if m is
// nested in another xml structurel. This is also why the first thing done
// on the first line is initialize it.
func (m *POMModelVersion) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = POMModel1
	for {
		var e POMModelVersionXmlEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		*m = ParseModelVersion(e.Value)
	}
	return nil
}

func NewPOMProjectWithValues(modelVersion POMModelVersion, groupId, artifactId, version string) *POMProject {
	pf := &POMProject{
		GroupId:    groupId,
		ArtifactId: artifactId,
		Version:    version,
		Properties: POMProperties{
			"maven.compiler.source": DefaultPOMJavaVersion,
			"maven.compiler.target": DefaultPOMJavaVersion,
		},
	}
	pf.SetModelVersion(modelVersion)
	return pf
}

func NewPOMProject() *POMProject {
	return NewPOMProjectWithValues(DefaultPOMModel, "", "", DefaultPOMVersion)
}

func (p *POMProject) Write() ([]byte, error) {
	if data, err := xml.MarshalIndent(*p, "", "  "); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (p *POMProject) Read(b []byte) error {
	return xml.Unmarshal(b, p)
}

func (p *POMProject) WriteFile(fname string) error {
	if xml, err := p.Write(); err != nil {
		return err
	} else {
		return os.WriteFile(fname, xml, 0755)
	}
}

func (p *POMProject) ReadFile(fname string) error {
	if _, err := os.Stat(fname); err == nil || os.IsExist(err) {
		if content, err := os.ReadFile(fname); err != nil {
			return err
		} else {
			if err := p.Read(content); err != nil {
				return err
			}
		}
	}
	return nil
}
