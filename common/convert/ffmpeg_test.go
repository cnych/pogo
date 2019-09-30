package convert

import "testing"

func TestFFMpeg_Duration(t *testing.T) {
	ffm := NewFFMpeg()
	duration, err := ffm.Duration("https://aweme.snssdk.com/aweme/v1/playwm/?s_vid=93f1b41336a8b7a442dbf1c29c6bbc5616c02e43ab032727a15dbaf57c04d21b623bb929f61b793fae3b726c7a867c8a6f1ca31d07b1c80a9f290ebc17899ee4&line=0")
	if err != nil {
		t.Errorf("get duration error: %v\n", err)
	} else {
		t.Logf("duration=%d\n", duration)
	}
}