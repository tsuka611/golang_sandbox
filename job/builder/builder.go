package builder

import (
	"github.com/tsuka611/golang_sandbox/job"
	"github.com/tsuka611/golang_sandbox/message"
	"github.com/tsuka611/golang_sandbox/util"
	"os/exec"
)

func BuildCmdJobFromMessage(m message.CmdJob, interrupt <-chan bool, workRoot string) (*job.CmdJob, error) {
	cmd := exec.Command(m.Command, m.Args...)
	cmd.Args = m.Args
	cmd.Env = append(m.BaseEnv, m.Env...)
	cmd.Dir = m.Dir
	workDir, err := util.NewDir(workRoot, string(m.JobID))
	if err != nil {
		return nil, err
	}
	return job.NewCmdJob(m.JobID, cmd, interrupt, workDir), nil
}
