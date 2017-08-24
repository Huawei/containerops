package model

type LogV1 struct {
	ID    string
	Level string
	//Phase must be one of 'flow','stage','action' or 'job'
	Phase   string
	PhaseID string
	Content  string
}
