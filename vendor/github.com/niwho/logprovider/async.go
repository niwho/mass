package logprovier

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

type IPROVIDER interface {
	io.Writer
	Flush() error
	Close() error
}

// AsyncFrame 异步写入, 考虑更通用的方式
type AsyncFrame struct {
	buf       chan []byte
	stop      chan struct{}
	flush     chan struct{}
	providers []IPROVIDER
	workerNum int
	isRunning int32
	stopFlag  bool
	idleNum   int32
}

func NewAsyncFrame(workerNum int, prds ...IPROVIDER) *AsyncFrame {
	af := &AsyncFrame{
		buf:       make(chan []byte, 1024),
		stop:      make(chan struct{}),
		flush:     make(chan struct{}),
		providers: prds,
		workerNum: workerNum,
	}
	if workerNum <= 0 {
		af.workerNum = 1
	}
	af.Run()

	return af
}

// cleanBuf
func (af *AsyncFrame) cleanBuf() {
	for {
		select {
		case dat := <-af.buf:
			for _, provider := range af.providers {
				provider.Write(dat)
			}
		default:
			return
		}
	}
}

func (af *AsyncFrame) Run() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 0, 1) {
		return
	}
	if af.providers == nil {
		fmt.Fprintln(os.Stderr, "logger's providers is nil.")
		// return // if return, af.Stop() will be blocked because no goroutine consumes af.stop
	}
	for i := 0; i < af.workerNum; i++ {
		go af.runOuter()
	}
}

func (af *AsyncFrame) runOuter() {
	for {
		if af.stopFlag {
			return
		}
		af.runrun()
	}

}

func (af *AsyncFrame) runrun() {
	atomic.AddInt32(&af.idleNum, 1)
	defer func() {
		atomic.AddInt32(&af.idleNum, -1)
		if err := recover(); err != nil {
			const size = 64 << 20
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("AsyncFrame panic=%v\n%s\n", err, buf)
		}
	}()
	for {
		select {
		case dat, ok := <-af.buf:
			if !ok {
				fmt.Fprintln(os.Stderr, "buf channel has been closed.")
				return
			}
			for _, provider := range af.providers {
				provider.Write(dat)
			}
		case <-af.flush:
			af.cleanBuf()
			for _, provider := range af.providers {
				provider.Flush()
			}
			af.flush <- struct{}{}
		case <-af.stop:
			af.stopFlag = true
			af.cleanBuf()
			for _, provider := range af.providers {
				provider.Flush()
				provider.Close()
			}
			af.stop <- struct{}{}
			return
		}
	}

}

func (af *AsyncFrame) Write(p []byte) (int, error) {
	select {
	case af.buf <- p:
		return 0, nil
	default:
		// warn write loss
		return 0, fmt.Errorf("%s", "AsyncFrame buf overflow")
	}
}

// Close safe clean
func (af *AsyncFrame) Close() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 1, 0) {
		return
	}
	af.stop <- struct{}{}
	for {
		if atomic.LoadInt32(&af.idleNum) == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

}

// Flush
func (af *AsyncFrame) Flush() {
	if atomic.LoadInt32(&af.isRunning) == 0 {
		return
	}
	af.flush <- struct{}{}
	<-af.flush
}
