package forward_channel

import (
	"fmt"

	"github.com/miquido/alertmanager-webhook-forwarder/pkg/config"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog"
)

type ForwardChannel struct {
	Name     string            `json:"channelName"`
	Provider string            `json:"provider"`
	Config   map[string]string `json:"config"`
}

func (forwardChannel *ForwardChannel) HasConfig(key string) (ok bool) {
	_, ok = forwardChannel.Config[key]

	return ok
}

func (forwardChannel *ForwardChannel) GetConfig(key string) string {
	if value, ok := forwardChannel.Config[key]; ok {
		return value
	}

	return ""
}

func NewForwardChannelFromConfig(provider string, channelName string) (forwardChannel *ForwardChannel, err error) {
	forwardChannel = &ForwardChannel{
		Name:     channelName,
		Provider: provider,
	}

	channelConfigObject, err := config.GetChannelObject(provider, channelName)
	if err != nil {
		err = fmt.Errorf("Config \"%s.%s\": %s", provider, channelName, err)
		klog.Error(err)
		return nil, err
	}

	err = mapstructure.Decode(channelConfigObject, &forwardChannel.Config)
	if err != nil {
		err = fmt.Errorf("Config \"%s.%s\" is invalid: %s", provider, channelName, err)
		klog.Error(err)
	}

	return forwardChannel, err
}
