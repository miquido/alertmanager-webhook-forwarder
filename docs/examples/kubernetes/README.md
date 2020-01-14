# Kubernetes examples

- [Kubernetes examples](#kubernetes-examples)
  - [Alertmanager configuration with Hangouts Chat with default template](#alertmanager-configuration-with-hangouts-chat-with-default-template)

## Alertmanager configuration with Hangouts Chat with default template

1. Create [Hangouts Chat Incoming Webhook](https://developers.google.com/hangouts/chat/how-tos/webhooks)
2. Copy Incoming Webhook URL
3. Copy files from this directory and update `config.yaml` according to your needs
4. Run `kustomize build | kubectl apply -f -`
5. Configure alertmanager's receiver like this:

    ```yaml
    receivers:
      - name: my-hangouts-chat-receiver
        webhook_configs:
          - url: http://alertmanager-webhook-forwarder/v1/webhook/hangouts-chat?channel=myMonitoringChannel&template=alertmanager
    ```
6. That's all!
