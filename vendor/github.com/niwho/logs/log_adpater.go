package logs

import (
	"fmt"
	"io"
	"os"
	"sync"

	logprovider "github.com/niwho/logprovider"
	log "github.com/sirupsen/logrus"
)

var LogAdapterInstance *LogAdapter

func InitLogAdapter(fileName string, reservedDays int) log.Hook {
	fileProvider := logprovider.NewFileProvider(fileName, logprovider.DayDur, reservedDays)
	asyncFrameLog := logprovider.NewAsyncFrame(1, fileProvider)
	LogAdapterInstance = &LogAdapter{
		asyncFrameLog: asyncFrameLog,
		entryPool: sync.Pool{
			New: func() interface{} {
				return &log.Entry{}
			},
		},
	}
	return LogAdapterInstance
}

type LogAdapter struct {
	asyncFrameLog *logprovider.AsyncFrame
	entryPool     sync.Pool
}

func (l *LogAdapter) GetWriter() io.Writer {
	return l.asyncFrameLog
}

func (*LogAdapter) Levels() []log.Level {
	return []log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel, log.WarnLevel, log.InfoLevel, log.DebugLevel}
}

func (la *LogAdapter) Fire(entry *log.Entry) error {
	if v, ok := entry.Logger.Formatter.(*log.TextFormatter); ok {
		v.DisableColors = true
	}
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	} else {
		_, err = la.asyncFrameLog.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
	}

	if v, ok := entry.Logger.Formatter.(*log.TextFormatter); ok {
		v.DisableColors = false
	}
	return nil
}

func (la *LogAdapter) Close() {
	la.asyncFrameLog.Close()
}
