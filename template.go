package junocfg

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func CheckTemplate(tmpl []byte) (*template.Template, error) {
	helpers := template.FuncMap{
		"copy": copyFunc(bytes.NewBuffer([]byte{})),
	}

	t, err := template.New("template").Funcs(helpers).Parse(string(tmpl))
	return t, err
}

func RenderAndCheckTemplate(tmpl []byte, settingsMap map[string]interface{}) (string, error) {
	buffer := bytes.NewBuffer([]byte{})

	helpers := template.FuncMap{
		"copy": copyFunc(buffer),
	}

	t, err := template.New("template").Funcs(helpers).Parse(string(tmpl))
	if err != nil {
		return "", fmt.Errorf("check template failed [%s]\n", err)
	}

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

//copy returns a template helper that is aware of it's current indent
//and copies given data to the result YAML according to this indent
func copyFunc(buf *bytes.Buffer) func(interface{}) (string, error) {
	return func(in interface{}) (result string, err error) {
		b, err := yaml.Marshal(in)
		if err != nil {
			return "", err
		}

		indent, err := indent(buf)
		if err != nil {
			return "", err
		}

		rr := []string{}
		scanner := bufio.NewScanner(bytes.NewBuffer(b))
		for scanner.Scan() {
			rr = append(rr, scanner.Text()+"\n")
		}

		return strings.Join(rr, indent), scanner.Err()
	}
}

var rxSpaces = regexp.MustCompile("^\\s*$")

func indent(b *bytes.Buffer) (string, error) {
	var line string
	scanner := bufio.NewScanner(bytes.NewBuffer(b.Bytes()))
	for scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if rxSpaces.MatchString(line) {
		return line, nil
	}

	return "", nil
}
