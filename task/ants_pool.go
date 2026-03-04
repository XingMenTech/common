package task

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var generic *ants.PoolWithFuncGeneric[AntsPoolParam]

var pool *ants.Pool

type AntsPoolParam struct {
	FuncCode string
	Arg      interface{}
}

var funcMap = make(map[string]func(interface{}))
var lock sync.Mutex

func RegisterAntsPoolFunc(code string, fn func(interface{})) {
	lock.Lock()
	defer lock.Unlock()
	funcMap[code] = fn
}

func InitAntsFuncPool(max int) (err error) {
	generic, err = ants.NewPoolWithFuncGeneric(max, func(i AntsPoolParam) {
		code := i.FuncCode
		if fu, ok := funcMap[code]; ok {
			fu(i.Arg)
		}

	}, ants.WithPreAlloc(true), ants.WithNonblocking(false))
	if err != nil {
		return
	}

	return
}

func Invoke(code string, arg interface{}) error {

	return generic.Invoke(AntsPoolParam{
		FuncCode: code,
		Arg:      arg,
	})
}

func InitAntsWorkPool(max int) {
	var err error
	pool, err = ants.NewPool(max, ants.WithPreAlloc(true), ants.WithNonblocking(false))
	if err != nil {
		panic(fmt.Sprintf("goroutine池启动失败,error: %+v", err))
	}
	return
}

func Submit(fn func()) error {
	return pool.Submit(fn)
}
