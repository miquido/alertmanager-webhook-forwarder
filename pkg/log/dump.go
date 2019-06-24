package log

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"k8s.io/klog"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func alignMultiline(content string) string {
	return strings.ReplaceAll(content, "\n", "\n ")
}

func DumpMultilineTextDepth(depth int, name string, format string, multilineText string) {
	klog.InfoDepth(depth, fmt.Sprintf("%s (%s):\n ---\n %s\n ---", name, format, alignMultiline(multilineText)))
}

func DumpMultilineText(name string, format string, multilineText string) {
	DumpMultilineTextDepth(2, name, format, multilineText)
}

func DumpYamlDepth(depth int, name string, obj interface{}) {
	parsed, err := yaml.Marshal(obj)
	if err != nil {
		klog.Error(err)
		return
	}
	DumpMultilineTextDepth(depth, name, "YAML", string(parsed))
}

func DumpYaml(name string, obj interface{}) {
	DumpYamlDepth(4, name, obj)
}

func DumpJsonDepth(depth int, name string, obj interface{}) {
	parsed, err := json.Marshal(obj)
	if err != nil {
		klog.Error(err)
	}
	DumpMultilineTextDepth(depth, name, "JSON", string(parsed))
}

func DumpJson(name string, obj interface{}) {
	DumpJsonDepth(4, name, obj)
}

func DumpRequestDepth(depth int, req *http.Request) {
	dumpedReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		klog.Error(err)
		return
	}

	klog.InfoDepth(depth, fmt.Sprintf("HTTP Request:\n ---\n %s\n ---", alignMultiline(string(dumpedReq))))
}

func DumpRequest(req *http.Request) {
	DumpRequestDepth(2, req)
}

func DumpRequestResponseDepth(depth int, req *http.Request, res *http.Response, err error, duration time.Duration) {
	if err != nil {
		klog.Error(err)
		return
	}

	dumpedRes, err := httputil.DumpResponse(res, true)
	if err != nil {
		klog.Error(err)
		return
	}

	klog.InfoDepth(depth, fmt.Sprintf(
		"HTTP Response:\n ---\n %s %s\n Duration: %sms\n Status: %s\n ---",
		req.Method,
		req.URL.String(),
		duration/time.Millisecond,
		alignMultiline(string(dumpedRes)),
	))
}

func DumpRequestResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	DumpRequestResponseDepth(2, req, res, err, duration)
}
