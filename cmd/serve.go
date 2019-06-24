package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/forwarder"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/spf13/cobra"
)

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		address := fmt.Sprintf("%s:%d", host, port)

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "OK")
		})
		http.HandleFunc("/v1/webhook/hangouts-chat", providerWebhookHandlerFactory(utils.ProviderHangoutsChat))
		klog.Infof("Listening on http://%s:%d", host, port)
		err = http.ListenAndServe(address, nil)
		if err != nil {
			klog.Error(err)
		}

		return err
	},
}

func providerWebhookHandlerFactory(provider string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		providerWebhookHandler(provider, w, r)
	}
}

func providerWebhookHandler(provider string, w http.ResponseWriter, r *http.Request) {
	conditional_runner.V(7).DumpRequest(r)
	if r.Method != "POST" {
		http.Error(w, "Only \"POST\" HTTP Method is allowed", http.StatusMethodNotAllowed)
		return
	}

	channelName := r.URL.Query().Get("channel")
	if channelName == "" {
		http.Error(w, "Query parameter \"channel\" is required", http.StatusBadRequest)
		return
	}

	templateName := r.URL.Query().Get("template")
	if templateName == "" {
		http.Error(w, "Query parameter \"template\" is required", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Currently accepted body \"content-type\" is only \"application/json\".", http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.Error(err)
		http.Error(w, fmt.Sprintf("Error while reading body: %s", err), http.StatusInternalServerError)
		return
	}

	var requestBodyMap map[string]interface{}
	err = json.Unmarshal(requestBody, &requestBodyMap)
	if err != nil {
		klog.Error(err)
		http.Error(w, fmt.Sprintf("Error while decoding json: %s", err), http.StatusInternalServerError)
		return
	}

	go func() {
		_ = forwarder.Forward(provider, channelName, templateName, requestBodyMap)
	}()

	w.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprint(w, "Accepted")
}

func init() {
	ServeCmd.Flags().String("host", "0.0.0.0", "Host that server should bind to.")
	ServeCmd.Flags().Int("port", 8080, "Port that server should listen on.")

	RootCmd.AddCommand(ServeCmd)

	_ = viper.BindPFlag("server.host", ServeCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", ServeCmd.Flags().Lookup("port"))
}
