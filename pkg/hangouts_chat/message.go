package hangouts_chat

import (
	"encoding/json"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/message_template"
	"google.golang.org/api/chat/v1"
	"k8s.io/klog"
)

func ParseMessageTemplate(input interface{}, template *message_template.MessageTemplate) (message *chat.Message, err error) {
	parsedTemplate, err := template.Parse(input)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	switch template.Type {
	case message_template.Jsonnet:
		err = json.Unmarshal([]byte(parsedTemplate), &message)
		if err == nil {
			return message, nil
		}
	case message_template.GoTemplateYAML:
		err = yaml.Unmarshal([]byte(parsedTemplate), &message)
		if err == nil {
			return message, nil
		}
	case message_template.GoTemplateText:
		return &chat.Message{
			Text: parsedTemplate,
		}, nil
	default:
		err = errors.New("not possible")
	}

	return nil, err
}
