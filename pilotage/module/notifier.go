package module

import (
	"fmt"
)

var Notifiers = make(map[string]Notifier)

type Notifier interface {
	Notify(flow *Flow, receivers []string) error
}

func Register(name string, notifier Notifier) error {
	if _, ok := Notifiers[name]; ok {
		return fmt.Errorf("Notifier %s already exist", name)
	}
	Notifiers[name] = notifier
	return nil
}
