package grpcx

import (
	"github.com/soheilhy/cmux"
	"net"
)

type Server interface {
	InstallShutdownHook(hook ShutdownHook)
	start() error
	stop() error
	Serve(l net.Listener) error
	getAddr() string
	typeCode() string
	Match() cmux.Matcher
}

// ShutdownHook /关闭钩子
type ShutdownHook interface {
	// Name /名字
	Name() string
	// ShutdownPriority /优先级
	ShutdownPriority() int
	// BeforeShutdown /应用程序退出前
	BeforeShutdown()
	// AfterShutdown /应用程序退出后
	AfterShutdown()
}

// ShutdownHooks /排序
type shutdownHooks []ShutdownHook

func (object shutdownHooks) Len() int {
	return len(object)
}
func (object shutdownHooks) Less(i, j int) bool {
	return object[i].ShutdownPriority() > object[j].ShutdownPriority()
}
func (object shutdownHooks) Swap(i, j int) {
	object[i], object[j] = object[j], object[i]
}
