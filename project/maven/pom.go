package maven

import (
	"encoding/xml"
	"io"
	"os"
)

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
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

type POMDependencies map[string]POMDependency

type POMDependenciesXmlEntry struct {
	XMLName xml.Name
	Value   POMDependency `xml:"dependency"`
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (m POMDependencies) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(POMDependenciesXmlEntry{XMLName: xml.Name{Local: k}, Value: v})
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
func (m *POMDependencies) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = POMDependencies{}
	for {
		var e POMDependenciesXmlEntry

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

type POMProject struct {
	Xmlns             string `xml:"xmlns,attr"`
	XmlnsXsi          string `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string `xml:"xsi:schemaLocation,attr"`

	ModelVersion string          `xml:"modelVersion"`
	GroupId      string          `xml:"groupId"`
	ArtifactId   string          `xml:"artifactId"`
	Version      string          `xml:"version"`
	Properties   POMProperties   `xml:"properties"`
	Dependencies POMDependencies `xml:"dependencies"`
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
