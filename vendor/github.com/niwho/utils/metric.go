package utils

import (
	"fmt"
	"net"
)

const METRIC_SERVICE_DEFAULT = "161.202.208.20:18694"

var metricConn net.Conn
var SERVICE string

func InitMetrics(service, metricSrv string) (err error) {
	SERVICE = service
	if metricSrv == "" {
		metricSrv = METRIC_SERVICE_DEFAULT
	}
	metricConn, err = net.Dial("udp", metricSrv)
	if err != nil {
		return err
	}
	return nil
}

func commonTags(tags map[string]string) {
	if _, found := tags["host"]; !found {
		tags["host"] = LookupIP()
	}
}

func MetricCounter(tags map[string]string, cnt int64) {
	tagStr := ""

	for k, v := range tags {
		tagStr += fmt.Sprintf(",%s=%s", k, v)
	}
	data := fmt.Sprintf("%s%s  cnt=%d", SERVICE, tagStr, cnt)
	metricConn.Write([]byte(data))
}

//单位毫秒
func MetricTimer(tags map[string]string, latency int64) {
	tagStr := ""
	for k, v := range tags {
		tagStr += fmt.Sprintf(",%s=%s", k, v)
	}
	data := fmt.Sprintf("%s%s latency=%d", SERVICE, tagStr, latency)
	metricConn.Write([]byte(data))
}
