package hangouts_chat

import (
	"errors"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/forward_channel"
	"google.golang.org/api/chat/v1"
	"google.golang.org/api/googleapi"
	"k8s.io/klog"
	"net/url"
	"regexp"
)

type IncomingWebHookConfig struct {
	ApiKey   string `json:"key"`
	Token    string `json:"token"`
	SpaceId  string `json:"spaceId"`
	ThreadId string `json:"threadId,omitempty"`
	Version  string `json:"version"`
}

func (config *IncomingWebHookConfig) Parent() string {
	return "spaces/" + config.SpaceId
}

// uri - e.g. https://chat.googleapis.com/v1/spaces/AAAAB3gQyCk/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=L1bm-eOmL6qN8wyUy44qbQA1wkhVICxYiMerszWoXas%3D
func NewIncomingWebHookConfigFromUrl(incomingWebHookUrl string, threadId string) (config *IncomingWebHookConfig, err error) {
	urlObj, err := url.Parse(incomingWebHookUrl)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	re := regexp.MustCompile("/([a-zA-Z0-9]+)/spaces/([a-zA-Z0-9]+)/messages")
	matches := re.FindStringSubmatch(urlObj.Path)

	config = &IncomingWebHookConfig{
		ApiKey:   urlObj.Query().Get("key"),
		Token:    urlObj.Query().Get("token"),
		ThreadId: threadId,
		SpaceId:  matches[2],
		Version:  matches[1],
	}

	return config, nil
}

func (config *IncomingWebHookConfig) ThreadName() string {
	if config.ThreadId != "" {
		return "spaces/" + config.SpaceId + "/threads/" + config.ThreadId
	}

	return ""
}

func (config *IncomingWebHookConfig) Configure(message *chat.Message) {
	if config.ThreadId != "" {
		message.Thread = &chat.Thread{
			Name: config.ThreadName(),
		}
	}
}

func NewIncomingWebHookConfigFromForwardChannel(forwardChannel *forward_channel.ForwardChannel) (incomingWebHookConfig *IncomingWebHookConfig, err error) {
	if !forwardChannel.HasConfig("url") {
		return nil, errors.New("forward channel must have key \"url\" defined")
	}

	return NewIncomingWebHookConfigFromUrl(forwardChannel.GetConfig("url"), forwardChannel.GetConfig("thread"))
}

func IncomingWebHookTokenOption(t string) googleapi.CallOption {
	return incomingWebHookToken(t)
}

type incomingWebHookToken string

func (t incomingWebHookToken) Get() (string, string) {
	return "token", string(t)
}

func IncomingWebHookApiKeyOption(k string) googleapi.CallOption {
	return incomingWebHookApiKey(k)
}

type incomingWebHookApiKey string

func (k incomingWebHookApiKey) Get() (string, string) {
	return "key", string(k)
}
