package forwarder

import (
	"errors"
	"fmt"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
)

type Registry struct {
	Forwarders map[string]*Forwarder
}

var mainRegistry = Registry{
	Forwarders: map[string]*Forwarder{},
}

func Attach(provider string, f *Forwarder) {
	mainRegistry.Forwarders[provider] = f
}

func Get(provider string) (f *Forwarder, err error) {
	if f, found := mainRegistry.Forwarders[provider]; found && f != nil {
		return f, nil
	}

	m, _ := utils.InterfaceToStringMapOfInterfaces(mainRegistry.Forwarders)
	return nil, errors.New(fmt.Sprintf("provider \"%s\" has not got any attached forwarder (currently ready providers: %s)", provider, utils.GetListOfKeys(m)))
}
