package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
	"github.com/spf13/viper"
)

const (
	viperTemplatesRootKey = "templates"
)

func getAvailableTemplateProviders() []string {
	return utils.GetListOfKeysBySubOfAllSettings(viperTemplatesRootKey, utils.ViperKeyDefaultDelimiter)
}

func getAvailableTemplates(provider string) []string {
	return utils.GetListOfKeysBySubOfAllSettings(getTemplatesProviderViperKey(provider), utils.ViperKeyDefaultDelimiter)
}

func getTemplatesProviderViperKey(provider string) string {
	return viperTemplatesRootKey + "." + provider
}

func getTemplateViperKey(provider string, name string) string {
	return getTemplatesProviderViperKey(provider) + "." + name
}

func GetTemplateObject(provider string, name string) (templateObj interface{}, err error) {
	if provider == "" {
		return nil, errors.New("provider name is empty")
	}

	if name == "" {
		return nil, errors.New("template name is empty")
	}

	templateObj = viper.Get(getTemplateViperKey(provider, name))
	if nil != templateObj {
		return templateObj, nil
	}

	availableProviders := getAvailableTemplateProviders()
	if len(availableProviders) == 0 {
		return nil, errors.New("currently there are no template providers defined")
	}

	availableTemplates := getAvailableTemplates(provider)
	if len(availableProviders) == 0 {
		return nil, fmt.Errorf(
			"currently there are no templates defined in %s provider (available template providers: %s)",
			provider,
			strings.Join(availableProviders, ", "),
		)
	}

	return nil, fmt.Errorf(
		"template \"%s\" from provider \"%s\" has not been found (available templates: %s)",
		name,
		provider,
		strings.Join(availableTemplates, ", "),
	)
}
