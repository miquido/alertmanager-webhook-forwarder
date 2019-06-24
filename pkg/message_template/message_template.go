package message_template

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/go-jsonnet"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/log"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog"
	"text/template"
)

type MessageTemplateType string

const (
	Jsonnet        MessageTemplateType = "jsonnet"
	GoTemplateYAML MessageTemplateType = "gotemplate-yaml"
	GoTemplateText MessageTemplateType = "gotemplate-text"
)

func NewMessageTemplateType(templateType string) (MessageTemplateType, error) {
	switch templateType {
	case "jsonnet":
		return Jsonnet, nil
	case "gotemplate-yaml":
		return GoTemplateYAML, nil
	case "gotemplate-text":
		return GoTemplateText, nil
	default:
		return "", errors.New("MessageTemplateType constructor accepts only values: jsonnet, gotemplate-yaml, gotemplate-text")
	}
}

type MessageTemplateInput struct {
	Text string `json:"text"`
}

func NewMessageTemplateInputFromArgs(args []string) *MessageTemplateInput {
	return &MessageTemplateInput{
		Text: args[0],
	}
}

type MessageTemplate struct {
	Name     string              `json:"name"`
	Provider string              `json:"provider"`
	Type     MessageTemplateType `json:"type"`
	Template string              `json:"template"`
}

func (tmpl *MessageTemplate) GetParseFn() func(interface{}, *MessageTemplate) (string, error) {
	switch tmpl.Type {
	case Jsonnet:
		return parseJsonnetMessageTemplate
	case GoTemplateYAML:
		return parseGoTemplateYamlMessageTemplate
	case GoTemplateText:
		return parseGoTemplateTextMessageTemplate
	default:
		klog.V(3).Infof("message_template::GetParseFn - not possible condition")
		panic("Not possible")
	}
}

func (tmpl *MessageTemplate) Parse(input interface{}) (string, error) {
	return tmpl.GetParseFn()(input, tmpl)
}

func parseGoTemplateYamlMessageTemplate(input interface{}, tmpl *MessageTemplate) (message string, err error) {
	if tmpl.Type != GoTemplateYAML {
		return "", errors.New("message_template::parseGoTemplateYamlMessageTemplate accepts only message templates with 'gotemplate-yaml' type")
	}

	messageTmpl, err := parseGoTemplate(tmpl)

	if err != nil {
		klog.Error(err)
		return "", err
	}

	parsedTemplateBuffer := new(bytes.Buffer)
	err = messageTmpl.ExecuteTemplate(parsedTemplateBuffer, "message", input)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	return parsedTemplateBuffer.String(), nil
}

func parseJsonnetMessageTemplate(input interface{}, tmpl *MessageTemplate) (message string, err error) {
	if tmpl.Type != Jsonnet {
		return "", errors.New("message_template::parseJsonnetMessageTemplate accepts only message templates with 'jsonnet' type")
	}

	extCodeInput, err := json.Marshal(input)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	extCodeInputString := string(extCodeInput)

	conditional_runner.V(7).Run(func() {
		log.DumpMultilineText("input.json", "JSON", extCodeInputString)
	})

	vm := jsonnet.MakeVM()
	vm.ExtCode("input", extCodeInputString)
	return vm.EvaluateSnippet("message", tmpl.Template)
}

func parseGoTemplate(tmpl *MessageTemplate) (goTemplate *template.Template, err error) {
	return template.New("message").Funcs(template.FuncMap{
		"toYaml": toYaml,
		"toJson": toJson,
		"indent": flexibleIndent,
	}).Parse(tmpl.Template)
}

func parseGoTemplateTextMessageTemplate(input interface{}, tmpl *MessageTemplate) (message string, err error) {
	if tmpl.Type != GoTemplateText {
		return "", errors.New("message_template::parseGoTemplateTextMessageTemplate accepts only message templates with 'gotemplate-text' type")
	}

	messageTmpl, err := parseGoTemplate(tmpl)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	parsedTemplateBuffer := new(bytes.Buffer)
	err = messageTmpl.ExecuteTemplate(parsedTemplateBuffer, "message", input)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	return parsedTemplateBuffer.String(), nil
}

func NewMessageTemplateFromConfig(provider string, name string, configObject interface{}) (template *MessageTemplate, err error) {
	err = mapstructure.Decode(configObject, &template)
	if err != nil {
		klog.Error(err)
	}

	template.Name = name
	template.Provider = provider

	return template, err
}
