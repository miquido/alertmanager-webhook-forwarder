package root

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

const (
	errorTypeHTTPMethodNotAllowed            = "HTTP_METHOD_NOT_ALLOWED"
	errorTypeHTTPBodyContentTypeNotSupported = "HTTP_BODY_CONTENT_TYPE_NOT_SUPPORTED"
	errorTypeHTTPBodyInvalid                 = "HTTP_BODY_INVALID"
	errorTypeRequiredQueryParameterMissing   = "REQUIRED_QUERY_PARAMETER_MISSING"
	errorTypeQueryParameterInvalid           = "QUERY_PARAMETER_INVALID"
)

func httpJSONError(w http.ResponseWriter, errorType string, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if statusCode < 400 || statusCode > 599 {
		statusCode = 500
	}
	w.WriteHeader(statusCode)
	errorTypeJSON, _ := json.Marshal(errorType)
	errorMessageJSON, _ := json.Marshal(errorMessage)
	_, _ = fmt.Fprintln(w, fmt.Sprintf("{\"error\":{\"type\":%s\"message\":%s}}", errorTypeJSON, errorMessageJSON))
}

func getSuccessStatusCode(r *http.Request) (statusCode int, err error) {
	successStatusCodeQuery := r.URL.Query().Get("successStatusCode")
	if successStatusCodeQuery != "" {
		statusCode, err = strconv.Atoi(successStatusCodeQuery)
		if err != nil {
			return 0, err
		}
	}

	if statusCode < 200 || statusCode > 299 {
		statusCode = http.StatusAccepted
	}

	return statusCode, nil
}

func providerWebhookHandler(provider string, w http.ResponseWriter, r *http.Request) {
	conditional_runner.V(7).DumpRequest(r)
	if r.Method != "POST" {
		httpJSONError(w, errorTypeHTTPMethodNotAllowed, "Only \"POST\" HTTP Method is allowed", http.StatusMethodNotAllowed)
		return
	}

	channelName := r.URL.Query().Get("channel")
	if channelName == "" {
		httpJSONError(w, errorTypeRequiredQueryParameterMissing, "Query parameter \"channel\" is required", http.StatusBadRequest)
		return
	}

	templateName := r.URL.Query().Get("template")
	if templateName == "" {
		httpJSONError(w, errorTypeRequiredQueryParameterMissing, "Query parameter \"template\" is required", http.StatusBadRequest)
		return
	}

	successStatusCode, err := getSuccessStatusCode(r)
	if err != nil {
		httpJSONError(w, errorTypeQueryParameterInvalid, "Value of query parameter \"successStatusCode\" is invalid. Accepted values are numbers between 200 and 299.", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		httpJSONError(w, errorTypeHTTPBodyContentTypeNotSupported, "Currently accepted body \"content-type\" is only \"application/json\".", http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.Error(err)
		httpJSONError(w, errorTypeHTTPBodyInvalid, fmt.Sprintf("Error while reading body: %s", err), http.StatusInternalServerError)
		return
	}

	var requestBodyMap map[string]interface{}
	err = json.Unmarshal(requestBody, &requestBodyMap)
	if err != nil {
		klog.Error(err)
		httpJSONError(w, errorTypeHTTPBodyInvalid, fmt.Sprintf("Error while decoding json: %s", err), http.StatusBadRequest)
		return
	}

	go func() {
		_ = forwarder.Forward(provider, channelName, templateName, requestBodyMap)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(int(successStatusCode))
	_, _ = fmt.Fprint(w, "{\"status\":\"accepted\"}")
}

func init() {
	ServeCmd.Flags().String("host", "0.0.0.0", "Host that server should bind to.")
	ServeCmd.Flags().Int("port", 8080, "Port that server should listen on.")

	RootCmd.AddCommand(ServeCmd)

	_ = viper.BindPFlag("server.host", ServeCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", ServeCmd.Flags().Lookup("port"))
}
