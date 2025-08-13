package gosctx

import (
	"flag"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/teoit/gosctx/configs"
	"gopkg.in/natefinch/lumberjack.v2"
)

var pathLogDefault = "logs"

type appLoggerDaily struct {
	id         string
	logger     *logrus.Logger
	logLevel   string
	pathLog    string
	logDaily   bool
	maxSize    int
	maxAge     int
	maxBackups int
	compress   bool
}

type AppLoggerDaily interface {
	GetLogger(id string) Logger
}

func NewAppLoggerDaily(id string) *appLoggerDaily {

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = logrus.Formatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return "", fileName
		},
	})

	return &appLoggerDaily{
		id:       id,
		logger:   logger,
		logDaily: true,
		pathLog:  pathLogDefault,
	}
}

func (logDaily *appLoggerDaily) ID() string {
	return logDaily.id
}

func (logDaily *appLoggerDaily) InitFlags() {
	logDaily.logLevel = configs.LogLevel
	flag.BoolVar(&logDaily.logDaily, "log-daily", false, "Log Daily default false")
	flag.StringVar(&logDaily.pathLog, "log-path", pathLogDefault, "Log path file: default path=logs")
	flag.IntVar(&logDaily.maxSize, "max-size", 100, "MaxSize is the maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.")
	flag.IntVar(&logDaily.maxAge, "max-age", 5, "MaxAge is the maximum number of days to retain old log files based")
	flag.IntVar(&logDaily.maxBackups, "max-backups", 5, "MaxBackups is the maximum number of old log files to retain")
	flag.BoolVar(&logDaily.compress, "compress", false, "Compress determines if the rotated log files should be compressed using gzip")
}

func (logDaily *appLoggerDaily) Activate(_ ServiceContext) error {
	lv := mustParseLevel(logDaily.logLevel)

	logDaily.logger.SetLevel(lv)

	return nil
}

func (logDaily *appLoggerDaily) Stop() error {
	return nil
}

func (logDaily *appLoggerDaily) GetLogger(name string) Logger {
	var entry *logrus.Entry
	if name == "" {
		name = "log"
	}

	name = strings.Trim(name, ".")
	entry = logDaily.logger.WithField("prefix", name)

	if logDaily.logDaily && logDaily.logLevel != "info" {
		now := time.Now().Local()
		dateString := now.Format("2006-01-02")
		filename := fmt.Sprintf("%s/log_%s.log", name, dateString)
		pathLog := logDaily.pathLog + "/" + filename

		lumberjackLogger := lumberjack.Logger{
			Filename:   pathLog,
			MaxSize:    logDaily.maxSize,
			MaxBackups: logDaily.maxBackups,
			MaxAge:     logDaily.maxAge,
			Compress:   logDaily.compress,
		}
		logDaily.logger.Out = &lumberjackLogger
		l := &logger{entry}
		var log Logger = l
		return log
	}

	return &logger{entry}
}
