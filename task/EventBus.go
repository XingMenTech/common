// 事件总线
// Neo
// 示例
// Ping事件
// type PingEvent struct {
// 	  Text string
// }

// Ping事件
// type PingEventHandlerV1 struct {
// }

// func (object *PingEventHandlerV1) Notify(param interface{}) {
// 	  if v, ok := param.(*PingEvent); ok {
//		  singleton.LOG.Infof("v1, ping event: %v", v.Text)
//	  }
// }
// 注册事件
// singleton.InstanceEventBus.Register(reflect.TypeOf(&PingEvent{}), &PingEventHandlerV1{})
// 通知事件
// singleton.InstanceEventBus.Notify(reflect.TypeOf(&PingEvent{}), &PingEvent{Text: "ping"})

package task

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.novgate.com/common/common/logger"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

const (
	EventBusWorkerSize = 32
	NotifyChanMaxSize  = 1024
)

var (
	InstanceEventBus *EventBus
)

type EventBeforeNotifyFilter func(eventType, event interface{}) bool

// Notifiable /可通知接口
type Notifiable interface {
	// Notify /通知
	Notify(param interface{})
}

type NotifyParam struct {
	NotifiableArray []Notifiable
	Param           interface{}
}

// /事件总线
type EventBus struct {
	sync.RWMutex
	exitFlag                 int32
	priority                 int
	log                      *logrus.Entry
	ctx                      context.Context
	cancel                   context.CancelFunc
	wg                       *sync.WaitGroup
	notifyCh                 chan *NotifyParam
	eventGroup               map[interface{}][]Notifiable
	eventBeforeNotifyFilters []EventBeforeNotifyFilter
}

// /工厂方法
func NewEventBus() *EventBus {
	object := &EventBus{
		exitFlag:   0,
		log:        logger.LOG.WithField("module", "EventBus"),
		wg:         &sync.WaitGroup{},
		notifyCh:   make(chan *NotifyParam, NotifyChanMaxSize),
		eventGroup: make(map[interface{}][]Notifiable),
	}
	object.ctx, object.cancel = context.WithCancel(context.Background())
	return object
}

// /安装发送消息过滤器
func (object *EventBus) InstallBeforeNotifyFilter(filter EventBeforeNotifyFilter) {
	object.Lock()
	object.eventBeforeNotifyFilters = append(object.eventBeforeNotifyFilters, filter)
	object.Unlock()
}

// /等待事件
func (object *EventBus) wait() {
	defer func() {
		if r := recover(); nil != r {
			//打印调用栈
			debug.PrintStack()
			//退还资源
			object.wg.Done()
			//事件循环非正常退出
			InstanceRoutinePool.PostTask(func(params []interface{}) interface{} {
				object.wait()
				return nil
			})
		}
	}()
	object.wg.Add(1)
loop:
	for {
		select {
		case <-object.ctx.Done():
			break loop
		case notifyParam := <-object.notifyCh:
			for _, notifiable := range notifyParam.NotifiableArray {
				notifiable.Notify(notifyParam.Param)
			}
		}
	}
	for 0 != len(object.notifyCh) {
		notifyParam := <-object.notifyCh
		for _, notifiable := range notifyParam.NotifiableArray {
			notifiable.Notify(notifyParam.Param)
		}
	}
	object.wg.Done()
}

// /通知事件
func (object *EventBus) notify(notifiableArray []Notifiable, param interface{}) {
	object.notifyCh <- &NotifyParam{
		NotifiableArray: notifiableArray,
		Param:           param,
	}
}

// /运行
func (object *EventBus) Start() {
	for i := 0; i < EventBusWorkerSize; i++ {
		InstanceRoutinePool.PostTask(func(params []interface{}) interface{} {
			object.wait()
			return nil
		})
	}
}

// /停止
func (object *EventBus) stop() {
	atomic.StoreInt32(&object.exitFlag, 1)
	object.cancel()
	object.wg.Wait()
	for 0 != len(object.notifyCh) {
		notifyParam := <-object.notifyCh
		for _, notifiable := range notifyParam.NotifiableArray {
			notifiable.Notify(notifyParam.Param)
		}
	}
	close(object.notifyCh)
	object.log.Infof("event bus stopped")
}

// 注册
func (object *EventBus) Register(event interface{}, notifiable Notifiable) *EventBus {
	object.Lock()
	var notifiableArray []Notifiable
	if v, ok := object.eventGroup[event]; !ok {
		notifiableArray = make([]Notifiable, 0)
	} else {
		notifiableArray = v
	}
	notifiableArray = append(notifiableArray, notifiable)
	object.eventGroup[event] = notifiableArray
	object.Unlock()
	return object
}

// /通知
func (object *EventBus) Notify(event, param interface{}) *EventBus {
	if 0 != atomic.LoadInt32(&object.exitFlag) {
		fmt.Fprintf(os.Stderr, "lost notify message: (%v,%v)", event, param)
		return object
	}

	//前置过滤
	object.RLock()
	for _, filter := range object.eventBeforeNotifyFilters {
		if !filter(event, param) {
			object.RUnlock()
			return object
		}
	}
	object.RUnlock()

	var notifiableArray []Notifiable
	object.RLock()
	if v, ok := object.eventGroup[event]; ok {
		notifiableArray = make([]Notifiable, len(v))
		copy(notifiableArray, v)
	}
	object.RUnlock()
	if nil != notifiableArray && 0 != len(notifiableArray) {
		object.notify(notifiableArray, param)
	}
	return object
}

// /同步通知
func (object *EventBus) SyncNotify(event, param interface{}) *EventBus {
	if 0 != atomic.LoadInt32(&object.exitFlag) {
		fmt.Fprintf(os.Stderr, "lost notify message: (%v,%v)", event, param)
		return object
	}

	//前置过滤
	object.RLock()
	for _, filter := range object.eventBeforeNotifyFilters {
		if !filter(event, param) {
			object.RUnlock()
			return object
		}
	}
	object.RUnlock()

	var notifiableArray []Notifiable
	object.RLock()
	if v, ok := object.eventGroup[event]; ok {
		notifiableArray = make([]Notifiable, len(v))
		copy(notifiableArray, v)
	}
	object.RUnlock()
	if nil != notifiableArray && 0 != len(notifiableArray) {
		for _, notifiable := range notifiableArray {
			notifiable.Notify(param)
		}
	}

	return object
}

// /名字
func (object *EventBus) Name() string {
	return "EventBus"
}

// SetShutdownPriority /设置关闭优先级
func (object *EventBus) SetShutdownPriority(priority int) {
	object.priority = priority
}

// ShutdownPriority /关闭优先级
func (object *EventBus) ShutdownPriority() int {
	return object.priority
}

// /开始停服
func (object *EventBus) BeforeShutdown() {
	object.stop()
}

// /结束停服
func (object *EventBus) AfterShutdown() {}
