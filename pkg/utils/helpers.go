package utils

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

const (
	ViperKeyDefaultDelimiter = "."
	ProviderHangoutsChat     = "hangoutsChat"
)

func GetListOfKeys(mapOfObjs map[string]interface{}) []string {
	keys := make([]string, len(mapOfObjs))
	i := 0
	for key := range mapOfObjs {
		keys[i] = key
		i += 1
	}
	return keys
}

func InterfaceToStringMapOfInterfaces(in interface{}) (out map[string]interface{}, err error) {
	out, ok := in.(map[string]interface{})
	if !ok {
		return nil, errors.New("error converting \"interface{}\" to \"map[string]interface{}\"")
	}

	return out, nil
}

func GetListOfKeysBySubOfAllSettings(key string, delimiter string) (keys []string) {
	settings := viper.AllSettings()
	for _, keyPart := range ParseKeyToParts(key, delimiter) {
		settings, _ = InterfaceToStringMapOfInterfaces(settings[keyPart])
	}

	return GetListOfKeys(settings)
}

func ParseKeyToParts(key string, delimiter string) []string {
	return strings.Split(strings.ToLower(key), delimiter)
}
