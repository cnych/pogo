package config

// 来自配置文件的配置项。
var (
	LogCap           int64  = 10000 // 日志缓存的容量
	LogLevel         int    = 8     // 全局日志打印级别（亦是日志文件输出级别）
	LogConsoleLevel  int    = 8     // 日志在控制台的显示级别
	LogFeedbackLevel int    = 7     // 客户端反馈至服务端的日志级别
	LogLineInfo      bool   = true  // 日志是否打印行信息
	LogSave          bool   = false // 是否保存所有日志到本地文件
	LogPath          string = "logs/pogo.log"
	LogAsync         bool   = false //日志是否异步输出
)

