package hangouts_chat

import (
	"context"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/log/client"
	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
	"k8s.io/klog"
	"net/http"
)

var defaultContext context.Context
var defaultHttpClient *http.Client

func getDefaultContext() context.Context {
	if nil == defaultContext {
		defaultContext = context.Background()
	}
	return defaultContext
}

func getDefaultHttpClient() *http.Client {
	if nil == defaultHttpClient {
		defaultHttpClient = client.NewKlogLoggedClient()
	}
	return defaultHttpClient
}

func NewChatSpacesMessagesService() (service *chat.SpacesMessagesService, err error) {
	chatService, err := chat.NewService(getDefaultContext(), option.WithHTTPClient(getDefaultHttpClient()))
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	service = chat.NewSpacesMessagesService(chatService)

	return service, nil
}

func SendIncomingWebHookMessage(service *chat.SpacesMessagesService, message *chat.Message, config *IncomingWebHookConfig) (*chat.Message, error) {
	return service.Create(config.Parent(), message).Do(
		IncomingWebHookTokenOption(config.Token),
		IncomingWebHookApiKeyOption(config.ApiKey))
}
