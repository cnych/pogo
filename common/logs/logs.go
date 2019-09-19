package logs

import (
	"fmt"
	"io"
	"pogo/common/config"
	"pogo/common/logs/logs"
)

type Logs interface {
	// 设置实时log信息显示终端
	SetOutput(show io.Writer) Logs
	// 暂停输出日志
	Rest()
	// 恢复暂停状态，继续输出日志
	GoOn()
	// 按先后顺序实时截获日志，每次返回1条，normal标记日志是否被关闭
	StealOne() (level int, msg string, normal bool)
	// 正常关闭日志输出
	Close()
	// 返回运行状态，如0,"RUN"
	Status() (int, string)
	DelLogger(adaptername string) error
	SetLogger(adaptername string, config map[string]interface{}) error
	SetLevel(l int)
	// 以下打印方法除正常log输出外，若为客户端或服务端模式还将进行socket信息发送
	Debug(format string, v ...interface{})
	Informational(format string, v ...interface{})
	App(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Critical(format string, v ...interface{})
	Alert(format string, v ...interface{})
	Emergency(format string, v ...interface{})
}

type pgLog struct {
	*logs.BeeLogger
}

func (log *pgLog) SetOutput(show io.Writer) Logs {
	log.BeeLogger.SetLogger("console", map[string]interface{}{
		"writer": show,
		"level":  config.LogConsoleLevel,
	})
	return log
}

var Log = func() Logs {
	//p, _ := path.Split(config.LogPath)
	//// 不存在目录时创建目录
	//d, err := os.Stat(p)
	//if err != nil || !d.IsDir() {
	//	if err := os.MkdirAll(p, 0777); err != nil {
	//		// Log.Error("Error: %v\n", err)
	//	}
	//}
	pglog := &pgLog{
		BeeLogger: logs.NewLogger(config.LogCap, config.LogFeedbackLevel),
	}

	// 是否打印行信息
	pglog.BeeLogger.EnableFuncCallDepth(config.LogLineInfo)
	// 全局日志打印级别（亦是日志文件输出级别）
	pglog.BeeLogger.SetLevel(config.LogLevel)
	// 是否异步输出日志
	pglog.BeeLogger.Async(config.LogAsync)
	// 设置日志显示位置
	pglog.BeeLogger.SetLogger("console", map[string]interface{}{
		"level": config.LogConsoleLevel,
	})

	// 是否保存所有日志到本地文件
	if config.LogSave {
		err := pglog.BeeLogger.SetLogger("file", map[string]interface{}{
			"filename": config.LogPath,
		})
		if err != nil {
			fmt.Printf("日志文档创建失败：%v", err)
		}
	}

	return pglog
}()