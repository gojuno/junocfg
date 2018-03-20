package junocfg

import (
	"bytes"
	"fmt"
	"text/template"
	//	"gopkg.in/yaml.v2"
)

func CheckTemplate(tmpl []byte) (*template.Template, error) {
	t, err := template.New("template").Parse(string(tmpl))
	return t, err
}

func RenderAndCheckTemplate(tmpl []byte, settingsMap map[string]interface{}) (string, error) {
	t, err := CheckTemplate(tmpl)
	if err != nil {
		return "", fmt.Errorf("check template failed [%s]\n", err)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := t.Execute(buffer, settingsMap); err != nil {
		return "", fmt.Errorf("failed to render template [%s]\n", err)
	}

	buffer = PreprocessYaml(buffer)
	// check yaml
	if err = CheckYaml(buffer.Bytes()); err != nil {
		fmt.Printf("%s\n", buffer.Bytes())
		return "", fmt.Errorf("Not valid output yaml: %s\n", err.Error())
	}
	return buffer.String(), nil
}
