package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/niwho/elasticsearch"
)

const (
	DateFmt = "_%d_%02d_%02d"
)

var (
	TdyIndex  string
	IndexLock sync.RWMutex
	EsClient  elasticsearch.EsClientV3

	// config
	Hosts           []string
	Index           string
	Type            string
	ShardNum        int
	ReplicaNum      int
	RefreshInterval int
)

func InitEsClient(hosts []string, index, ttype string, shardNum, replicaNum, refreshInterval int) {
	Hosts = hosts
	Index = index
	Type = ttype
	ShardNum = shardNum
	ReplicaNum = replicaNum
	RefreshInterval = refreshInterval
	var err error
	EsClient, err = elasticsearch.CreateEsClientV3(Hosts)
	if err != nil {
		panic(err)
	}
	GetIndex()
}

func initEsIndex(index string) {

	err := EsClient.CreateEsIndex(
		index,
		int32(ShardNum),
		int32(ReplicaNum),
		int32(RefreshInterval),
	)
	if err != nil {
		panic(err)
	}
}

func GetIndex() string {

	t := time.Now()
	currentIndex := Index + fmt.Sprintf(DateFmt, t.Year(), t.Month(), t.Day())
	if TdyIndex != currentIndex {

		IndexLock.Lock()
		// 二次判断
		if TdyIndex != currentIndex {
			initEsIndex(currentIndex)
			TdyIndex = currentIndex
		}
		IndexLock.Unlock()
	}
	return TdyIndex
}

func Insert2Es(dat string) error {
	return EsClient.Insert(GetIndex(), Type, dat)
}

func Insert2EsBulk(dats []interface{}) error {
	return EsClient.BulkInsert(GetIndex(), Type, dats)
}
