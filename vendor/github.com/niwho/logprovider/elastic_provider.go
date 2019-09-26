package logprovier

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/niwho/elasticsearch"
)

const (
	DateFmt = "_%d_%02d_%02d"
)

type EsProvider struct {
	*elasticsearch.EsClient
	sync.Mutex
	ticker <-chan time.Time
	level  int

	indexBase    string
	currentIndex string
	typ          string
}

func NewEsProvider(indexBase, typ string) *EsProvider {
	return &EsProvider{
		ticker:    NewDayHourTicker(DayDur),
		indexBase: indexBase,
		typ:       typ,
	}
}

func (ep *EsProvider) Init() (err error) {
	t := time.Now()
	needIndex := ep.indexBase + fmt.Sprintf(DateFmt, t.Year(), t.Month(), t.Day())
	if needIndex != ep.currentIndex {
		err = ep.CreateEsIndex(needIndex, 2, 0, 60)
		if err == nil {
			ep.currentIndex = needIndex
		}
	}

	return
}

func (ep *EsProvider) doCheck() error {
	// 不需要加锁
	//	ep.Lock()
	//	defer ep.Unlock()
	select {
	case <-ep.ticker:
		if err := ep.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "create index error: %s\n", err)
			return err
		}
	default:
	}
	return nil
}

func (ep *EsProvider) SetLevel(l int) {
	ep.level = l
}

func (ep *EsProvider) WriteMsg(msg string, level int) error {
	if level < ep.level {
		return nil
	}
	ep.doCheck()
	return ep.Insert(ep.currentIndex, ep.typ, msg)
}

func (ep *EsProvider) Flush() error {
	return nil
}

func (ep *EsProvider) Close() error {
	return nil
}
