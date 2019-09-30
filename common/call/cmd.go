package call

import (
	"bytes"
	"os/exec"
)

func Cmd(cmds string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/sh", "-c", cmds)
	var domifstat bytes.Buffer
	cmd.Stdout = &domifstat
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return domifstat.String(), nil
}

