package conditional_runner

import (
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/log"
	"github.com/spf13/viper"
	"k8s.io/klog"
	"net/http"
	"time"
)

type conditionalRunner struct {
	condition func() bool
}

func NotDryRun() *conditionalRunner {
	return &conditionalRunner{
		condition: func() bool {
			return !viper.GetBool("dryRun")
		},
	}
}

func V(verbosity int) *conditionalRunner {
	if verbosity < 0 {
		klog.Errorln("Verbosity has to be greater than 0")
	}

	return &conditionalRunner{
		condition: func() bool {
			return viper.GetInt("verbosity") >= verbosity
		},
	}
}

func (ex *conditionalRunner) RunE(cmd func() error) (err error) {
	if ex.condition() {
		return cmd()
	}

	return nil
}

func (ex *conditionalRunner) Run(cmd func()) {
	if ex.condition() {
		cmd()
	}
}

func (ex *conditionalRunner) Check() bool {
	return ex.condition()
}

func (ex *conditionalRunner) DumpMultilineText(name string, format string, multilineText string) {
	if ex.condition() {
		log.DumpMultilineTextDepth(3, name, format, multilineText)
	}
}

func (ex *conditionalRunner) DumpYaml(name string, obj interface{}) {
	if ex.condition() {
		log.DumpYamlDepth(3, name, obj)
	}
}

func (ex *conditionalRunner) DumpJson(name string, obj interface{}) {
	if ex.condition() {
		log.DumpJsonDepth(3, name, obj)
	}
}

func (ex *conditionalRunner) DumpRequestResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	log.DumpRequestResponseDepth(2, req, res, err, duration)
}
func (ex *conditionalRunner) DumpRequest(req *http.Request) {
	log.DumpRequestDepth(2, req)
}
