package grpcx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/valyala/fasthttp"
	"gitlab.novgate.com/common/common/logger"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sort"
	"sync"
)

type ServerImpl struct {
	addr               string
	lis                net.Listener
	shutdownHooks      []ShutdownHook
	shutdownHooksMutex sync.RWMutex
}

func (s *ServerImpl) InstallShutdownHook(hook ShutdownHook) {
	s.shutdownHooksMutex.Lock()
	if s.shutdownHooks == nil {
		s.shutdownHooks = make([]ShutdownHook, 0)
	}
	s.shutdownHooks = append(s.shutdownHooks, hook)
	s.shutdownHooksMutex.Unlock()
}

func (s *ServerImpl) beforeHook() {
	s.shutdownHooksMutex.Lock()
	sort.Sort(shutdownHooks(s.shutdownHooks))
	for _, v := range s.shutdownHooks {
		logger.LOG.Infof("app before shutdown: %s", v.Name())
		v.BeforeShutdown()
	}
	s.shutdownHooksMutex.Unlock()
}

func (s *ServerImpl) afterHook() {
	s.shutdownHooksMutex.Lock()
	sort.Sort(shutdownHooks(s.shutdownHooks))
	for _, v := range s.shutdownHooks {
		logger.LOG.Infof("app after shutdown: %s", v.Name())
		v.AfterShutdown()
	}
	s.shutdownHooksMutex.Unlock()
}

// ServiceRegister registry grpc services
type ServiceRegister func(*grpc.Server) []Service

// GRPCServer is a grpc server
type GRPCServer struct {
	servers []Server
	log     *logrus.Entry
}

func NewServer(servers []Server, log *logrus.Logger) *GRPCServer {
	return &GRPCServer{
		servers: servers,
		log:     log.WithField("model", "GatewayGrpcServer"),
	}
}

func (s *GRPCServer) Start() error {
	defer func() {
		if err := recover(); err != nil {
			s.log.Errorf("rpc: grpc server crash, errors:\n %+v", err)
		}
	}()

	for _, server := range s.servers {
		s.log.Infof("rpc: start a %s server at %s \n", server.typeCode(), server.getAddr())
		ser := server
		go func() {
			err := ser.start()
			if err != nil {
				s.log.Fatalf("rpc: start a %s server failed, errors:\n%+v", ser.typeCode(), err)
			}
		}()

	}

	return nil
}

func (s *GRPCServer) StartWithListener(l net.Listener) error {
	defer func() {
		if err := recover(); err != nil {
			s.log.Errorf("rpc: grpc server crash, errors:\n %+v", err)
		}
	}()

	m := cmux.New(l)

	for _, server := range s.servers {
		s.log.Infof("rpc: start a %s server at %s \n", server.typeCode(), server.getAddr())
		ser := server
		go func() {
			sl := m.Match(ser.Match())
			err := ser.Serve(sl)
			if err != nil {
				s.log.Fatalf("rpc: start a %s server failed, errors:\n%+v", ser.typeCode(), err)
			}
		}()

	}

	return m.Serve()
}

// Stop stop servers
func (s *GRPCServer) Stop() {
	for _, server := range s.servers {
		s.log.Infof("rpc: stop a %s server at %s \n", server.typeCode(), server.getAddr())
		err := server.stop()
		if err != nil {
			s.log.Fatalf("rpc: stop a %s server failed, errors:\n%+v", server.typeCode(), err)
		}
	}
}

func adjustAddr(addr string) string {
	if addr[0] == ':' {
		ips, err := intranetIP()
		if err != nil {
			logrus.Fatalf("get intranet ip failed, errors:\n%+v", err)
		}

		return fmt.Sprintf("%s%s", ips[0], addr)
	}

	return addr
}

type httpServer struct {
	ServerImpl
	server *http.Server
}

func WithHTTPServer(addr string, httpSetup func(engine *gin.Engine)) Server {
	handler := gin.Default()
	httpSetup(handler)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &httpServer{
		ServerImpl: ServerImpl{
			addr: addr,
		},
		server: srv,
	}
}

func (s *httpServer) start() error {
	return s.server.ListenAndServe()
}

func (s *httpServer) stop() error {
	s.beforeHook()
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	s.afterHook()
	return nil
}
func (s *httpServer) typeCode() string {
	return "http"
}
func (s *httpServer) getAddr() string {
	return s.addr
}

func (s *httpServer) Serve(l net.Listener) error {
	return s.server.Serve(l)
}
func (s *httpServer) Match() cmux.Matcher {
	return cmux.Any()
}

type fastHttpServer struct {
	ServerImpl
	server *fasthttp.Server
}

func WithFastHttpServer(addr string, handler fasthttp.RequestHandler) Server {
	httpS := &fasthttp.Server{
		Handler: handler,
	}

	return &fastHttpServer{
		ServerImpl: ServerImpl{
			addr: addr,
		},
		server: httpS,
	}
}

func (s *fastHttpServer) start() error {
	return s.server.ListenAndServe(s.addr)
}

func (s *fastHttpServer) stop() error {
	s.beforeHook()
	err := s.server.Shutdown()
	if err != nil {
		return err
	}
	s.afterHook()
	return nil
}
func (s *fastHttpServer) typeCode() string {
	return "fast http"
}
func (s *fastHttpServer) getAddr() string {
	return s.addr
}
func (s *fastHttpServer) Serve(l net.Listener) error {
	return s.server.Serve(l)
}
func (s *fastHttpServer) Match() cmux.Matcher {
	return cmux.HTTP1Fast()
}

type grpcServer struct {
	ServerImpl
	server *grpc.Server
}

func WithGrpcServer(addr string, register ServiceRegister, opt ...grpc.ServerOption) Server {

	server := grpc.NewServer(opt...)

	register(server)
	return &grpcServer{
		ServerImpl: ServerImpl{
			addr: addr,
		},
		server: server,
	}
}
func (s *grpcServer) start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.server.Serve(lis)
}

func (s *grpcServer) stop() error {
	s.beforeHook()
	s.server.GracefulStop()
	s.afterHook()
	return nil
}
func (s *grpcServer) Serve(l net.Listener) error {
	return s.server.Serve(l)
}

func (s *grpcServer) typeCode() string {
	return "grpc"
}
func (s *grpcServer) getAddr() string {
	return s.addr
}
func (s *grpcServer) Match() cmux.Matcher {
	return cmux.HTTP2HeaderField("content-type", "application/grpc")
}
