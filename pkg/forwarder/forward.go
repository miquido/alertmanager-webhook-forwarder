package forwarder

import (
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/config"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/forward_channel"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/message_template"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
	"k8s.io/klog"
)

type Forwarder interface {
	Forward(interface{}, *message_template.MessageTemplate, *forward_channel.ForwardChannel) error
}

func Forward(provider string, channelName string, templateName string, input interface{}) (err error) {
	klog.V(7).Infof("Provider: \"%s\"", provider)
	klog.V(7).Infof("Channel: \"%s\"", channelName)
	klog.V(7).Infof("Template: \"%s\"", templateName)
	conditional_runner.V(7).DumpYaml("Input", input)

	forwardChannel, err := forward_channel.NewForwardChannelFromConfig(provider, channelName)
	if err != nil {
		klog.Error(err)
		return err
	}

	conditional_runner.V(7).DumpYaml("Forward channel", forwardChannel)

	templateConfigObject, err := config.GetTemplateObject(provider, templateName)
	if err != nil {
		klog.Error(err)
		return err
	}

	template, err := message_template.NewMessageTemplateFromConfig(utils.ProviderHangoutsChat, templateName, templateConfigObject)
	if err != nil {
		klog.Error(err)
		return err
	}

	conditional_runner.V(7).DumpYaml("Message template", template)

	f, err := Get(provider)
	if err != nil {
		klog.Error(err)
		return err
	}

	return (*f).Forward(input, template, forwardChannel)
}
