package hangouts_chat

import (
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/forward_channel"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/message_template"
	"google.golang.org/api/chat/v1"
	"k8s.io/klog"
)

type HangoutsChatForwarder struct {
	SpacesMessageService *chat.SpacesMessagesService
}

func NewHangoutsChatForwarder() (forwarder *HangoutsChatForwarder, err error) {
	service, err := NewChatSpacesMessagesService()
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return &HangoutsChatForwarder{
		SpacesMessageService: service,
	}, nil
}

func (f *HangoutsChatForwarder) Forward(input interface{}, template *message_template.MessageTemplate, forwardChannel *forward_channel.ForwardChannel) (err error) {
	incomingWebHookConfig, err := NewIncomingWebHookConfigFromForwardChannel(forwardChannel)
	if err != nil {
		klog.Error(err)
		return err
	}

	conditional_runner.V(7).DumpYaml("Hangouts Chat IncomingWebHookConfig", incomingWebHookConfig)

	message, err := ParseMessageTemplate(input, template)
	if err != nil {
		klog.Error(err)
		return err
	}
	incomingWebHookConfig.Configure(message)

	conditional_runner.V(7).DumpYaml("Hangouts Chat Message", message)

	return conditional_runner.NotDryRun().RunE(func() (err error) {
		_, err = SendIncomingWebHookMessage(f.SpacesMessageService, message, incomingWebHookConfig)
		if err != nil {
			klog.Error(err)
		}
		return err
	})
}
