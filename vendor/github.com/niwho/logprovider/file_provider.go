package logprovier

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DayHourClock struct {
	stop chan struct{}
}

func NewDayHourTicker(dur SegDuration) <-chan time.Time {
	hourClock := &DayHourClock{stop: make(chan struct{}, 1)}
	return hourClock.C(dur)
}

func (hc *DayHourClock) C(dur SegDuration) <-chan time.Time {
	ch := make(chan time.Time)
	go func(dur SegDuration) {
		var checkPoint, currentPoint int
		switch dur {
		case HourDur:
			checkPoint = time.Now().Hour()
		case DayDur:
			checkPoint = time.Now().Day()
		default:
			return // just exit to block
		}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				if dur == HourDur {
					currentPoint = t.Hour()
				} else if dur == DayDur {
					currentPoint = t.Day()
				}
				if currentPoint != checkPoint {
					ch <- t
					checkPoint = currentPoint
				}
			case <-hc.stop:
				return
			}
		}
	}(dur)
	return ch
}

func (hc *DayHourClock) Stop() {
	hc.stop <- struct{}{}
}

type SegDuration string

const (
	HourDur SegDuration = "Hour"
	DayDur  SegDuration = "Day"
	NoDur   SegDuration = "No"
)

type FileProvider struct {
	sync.Mutex
	enableRotate bool
	hourTicker   <-chan time.Time

	fd           *os.File
	filename     string
	dur          SegDuration
	reservedDays int

	buf chan []byte
}

func NewFileProvider(filename string, dur SegDuration, reservedDays int) *FileProvider {
	rotate := false
	if dur != NoDur {
		rotate = true
	}

	provider := &FileProvider{
		enableRotate: rotate,
		filename:     filename,
		dur:          dur,
		buf:          make(chan []byte, 16),
		hourTicker:   NewDayHourTicker(dur),
		reservedDays: reservedDays,
	}
	provider.Init()
	return provider
}

func (fp *FileProvider) Init() error {
	var (
		fd  *os.File
		err error
	)
	realFile, err := fp.timeFilename()
	if err != nil {
		return err
	}
	if env := os.Getenv("IS_PROD_RUNTIME"); len(env) == 0 {
		fd, err = os.OpenFile(realFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	} else {
		fd, err = os.OpenFile(realFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	}
	fp.fd = fd
	_, err = os.Lstat(fp.filename)
	if err == nil || os.IsExist(err) {
		os.Remove(fp.filename)
	}
	os.Symlink("./"+filepath.Base(realFile), fp.filename)
	return nil
}

func (fp *FileProvider) doCheck() error {
	fp.Lock()
	defer fp.Unlock()

	if !fp.enableRotate {
		return nil
	}

	select {
	case <-fp.hourTicker:
		if err := fp.truncate(); err != nil {
			fmt.Fprintf(os.Stderr, "truncate file %s error: %s\n", fp.filename, err)
			return err
		}
	default:
	}
	return nil
}

func (fp *FileProvider) Write(p []byte) (n int, err error) {
	fp.doCheck()
	n, err = fmt.Fprint(fp.fd, string(p))
	return
}

func (fp *FileProvider) Close() error {
	return fp.fd.Close()
}

func (fp *FileProvider) Flush() error {
	return fp.fd.Sync()
}

// 1: 拼接出新的日志文件的名字
// 2: 拷贝当前日志文件到新的文件
// 3: Truncate当前日志文件
func (fp *FileProvider) truncate() error {
	fp.fd.Sync()
	fp.fd.Close()
	return fp.Init()
}

func (fp *FileProvider) timeFilename() (string, error) {
	absPath, err := filepath.Abs(fp.filename)
	if err != nil {
		return "", err
	}
	fmtStr := "2006-01-02_15"
	if fp.dur == DayDur {
		fmtStr = "2006-01-02"
	}
	removedFile := absPath + "." + time.Now().Add(-time.Duration(fp.reservedDays)*time.Hour*24).Format(fmtStr)
	os.Remove(removedFile)
	return absPath + "." + time.Now().Format(fmtStr), nil
}
