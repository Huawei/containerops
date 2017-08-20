package module

import (
	"fmt"
	"time"
	"strings"

	. "github.com/logrusorgru/aurora"
)

const (
	// Action Type
	Sequencing = "sequence"
	Parallel   = "parallel"
)


// Action is
type Action struct {
	Name   string   `json:"name" yaml:"name"`
	Title  string   `json:"title" yaml:"title"`
	Status string   `json:"status,omitempty" yaml:"status,omitempty"`
	Jobs   []Job    `json:"jobs,omitempty" yaml:"jobs,omitempty"`
	Logs   []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}


// TODO filter the log print with different color.
func (a *Action) Log(log string, verbose, timestamp bool) {
	a.Logs = append(a.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (a *Action) Run(verbose, timestamp bool, f *Flow) (string, error) {
	a.Status = Running

	a.Log(fmt.Sprintf("Action [%s] status change to %s", a.Name, a.Status), false, timestamp)
	f.Log(fmt.Sprintf("Action [%s] status change to %s", a.Name, a.Status), verbose, timestamp)

	for i, _ := range a.Jobs {
		job := &a.Jobs[i]

		a.Log(fmt.Sprintf("The Number [%d] job is running: %s", i, a.Title), false, timestamp)
		f.Log(fmt.Sprintf("The Number [%d] job is running: %s", i, a.Title), verbose, timestamp)

		if status, err := job.Run(a.Name, verbose, timestamp, f); err != nil {
			a.Status = Failure

			a.Log(fmt.Sprintf("Job [%d] run error: %s", i, err.Error()), false, timestamp)
			f.Log(fmt.Sprintf("Job [%d] run error: %s", i, err.Error()), verbose, timestamp)

		} else {
			a.Status = status
		}

		if a.Status == Failure || a.Status == Cancel {
			break
		}

	}

	return a.Status, nil
}

