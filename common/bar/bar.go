package bar

import (
	"fmt"
	"time"
)

type Bar struct {
	Total int64
	Size int64
	Resize func(bar *Bar) error
	finish chan bool
	err chan bool
	start int64
}

func NewBar(total int64) *Bar {
	return &Bar{
		Total: total,
		finish: make(chan bool, 1),
	}
}

func (bar *Bar) Start() {
	bar.start = time.Now().Unix()
	go func() {
		for {
			if bar.Size < bar.Total {
				if err := bar.Resize(bar); err != nil {
					fmt.Printf("resize bar error: %s\n", err.Error())
					bar.err <- true
					break
				}
				time.Sleep(100 * time.Millisecond)
			} else {
				bar.finish <- true
				return
			}
		}
	}()
}

func (bar *Bar) ShowProgress() {
	t := time.Tick(100*time.Millisecond)
	for {
		select {
		case <- bar.err:
		case <- bar.finish:
			goto ForEnd
		case <- t:
			bar.print()
		}
	}
ForEnd:
	bar.print()
	fmt.Println()
}

func formatSize(size int64) string {
	s := fmt.Sprintf("%.2f MiB", float64(size) / 1024 /1024)
	return s
}

func (bar *Bar) print() {
	str := ""
	progress := bar.Size * 100 / bar.Total
	for i:=int64(0);i<int64(100);i++ {
		if i < progress{
			str += "="
		} else if i == progress {
			str += ">"
		} else {
			str += "_"
		}
	}
	cost := fmt.Sprintf("%ds", time.Now().Unix() - bar.start)
	str = fmt.Sprintf("%s/%s [%s] %d%% %s", formatSize(bar.Size), formatSize(bar.Total),
		str, progress, cost)
	fmt.Printf("\r%s", str)
}


