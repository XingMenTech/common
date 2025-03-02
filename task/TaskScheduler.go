// 任务调度器
// Neo

package task

import (
	"context"
	"sync"
	"time"
)

var (
	InstanceTaskScheduler *TaskScheduler
)

// TaskScheduler /任务调度器
type TaskScheduler struct {
	sync.RWMutex
	allTriggers map[Trigger]interface{}
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
	priority    int
}

// /初始化
func init() {
	InstanceTaskScheduler = NewTaskScheduler()
}

// NewTaskScheduler /工厂方法
func NewTaskScheduler() *TaskScheduler {
	object := &TaskScheduler{
		allTriggers: make(map[Trigger]interface{}, 0),
		wg:          &sync.WaitGroup{},
		priority:    4,
	}
	object.ctx, object.cancel = context.WithCancel(context.Background())
	//协程池提交任务
	InstanceRoutinePool.PostTask(func(params []interface{}) interface{} {
		object.schedule()
		return nil
	})
	return object
}

// /循环
func (object *TaskScheduler) schedule() {
	object.wg.Add(1)
loop:
	for {
		select {
		case <-object.ctx.Done():
			break loop
		case now := <-time.After(1 * time.Second):
			object.RLock()
			for trigger := range object.allTriggers {
				if trigger.CanTrigger(now) {
					InstanceRoutinePool.PostTask(func(params []interface{}) interface{} {
						trigger := params[0].(Trigger)
						trigger.Trigger()
						return nil
					}, trigger)
					//不能周期性触发的，直接删除
					if !trigger.CanPeriodic() {
						delete(object.allTriggers, trigger)
					}
				}
			}
			object.RUnlock()
		}
	}
	object.wg.Done()
}

// AddTrigger /设置触发器
func (object *TaskScheduler) AddTrigger(trigger Trigger) {
	object.Lock()
	object.allTriggers[trigger] = nil
	object.Unlock()
}

// DeleteTrigger /删除触发器
func (object *TaskScheduler) DeleteTrigger(trigger Trigger) {
	object.Lock()
	if _, ok := object.allTriggers[trigger]; ok {
		delete(object.allTriggers, trigger)
	}
	object.Unlock()
}

// DeleteTriggers /删除多个触发器
func (object *TaskScheduler) DeleteTriggers(canDelete func(trigger Trigger) bool) {
	object.Lock()
	for k := range object.allTriggers {
		if canDelete(k) {
			delete(object.allTriggers, k)
		}
	}
	object.Unlock()
}

// Name /名字
func (object *TaskScheduler) Name() string {
	return "TaskScheduler"
}

// SetShutdownPriority /设置关闭优先级
func (object *TaskScheduler) SetShutdownPriority(priority int) {
	object.priority = priority
}

// ShutdownPriority /关闭优先级
func (object *TaskScheduler) ShutdownPriority() int {
	return object.priority
}

// BeforeShutdown /关闭之前
func (object *TaskScheduler) BeforeShutdown() {
	object.cancel()
	object.wg.Wait()
}

func (object *TaskScheduler) AfterShutdown() {

}
