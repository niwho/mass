package utils

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// AsyncJob 异步写入, 考虑更通用的方式
type AsyncJob struct {
	ffs       chan func() error
	dataCh    chan interface{}
	stop      chan struct{}
	flush     chan struct{}
	stopFlag  bool
	workerNum int
	isRunning int32

	batchNum  int
	batchFunc func([]interface{})
	wait      time.Duration

	idleNum int32
}

func NewAsyncJob(workerNum, batchNum int, bf func([]interface{}), wait time.Duration) *AsyncJob {
	return NewAsyncJobA(1024, workerNum, batchNum, bf, wait)
}

func NewAsyncJobA(funChanSize, workerNum, batchNum int, bf func([]interface{}), wait time.Duration) *AsyncJob {

	if funChanSize < 1024 {
		funChanSize = 1024
	}

	af := &AsyncJob{
		ffs:       make(chan func() error, funChanSize),
		stop:      make(chan struct{}, 1),
		flush:     make(chan struct{}, 1),
		workerNum: workerNum,
		dataCh:    make(chan interface{}, batchNum*workerNum),
		batchNum:  batchNum,
	}
	af.batchFunc = bf
	if workerNum <= 0 {
		af.workerNum = 1
	}
	if af.wait <= 0 {
		af.wait = 5 * time.Second
	}
	af.Run()

	return af
}

// cleanBuf
func (af *AsyncJob) cleanBuf() {
	for {
		select {
		case ff := <-af.ffs:
			ff()
		case dt, ok := <-af.dataCh:
			if ok {
				af.batchFunc([]interface{}{dt})
			}
		default:
			return
		}
	}
}

func (af *AsyncJob) Run() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 0, 1) {
		return
	}
	for i := 0; i < af.workerNum; i++ {
		go af.runOuter()
	}
}

func (af *AsyncJob) runOuter() {
	for {
		if af.stopFlag {
			return
		}
		af.runrun()
	}

}

func (af *AsyncJob) runrun() {
	atomic.AddInt32(&af.idleNum, 1)
	defer func() {
		atomic.AddInt32(&af.idleNum, -1)
		if err := recover(); err != nil {
			const size = 64 << 20
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("AsyncJob panic=%v\n%s\n", err, buf)
		}
	}()
	var arr []interface{}
	var ticker *time.Ticker = time.NewTicker(af.wait)
	var tickerFlag bool

	for {
		select {
		case ff, ok := <-af.ffs:
			if !ok {
				fmt.Fprintln(os.Stderr, "buf channel has been closed.")
				// af.stop <- struct{}{}
				time.Sleep(time.Second)
				return
			}
			ff()
		case <-af.stop:
			af.cleanBuf()
			af.stopFlag = true
			af.stop <- struct{}{}
			return
		case dt, ok := <-af.dataCh:
			if !ok {
				fmt.Fprintln(os.Stderr, "dataCh channel has been closed.")
				// af.stop <- struct{}{}
				time.Sleep(time.Second)
				return
			}
			if af.batchFunc != nil {
				arr = append(arr, dt)
			}
		case <-ticker.C:
			tickerFlag = true
		case <-af.flush:
			tickerFlag = true

		}
		if af.batchNum > 0 && af.batchFunc != nil {
			if len(arr) > af.batchNum || (tickerFlag && len(arr) > 0) {
				af.batchFunc(arr)
				tickerFlag = false
				arr = arr[:0]
			}
		}
	}

}

func (af *AsyncJob) Flush() {
	// 忽略同一时刻的多个
	select {
	case af.flush <- struct{}{}:
	default:

	}
}

func (af *AsyncJob) AddData(dt interface{}) {
	af.dataCh <- dt
}
func (af *AsyncJob) AddJob(ff func() error) error {
	select {
	case af.ffs <- ff:
		return nil
	default:
		// warn write loss
		return fmt.Errorf("%s", "AsyncJob overflow")

	}
}

func (af *AsyncJob) AddJobNoLoss(ff func() error) {
	af.ffs <- ff
}

// Close safe clean
func (af *AsyncJob) Close() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 1, 0) {
		return
	}
	af.stop <- struct{}{}
	// <-af.stop // bug
	// 每个协程给一个编号，空闲计数
	for {
		if atomic.LoadInt32(&af.idleNum) == 0 {
			break
		}
		fmt.Println("aysnc close", af.idleNum, af.stopFlag)
		time.Sleep(100 * time.Millisecond)
	}
}
