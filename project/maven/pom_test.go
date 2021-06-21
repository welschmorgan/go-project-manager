package maven

import (
	"strings"
	"testing"
)

func TestPOMProject(t *testing.T) {
	pf := NewPOMProjectWithValues(POMModel4, "com.test.app", "MyApp", "1.0.0")
	if xml, err := pf.Write(); err != nil {
		t.Fatalf("failed to write POM's XML, %s", err.Error())
	} else {
		actual := strings.TrimSpace(string(xml))
		expectedLines := []string{
			`<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">`,
			`<modelVersion>4.0.0</modelVersion>`,
			`<groupId>com.test.app</groupId>`,
			`<artifactId>MyApp</artifactId>`,
			`<version>1.0.0</version>`,
			`</project>`,
		}
		for _, line := range expectedLines {
			if !strings.Contains(actual, line) {
				t.Errorf("Missing project declaration, expected to find '%s' in:\n%s", line, actual)
			}
		}
	}
}
