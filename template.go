package junocfg

import (
	// "fmt"
	// "os"
	"text/template"
	//	"gopkg.in/yaml.v2"
)

func CheckTemplate(tmpl []byte) (*template.Template, error) {
	t, err := template.New("template").Parse(string(tmpl))
	return t, err
}
