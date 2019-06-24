package message_template

import (
	"encoding/json"
	"errors"
	"github.com/ghodss/yaml"
	"reflect"
	"strings"
)

func toJson(anyObj interface{}) (string, error) {
	result, err := json.Marshal(anyObj)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func toYaml(anyObj interface{}) (string, error) {
	result, err := yaml.Marshal(anyObj)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func indent(input string, indentSize int) string {
	spaces := strings.Repeat(" ", indentSize)
	newline := "\n" + spaces
	return spaces + strings.ReplaceAll(input, "\n", newline)
}

func flexibleIndent(val1 interface{}, val2 interface{}) (string, error) {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)
	if v1.Kind() == reflect.String && v2.Kind() == reflect.Int {
		return indent(v1.String(), int(v2.Int())), nil
	} else if v1.Kind() == reflect.Int && v2.Kind() == reflect.String {
		return indent(v2.String(), int(v1.Int())), nil
	}
	return "", errors.New("indent: expected (string, int) or (int, string) arguments")
}
