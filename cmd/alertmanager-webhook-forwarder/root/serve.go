package root

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
	Short: "Run as HTTP server",
	Long:  `
Run as golang HTTP server with /healthz and /v1/webhook/hangouts-chat endpoints.

Example:

	# first terminal
	alertmanager-webhook-forwarder serve

	# second terminal
	curl http://localhost:8080/healthz
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		listenHost := viper.GetString("server.host")
		listenPort := viper.GetInt("server.port")
		gracefulShutdownTimeout := viper.GetInt("server.gracefulShutdownTimeout")
		accessLogEnabled := viper.GetBool("server.accessLog.enabled")
		listenAddress := fmt.Sprintf("%s:%d", listenHost, listenPort)

		// Listen for shutdown signals
		done := make(chan bool, 1)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Configure routes
		router := http.NewServeMux()
		router.HandleFunc("/healthz", handleHealthz)
		router.HandleFunc("/v1/webhook/hangouts-chat", providerWebhookHandlerFactory(utils.ProviderHangoutsChat))

		var handler http.Handler

		// Optionally configure Apache-style access logs
		if accessLogEnabled {
			handler = handlers.LoggingHandler(os.Stdout, router)
		} else {
			handler = router
		}

		// Support HTTP 2 (no TLS)
		// Note: Graceful Shutdown does not work on h2c (https://github.com/golang/go/issues/26682)
		http2server := &http2.Server{}
		handler = h2c.NewHandler(handler, http2server)

		server := &http.Server{
			Addr:         listenAddress,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}

		go gracefulShutdown(server, quit, done, gracefulShutdownTimeout)

		klog.Infof("Server is listening on http://%s:%d", listenHost, listenPort)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.Fatalf("Could not listen on %s: %v\n", listenAddress, err)
		}

		<-done
		klog.Info("Server has been stopped")

		return nil
	},
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "OK")
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

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool, timeoutSeconds int) {
	<-quit
	klog.Infoln("Server is gracefully shutting down...")

	// Cancel gracefulShutdown after `timeoutSeconds` period
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds) * time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		klog.Fatalf("Could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

func init() {
	ServeCmd.Flags().String("host", "0.0.0.0", "Host that HTTP server should bind to.")
	ServeCmd.Flags().Int("port", 8080, "Port that HTTP server should listen on.")
	ServeCmd.Flags().Int("grace", 30, "Graceful shutdown timeout seconds.")
	ServeCmd.Flags().Bool("access-log", false, "Enable Apache-style access logs.")

	RootCmd.AddCommand(ServeCmd)

	_ = viper.BindPFlag("server.host", ServeCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", ServeCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("server.gracefulShutdownTimeout", ServeCmd.Flags().Lookup("grace"))
	_ = viper.BindPFlag("server.accessLog.enabled", ServeCmd.Flags().Lookup("access-log"))
}
