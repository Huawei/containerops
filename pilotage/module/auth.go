package module

const (
	// AuthTypePipelineDefault is default auth type that pipeline doesn't spec a target
	AuthTypePipelineDefault = "default"
	// AuthTypePipelineStartAPI is auth that given by api request
	AuthTypePipelineStartAPI = "PipelineStartAPI"
	// AuthTypePipelineStartDone is auth that given by system when pipeline start auth is all done
	AuthTypePipelineStartDone = "PipelineStartDone"

	// AuthTypePipelineDefault is default auth type that pipeline doesn't spec a target
	AuthTypeStageDefault = "default"
	// AuthTypePreStageDone is auth that given by system when pre stage run success
	AuthTypePreStageDone = "PreStageDone"
	// AuthTyptStageStartDon is auth that given by system when a stage get all auth it request
	AuthTyptStageStartDone = "StageStartDone"

	// AuthAuthorizerDefault is default auth authorizer that pipeline doesn't spec a authorizer
	AuthTokenDefault = "default"
)
