package job

import "os/exec"

type CmdJob struct {
	*baseJob
}

func NewCmdJob(jobID JobID, cmd *exec.Cmd, interrupt <-chan bool, workDir string) *CmdJob {
	return &CmdJob{baseJob: newBaseJob(jobID, cmd, interrupt, workDir)}
}
