// 协程池
// Neo

package task

import (
	"context"
	"sync"
	"sync/atomic"
)

var (
	routinePool *RoutinePool
	rpOnce      sync.Once
)

const (
	DefaultPoolSize = 128
)

// TaskMethod /任务方法
type TaskMethod func(params []interface{}) interface{}

// TaskParam /任务参数
type TaskParam struct {
	TaskMethod TaskMethod
	TaskParam  []interface{}
}

// Context /协程上下文
type Context struct {
	sync.RWMutex
	idleFlag        int32
	ctx             context.Context
	cancel          context.CancelFunc
	taskCh          chan *TaskParam
	recycleNotifyCh chan interface{}
	wg              *sync.WaitGroup
	parentWG        *sync.WaitGroup
}

// /协程方法
func (object *Context) loop() {

	object.wg.Add(1)
loop:
	for {
		select {
		case <-object.ctx.Done():
			break loop
		case taskParam, ok := <-object.taskCh:
			if !ok {
				continue
			}
			object.parentWG.Add(1)
			taskParam.TaskMethod(taskParam.TaskParam)
			object.parentWG.Done()
			atomic.StoreInt32(&object.idleFlag, 1)

			object.RLock()
			if nil != object.recycleNotifyCh &&
				0 < cap(object.recycleNotifyCh) &&
				len(object.recycleNotifyCh) < cap(object.recycleNotifyCh) {
				object.recycleNotifyCh <- nil
			}
			object.RUnlock()
		}
	}
	object.wg.Done()
}

// RoutinePool /协程池
type RoutinePool struct {
	sync.RWMutex
	minRoutine      int64
	contextPool     []*Context
	recycleNotifyCh chan interface{}
	ctx             context.Context
	cancel          context.CancelFunc
	wg              *sync.WaitGroup
	priority        int
}

// NewRoutinePool /工厂方法
func NewRoutinePool() *RoutinePool {
	rpOnce.Do(func() {
		object := &RoutinePool{
			minRoutine:      DefaultPoolSize,
			recycleNotifyCh: make(chan interface{}, 2*DefaultPoolSize),
			wg:              &sync.WaitGroup{},
			priority:        2,
		}
		object.ctx, object.cancel = context.WithCancel(context.Background())
		object.contextPool = make([]*Context, DefaultPoolSize)
		for i := int64(0); i < DefaultPoolSize; i++ {
			ctx := &Context{
				idleFlag:        1,
				taskCh:          make(chan *TaskParam, 16),
				recycleNotifyCh: object.recycleNotifyCh,
				parentWG:        object.wg,
				wg:              &sync.WaitGroup{},
			}
			ctx.ctx, ctx.cancel = context.WithCancel(context.Background())
			object.contextPool[i] = ctx
			go ctx.loop()
		}
		go object.recycleRouting()
		routinePool = object
	})
	return routinePool
}

// /回收上下文
func (object *RoutinePool) recycleContext(ctx *Context) {
	ctx.cancel()
	ctx.wg.Wait()

	ctx.Lock()
	close(ctx.taskCh)
	ctx.taskCh = nil
	ctx.Unlock()
}

// /检查空闲协程数量
func (object *RoutinePool) recycleRouting() {
	object.wg.Add(1)
loop:
	for {
		select {
		case <-object.ctx.Done():
			break loop
		case <-object.recycleNotifyCh:
			object.Lock()
			if object.minRoutine < int64(len(object.contextPool)) {
				for i := len(object.contextPool) - 1; i >= 0; i-- {
					ctx := object.contextPool[i]
					if 1 == atomic.LoadInt32(&ctx.idleFlag) {
						object.recycleContext(ctx)
						object.contextPool = append(object.contextPool[:i], object.contextPool[i+1:]...)
					}
				}
			}
			object.Unlock()
		}
	}
	object.wg.Done()
}

// PostTask /提交任务
func (object *RoutinePool) PostTask(task TaskMethod, params ...interface{}) {
	object.Lock()
	i := 0
	for i < len(object.contextPool) {
		ctx := object.contextPool[i]
		if 1 == atomic.LoadInt32(&(ctx.idleFlag)) {
			break
		}
		i++
	}
	if i < len(object.contextPool) {
		ctx := object.contextPool[i]
		atomic.StoreInt32(&(ctx.idleFlag), 0)

		ctx.RLock()
		if nil != ctx.taskCh && 0 < cap(ctx.taskCh) {
			ctx.taskCh <- &TaskParam{TaskMethod: task, TaskParam: params}
		}
		ctx.RUnlock()
	} else {
		ctx := &Context{
			idleFlag:        0,
			taskCh:          make(chan *TaskParam, 16),
			recycleNotifyCh: object.recycleNotifyCh,
			parentWG:        object.wg,
			wg:              &sync.WaitGroup{},
		}
		ctx.ctx, ctx.cancel = context.WithCancel(context.Background())
		object.contextPool = append(object.contextPool, ctx)
		go ctx.loop()
		ctx.taskCh <- &TaskParam{TaskMethod: task, TaskParam: params}
	}
	object.Unlock()
}

// Name /名字
func (object *RoutinePool) Name() string {
	return "RoutinePool"
}

// ShutdownPriority /设置关闭优先级
func (object *RoutinePool) SetShutdownPriority(priority int) {
	object.priority = priority
}

// ShutdownPriority /关闭优先级
func (object *RoutinePool) ShutdownPriority() int {
	return object.priority
}

// BeforeShutdown /应用退出前
func (object *RoutinePool) BeforeShutdown() {
	object.cancel()
	object.wg.Wait()
	object.Lock()
	for _, ctx := range object.contextPool {
		object.recycleContext(ctx)
	}
	object.contextPool = make([]*Context, 0)
	close(object.recycleNotifyCh)
	object.recycleNotifyCh = nil
	object.Unlock()
}
func (object *RoutinePool) AfterShutdown() {

}
