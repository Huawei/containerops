package module

import (
	"fmt"

	"github.com/Huawei/containerops/pilotage/models"

	log "github.com/Sirupsen/logrus"
)

func RecordOutcom(pipelineId, fromPiipelineId, stageId, fromStageId, actionId, fromActionId, sequence, evnetId int64, status bool, result, output string) error {
	outcome := new(models.Outcome)
	outcome.Pipeline = pipelineId
	outcome.RealPipeline = fromPiipelineId
	outcome.Stage = stageId
	outcome.RealStage = fromStageId
	outcome.Action = actionId
	outcome.RealAction = fromActionId
	outcome.Sequence = sequence
	outcome.Event = evnetId
	outcome.Status = status
	outcome.Result = result
	outcome.Output = output

	err := outcome.GetOutcome().Save(outcome).Error
	if err != nil {
		log.Error("[outcome's RecordOutcome]:error when save outcome info:", fmt.Sprintf("%#v", outcome), " ===>error is:", err.Error())
		return err
	}

	return nil
}
