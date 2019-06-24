package client

import (
	"github.com/ernesto-jimenez/httplogger"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"net/http"
	"time"
)

func NewKlogLoggedClient() *http.Client {
	return &http.Client{
		Transport: httplogger.NewLoggedTransport(http.DefaultTransport, &klogHttpLogger{}),
	}
}

type klogHttpLogger struct{}

func (l *klogHttpLogger) LogRequest(req *http.Request) {
	conditional_runner.V(7).DumpRequest(req)
}

func (l *klogHttpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	conditional_runner.V(7).DumpRequestResponse(req, res, err, duration)
}
