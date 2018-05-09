/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/seeleteam/monitor-api/config"
)

const (
	defaultLogPath = "monitor-api-logs"
	defaultLogFile = "monitor-api-all.logs"
)

type MonitorLog struct {
	log *logrus.Logger
}

var logMap map[string]*MonitorLog
var getLogMutex sync.Mutex

// Debug wrapper Debug logger
func (m *MonitorLog) Debug(f interface{}, args ...interface{}) {
	m.log.Debug(formatLog(f, args...))
}

// Info wrapper Info logger
func (m *MonitorLog) Info(f interface{}, args ...interface{}) {
	m.log.Info(formatLog(f, args...))
}

// Warn wrapper Warn logger
func (m *MonitorLog) Warn(f interface{}, args ...interface{}) {
	m.log.Warn(formatLog(f, args...))
}

// Printf wrapper Printf logger
func (m *MonitorLog) Printf(f interface{}, args ...interface{}) {
	m.log.Print(formatLog(f, args...))
}

// Panic wrapper Panic logger
func (m *MonitorLog) Panic(f interface{}, args ...interface{}) {
	m.log.Panic(formatLog(f, args...))
}

// Fatal wrapper Fatal logger
func (m *MonitorLog) Fatal(f interface{}, args ...interface{}) {
	m.log.Fatal(formatLog(f, args...))
}

// Error wrapper Error logger
func (m *MonitorLog) Error(f interface{}, args ...interface{}) {
	m.log.Error(formatLog(f, args...))
}

// Debugln wrapper Debugln logger
func (m *MonitorLog) Debugln(v ...interface{}) {
	m.log.Debugln(v...)
}

// Infoln wrapper Infoln logger
func (m *MonitorLog) Infoln(args ...interface{}) {
	m.log.Infoln(args...)
}

// Warnln wrapper Warnln logger
func (m *MonitorLog) Warnln(args ...interface{}) {
	m.log.Warnln(args...)
}

// Printfln wrapper Printfln logger
func (m *MonitorLog) Printfln(args ...interface{}) {
	m.log.Println(args...)
}

// Panicln wrapper Panicln logger
func (m *MonitorLog) Panicln(args ...interface{}) {
	m.log.Panicln(args...)
}

// Fatalln wrapper Fatalln logger
func (m *MonitorLog) Fatalln(args ...interface{}) {
	m.log.Fatalln(args...)
}

// Errorln wrapper Errorln logger
func (m *MonitorLog) Errorln(args ...interface{}) {
	m.log.Errorln(args...)
}

// GetLogger gets logrus.Logger object according to logName
// each module can have its own logger
func GetLogger(logName string, bConsole bool) *MonitorLog {
	getLogMutex.Lock()
	defer getLogMutex.Unlock()
	if logMap == nil {
		logMap = make(map[string]*MonitorLog)
	}
	curLog, ok := logMap[logName]
	if ok {
		return curLog
	}

	log := logrus.New()

	// get logLevel
	logLevel := config.SeeleConfig.ServerConfig.LogLevel
	log.SetLevel(logLevel)

	writeLog := config.SeeleConfig.ServerConfig.EngineConfig.WriteLog
	if writeLog && !bConsole {
		err := os.MkdirAll(defaultLogPath, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("creating log file failed: %s", err.Error()))
		}

		path := defaultLogPath + string(os.PathSeparator) + defaultLogFile
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),       // 24 hours
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second), // 1 days
		)
		if err != nil {
			panic(fmt.Sprintf("rotatelogs log failed: %s", err.Error()))
			return nil
		}

		log.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
			},
			&logrus.JSONFormatter{},
		))

		pathMap := lfshook.PathMap{
			logrus.DebugLevel: fmt.Sprintf("%s/debug.log", defaultLogPath),
			logrus.InfoLevel:  fmt.Sprintf("%s/info.log", defaultLogPath),
			logrus.WarnLevel:  fmt.Sprintf("%s/warn.log", defaultLogPath),
			logrus.ErrorLevel: fmt.Sprintf("%s/error.log", defaultLogPath),
			logrus.FatalLevel: fmt.Sprintf("%s/fatal.log", defaultLogPath),
		}
		log.AddHook(lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{},
		))
	} else {
		log.Out = os.Stdout
	}

	log.AddHook(&CallerHook{}) // add caller hook to print caller's file and line number
	curLog = &MonitorLog{
		log: log,
	}
	logMap[logName] = curLog
	return curLog
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}
