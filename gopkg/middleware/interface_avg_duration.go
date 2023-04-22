package middleware

/*
** 每个接口每分钟平均耗时
 */

import (
	"fmt"
	"go/gopkg/logger/vglog"
	"runtime"
	"sync"
	"time"
)

var (
	interfaceAvgDurationInstall *InterfaceAvgDuration
)

type InterfaceDetail struct {
	cnt      int64
	duration int64
}

type InterfaceAvgDuration struct {
	lock    sync.Mutex
	details map[string]*InterfaceDetail
}

//init
func InterfaceAverageDurationInit() {
	interfaceAvgDurationInstall = &InterfaceAvgDuration{
		details: make(map[string]*InterfaceDetail),
	}

	go runCalAvgTimer()
}

func AddInterfaceAccessDuration(interfaceName string, duration int64) {
	if interfaceAvgDurationInstall != nil {
		interfaceAvgDurationInstall.Add(interfaceName, duration)
	}
}

//记录接口耗时
func (i *InterfaceAvgDuration) Add(interfaceName string, duration int64) bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	_, ok := i.details[interfaceName]
	if false == ok {
		cur := &InterfaceDetail{
			cnt:      1,
			duration: duration,
		}
		i.details[interfaceName] = cur
	} else {
		i.details[interfaceName].cnt += 1
		i.details[interfaceName].duration += duration
	}

	return true
}

//计算平均耗时
func (i *InterfaceAvgDuration) calAvg() {
	i.lock.Lock()
	defer i.lock.Unlock()

	numGo := runtime.NumGoroutine()

	for k, v := range i.details {
		var avg int64
		if v.cnt == 0 {
			avg = 0
		} else {
			avg = v.duration / v.cnt
		}

		s := fmt.Sprintf("{\"interface\":\"%s\", \"cnt\":%d, \"duration\":%d, \"avg\":%d, \"numGoroutine\":%d}", k, v.cnt, v.duration, avg, numGo)
		vglog.InterfaceAverageDuration("%s", s)
		delete(i.details, k)
	}
}

//定时每分钟计算平均耗时
func runCalAvgTimer() {
	for range time.NewTicker(time.Second * 60).C {
		interfaceAvgDurationInstall.calAvg()
	}
}
