package call

import "testing"

func TestCmd(t *testing.T) {
	cmds := "go version"
	result, err := Cmd(cmds)
	if err != nil {
		t.Errorf("run cmd error: %v", err)
	} else {
		t.Logf("run cmd result: %s", result)
	}
}
