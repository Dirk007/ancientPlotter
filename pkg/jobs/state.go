package jobs

type JobState string

const (
	JobStateNew     JobState = "new"
	JobStateWorking JobState = "working"
	JobStateDone    JobState = "done"
	JobStateError   JobState = "error"
)
