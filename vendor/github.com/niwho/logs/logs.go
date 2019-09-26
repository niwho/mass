package logs

import (
	"fmt"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// 默认debu输出
func init() {
	log.SetLevel(log.Level(5))
	log.SetFormatter(&log.JSONFormatter{})
}

func SetFormatter(formatter log.Formatter) {
	log.SetFormatter(formatter)
}

func InitLog(fileName string, level L, reservedDays int) {
	InitLogAdapter(fileName, reservedDays)
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(level))
	log.AddHook(LogAdapterInstance)
}

type F log.Fields
type L log.Level

func commonFileds(metaFields, dataFiled F) *log.Entry {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	if metaFields == nil {
		metaFields = map[string]interface{}{}
	}
	if dataFiled == nil {
		dataFiled = map[string]interface{}{}
	}
	realField := F{
		"data":     dataFiled,
		"metadata": metaFields,
	}
	metaFields["timestamp"] = time.Now().Format("2006-01-02 15:04:05.000")
	dataFiled["pos"] = fmt.Sprintf("%s:%d", file, line)
	return log.WithFields(log.Fields(realField))
}

func Log(kvs ...map[string]interface{}) *log.Entry {
	if len(kvs) == 2 {
		return commonFileds(kvs[0], kvs[1])
	} else if len(kvs) == 1 {
		return commonFileds(nil, kvs[0])
	} else {
		return commonFileds(nil, nil)
	}
}

func WithField(key string, val interface{}) *log.Entry {
	return commonFileds(nil, F{key: val})
}

func AddHook(hook log.Hook) {
	log.AddHook(hook)
}

/**
{"count":1,"current_time":"2019-05-31 15:32:58.745","level":"warning","msg":"GetAudienceList show count","needCount":6,
"pos":"/data/golang/src/code.rightpaddle.cn/backend-golang/live/rooms/room_impl.go:252","roomId":2561665533925658753,"time":"2019-05-31T15:32:58+08:00"}
*/
