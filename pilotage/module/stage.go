package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/model"
	. "github.com/logrusorgru/aurora"
)

const (
	// Stage Type
	StartStage  = "start"
	EndStage    = "end"
	NormalStage = "normal"
	PauseStage  = "pause"
)

// Stage is
type Stage struct {
	ID         int64    `json:"-" yaml:"-"`
	T          string   `json:"type" yaml:"type"`
	Name       string   `json:"name" yaml:"name"`
	Title      string   `json:"title" yaml:"title"`
	Sequencing string   `json:"sequencing,omitempty" yaml:"sequencing,omitempty"`
	Status     string   `json:"status,omitempty" yaml:"status,omitempty"`
	Logs       []string `json:"logs,omitempty" yaml:"logs,omitempty"`
	Actions    []Action `json:"actions,omitempty" yaml:"actions,omitempty"`
}

// TODO filter the log print with different color.
func (s *Stage) Log(log string, verbose, timestamp bool) {
	s.Logs = append(s.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))
	l := new(model.LogV1)
	l.Create(model.INFO, model.STAGE, s.ID, log)

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (s *Stage) SequencingRun(verbose, timestamp bool, f *Flow, stageIndex int) (string, error) {
	s.Status = Running

	s.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), false, timestamp)
	f.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), verbose, timestamp)

	// Save Stage into database
	stage := new(model.StageV1)
	stageID, err := stage.Put(f.ID, s.T, s.Name, s.Title, s.Sequencing)
	if err != nil {
		s.Log(fmt.Sprintf("Save Stage [%s] error: %s", s.Name, err.Error()), false, timestamp)
	}
	s.ID = stageID

	// Record stage data
	stageData := new(model.StageDataV1)
	startTime := time.Now()

	for i, _ := range s.Actions {
		action := &s.Actions[i]

		s.Log(fmt.Sprintf("The Number [%d] action is running: %s", i, s.Title), false, timestamp)
		f.Log(fmt.Sprintf("The Number [%d] action is running: %s", i, s.Title), verbose, timestamp)

		if status, err := action.Run(verbose, timestamp, f, stageIndex, i); err != nil {
			s.Status = Failure

			s.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), false, timestamp)
			f.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), verbose, timestamp)

		} else {
			s.Status = status
		}

		if s.Status == Failure || s.Status == Cancel {
			break
		}
	}

	currentNumber, err := stageData.GetNumbers(stageID)
	if err != nil {
		s.Log(fmt.Sprintf("Get Stage Data [%s] Numbers error: %s", s.Name, err.Error()), verbose, timestamp)
	}
	if err := stageData.Put(s.ID, currentNumber+1, s.Status, startTime, time.Now()); err != nil {
		s.Log(fmt.Sprintf("Save Stage Data [%s] error: %s", s.Name, err.Error()), false, timestamp)
	}

	return s.Status, nil
}

func (s *Stage) ParallelRun(verbose, timestamp bool, f *Flow, stageIndex int) (string, error) {
	s.Status = Running

	s.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), false, timestamp)
	f.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), verbose, timestamp)

	// Save Stage into database
	stage := new(model.StageV1)
	stageID, err := stage.Put(f.ID, s.T, s.Name, s.Title, s.Sequencing)
	if err != nil {
		s.Log(fmt.Sprintf("Save Stage [%s] error: %s", s.Name, err.Error()), false, timestamp)
	}
	s.ID = stageID

	// Record stage data
	stageData := new(model.StageDataV1)
	startTime := time.Now()

	resultChan := make(chan string)
	defer close(resultChan)
	count := 0
	for i, _ := range s.Actions {
		go func(index int) {
			action := &s.Actions[index]
			var tempStaus string

			s.Log(fmt.Sprintf("The Number [%d] action is running: %s", index, s.Title), false, timestamp)
			f.Log(fmt.Sprintf("The Number [%d] action is running: %s", index, s.Title), verbose, timestamp)

			if status, err := action.Run(verbose, timestamp, f, stageIndex, index); err != nil {
				tempStaus = Failure

				s.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), false, timestamp)
				f.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), verbose, timestamp)
			} else {
				tempStaus = status
			}
			resultChan <- tempStaus

		}(i)
	}

	for {
		if result, ok := <-resultChan; ok {
			count++
			s.Status = result
			if result == Failure || result == Cancel || count == len(s.Actions) {

				currentNumber, err := stageData.GetNumbers(stageID)
				if err != nil {
					s.Log(fmt.Sprintf("Get Stage Data [%s] Numbers error: %s", s.Name, err.Error()), verbose, timestamp)
				}
				if err := stageData.Put(s.ID, currentNumber+1, s.Status, startTime, time.Now()); err != nil {
					s.Log(fmt.Sprintf("Save Stage Data [%s] error: %s", s.Name, err.Error()), false, timestamp)
				}

				return s.Status, nil
			}
		}

	}
}
