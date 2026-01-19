package internal

import (
	"errors"
	"fmt"
	"log/syslog"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type ConfigOverrides struct {
	Debug bool
	Log   string
}

func (o *ConfigOverrides) AsString() string {
	var out []string
	out = append(out, fmt.Sprintf("debug:'%t'", o.Debug))
	out = append(out, fmt.Sprintf("log:'%s'", o.Log))

	return fmt.Sprintf("%s", strings.Join(out, ","))
}

const (
	LOGTYPE_STDOUT  = 1
	LOGTYPE_FILE    = 2
	LOGTYPE_UNKNOWN = -1
)

func LogTypeAsString(logtype int) string {
	types := map[int]string{
		LOGTYPE_STDOUT:  "log+stdout",
		LOGTYPE_FILE:    "log+file",
		LOGTYPE_UNKNOWN: "log+unknown",
	}

	if _, ok := types[logtype]; !ok {
		return types[LOGTYPE_UNKNOWN]
	}
	return types[logtype]
}

type TLogOptions struct {
	Format string
	Level  string
	Log    string
}

func (t *TLogOptions) AsString() string {
	var out []string
	out = append(out, fmt.Sprintf("level:'%s'", t.Level))
	out = append(out, fmt.Sprintf("log:'%s'", t.Log))
	return fmt.Sprintf("%s", strings.Join(out, ","))
}

type TLog struct {
	// LogType could be LOG_STDOUT or LOG_FILE
	LogType int

	Log        *logrus.Logger
	LogOptions *TLogOptions

	File *os.File

	syslog *syslog.Writer
}

func (log *TLog) Debug(str string) {
	if log.LogType == LOGTYPE_STDOUT || log.LogType == LOGTYPE_FILE {
		if log.Log != nil {
			log.Log.Debug(str)
		}
	}
}

func (log *TLog) Info(str string) {
	if log.LogType == LOGTYPE_STDOUT || log.LogType == LOGTYPE_FILE {
		if log.Log != nil {
			log.Log.Info(str)
		}

	}
}

func (log *TLog) Error(str string) {
	if log.LogType == LOGTYPE_STDOUT || log.LogType == LOGTYPE_FILE {
		if log.Log != nil {
			log.Log.Error(str)
		}
	}
}

func (log *TLog) SetLogLevel(level string) error {
	if log.Log == nil {
		return errors.New("log not initialized")
	}
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	log.Log.Level = l
	return nil
}

func (log *TLog) SetLogOut(out string) error {
	if log.Log == nil {
		return errors.New("log not initialized")
	}

	if log.LogType == LOGTYPE_STDOUT || log.LogType == LOGTYPE_FILE {
		switch out {
		case "stdout":
			log.Log.Out = os.Stdout
		case "stderr":
			log.Log.Out = os.Stderr
		default:
			f, err := os.OpenFile(log.LogOptions.Log,
				os.O_RDWR|os.O_CREATE|os.O_APPEND,
				0644)
			if err != nil {
				return err
			}
			log.Log.Out = f
			log.File = f

			logrus.SetOutput(log.File)
			logrus.SetLevel(logrus.DebugLevel)
		}
	}
	return nil
}

type LogrotateOptions struct {
	// if we need move file to old
	// named file or not, if not
	// we just create new file and use it
	Move bool `json:"move"`
}

func (l *LogrotateOptions) AsString() string {
	return fmt.Sprintf("move:'%t'", l.Move)
}

func (log *TLog) RotateLog(LogrotateOptions *LogrotateOptions) error {
	var err error

	id := "(rotate)"
	name := log.LogOptions.Log
	log.Debug(fmt.Sprintf("%s request to rotate log:'%s' logtype:'%s'", id, name, LogTypeAsString(log.LogType)))

	if LogrotateOptions != nil {
		log.Debug(fmt.Sprintf("%s logtate option: %s", id, LogrotateOptions.AsString()))
	}

	if log.Log == nil {
		return errors.New("log not initialized")
	}

	if log.LogType == LOGTYPE_FILE || name != "stdout" {
		f := log.File
		if f == nil {
			log.Error(fmt.Sprintf("%s error file description is not found", id))
			return errors.New("log not initialized")
		}
		if err := f.Close(); err != nil {
			log.Error(fmt.Sprintf("%s error closing file, err:'%s'", id, err))
			return err
		}

		if LogrotateOptions != nil && LogrotateOptions.Move {
			rotatename := fmt.Sprintf("%s.old", name)

			err := os.Rename(name, rotatename)
			if err != nil {
				return err
			}
		}

		log.SetLogOut(log.LogOptions.Log)
	}
	return err
}

func (log *TLog) UpdateLog() error {
	var err error
	if log.LogOptions == nil {
		return errors.New("log file not found")
	}
	if log.LogType == LOGTYPE_STDOUT || log.LogType == LOGTYPE_FILE {
		if err = log.SetLogLevel(log.LogOptions.Level); err != nil {
			return err
		}
		if err = log.SetLogOut(log.LogOptions.Log); err != nil {
			return err
		}
	}
	return err
}

func (log *TLog) UpdateOverrides(Overrides ConfigOverrides) error {
	if log.LogOptions == nil {
		return errors.New("log file not found")
	}
	if log.LogType == LOGTYPE_STDOUT || log.LogType == LOGTYPE_FILE {

		if Overrides.Debug {
			log.LogOptions.Level = "debug"
		}
		if len(Overrides.Log) > 0 {
			log.LogOptions.Log = Overrides.Log
		}
	}
	return log.UpdateLog()
}

func CreateLog(opts *ConfYaml) (*TLog, error) {
	var err error

	var Log TLog
	Log.Log = logrus.New()

	LogOptions := opts.LogOptions

	Log.LogType = LOGTYPE_STDOUT
	if LogOptions.Log != "stdout" {
		Log.LogType = LOGTYPE_FILE
	}

	Log.LogOptions = LogOptions
	Log.Log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006/01/02 - 15:04:05.000",
		FullTimestamp:   true,
	}

	if err = Log.UpdateLog(); err != nil {
		return nil, errors.New(fmt.Sprintf("log error, err:'%s'", err.Error()))
	}

	return &Log, nil
}
