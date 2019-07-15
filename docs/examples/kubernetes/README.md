# Kubernetes examples

- [Kubernetes examples](#Kubernetes-examples)
  - [Alertmanager configuration with Hangouts Chat with default template](#Alertmanager-configuration-with-Hangouts-Chat-with-default-template)

## Alertmanager configuration with Hangouts Chat with default template

1. Create [Hangouts Chat Incoming Webhook](https://developers.google.com/hangouts/chat/how-tos/webhooks)
2. Copy Incoming Webhook URL
3. Copy file `configmap.example.yaml` to `configmap.yaml` and replace `%CHANGE_ME_URL%` with copied Incoming Webhook URL
4. Run `kubectl appply -f configmap.yaml deployment.yaml`
5. Configure alertmanager's receiver like this:

    ```yaml
    receivers:
      - name: my-hangouts-chat-receiver
        webhook_configs:
          - url: http://alertmanager-webhook-forwarder/v1/webhook/hangouts-chat?channel=myMonitoringChannel&template=alertmanager
    ```
6. That's all!
