package call

import (
	"bytes"
	"os/exec"
)

func Cmd(cmds string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/sh", "-c", cmds)
	var resultBuffer bytes.Buffer
	cmd.Stdout = &resultBuffer
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return resultBuffer.String(), nil
}
