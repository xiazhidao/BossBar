package utils

import (
	"github.com/astaxie/beego"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
	"time"
)

const LogPath = "logs/"
const FileSuffix = ".log"

func InitLogs() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceQuote:      true,
		ForceColors:     true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetReportCaller(true)

	runmode := strings.TrimSpace(strings.ToLower(beego.AppConfig.String("runmode")))
	if runmode == "" {
		runmode = "dev"
	}

	//写文本日志
	if runmode == "dev" {
		newfileLogger(log.StandardLogger(), LogPath, 8)
	}
}

/**
  文件日志
*/
func newfileLogger(logger *log.Logger, logPath string, save uint) {

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer(logPath, "debug", save), // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer(logPath, "info", save),
		log.WarnLevel:  writer(logPath, "warn", save),
		log.ErrorLevel: writer(logPath, "error", save),
		log.FatalLevel: writer(logPath, "fatal", save),
		log.PanicLevel: writer(logPath, "panic", save),
	}, &log.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.AddHook(lfHook)
}

/**
文件设置
*/
func writer(logPath string, level string, save uint) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	fileSuffix := time.Now().In(cstSh).Format("2006-01-02") + FileSuffix

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		//rotatelogs.WithLinkName(logFullPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(int(save)),   // 文件最大保存份数
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	return logier
}
