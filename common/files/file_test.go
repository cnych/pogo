package files

import "testing"

func TestWhich(t *testing.T) {
	Which("ffmpeg")
}


func TestXOk(t *testing.T) {
	t.Log(XOk("/usr/local/homebrew/bin/ffmpeg"))
	//t.Log(XOk("/Users/ych/devs/workspace/yidianzhishi/course/pogo/main.go"))
}