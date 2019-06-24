package config

import (
	"errors"
	"fmt"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
	"github.com/spf13/viper"
	"strings"
)

const (
	viperChannelsRootKey = "channels"
)

func getAvailableChannelProviders() []string {
	return utils.GetListOfKeysBySubOfAllSettings(viperChannelsRootKey, utils.ViperKeyDefaultDelimiter)
}

func getAvailableChannels(provider string) []string {
	return utils.GetListOfKeysBySubOfAllSettings(getChannelsProviderViperKey(provider), utils.ViperKeyDefaultDelimiter)
}

func getChannelsProviderViperKey(provider string) string {
	return viperChannelsRootKey + "." + provider
}

func getChannelViperKey(provider string, name string) string {
	return getChannelsProviderViperKey(provider) + "." + name
}

func GetChannelObject(provider string, name string) (channelObj interface{}, err error) {
	if provider == "" {
		return nil, errors.New("provider name is empty")
	}

	if name == "" {
		return nil, errors.New("channel name is empty")
	}

	channelObj = viper.Get(getChannelViperKey(provider, name))
	if nil != channelObj {
		return channelObj, nil
	}

	availableProviders := getAvailableChannelProviders()
	if len(availableProviders) == 0 {
		return nil, errors.New("currently there are no channel providers defined")
	}

	availableChannels := getAvailableChannels(provider)
	if len(availableProviders) == 0 {
		return nil, errors.New(fmt.Sprintf(
			"currently there are no templates defined in %s provider (available template providers: %s)",
			provider,
			strings.Join(availableProviders, ", "),
		))
	}

	return nil, errors.New(fmt.Sprintf(
		"channel \"%s\" from provider \"%s\" has not been found (available templates: %s)",
		name,
		provider,
		strings.Join(availableChannels, ", "),
	))
}
