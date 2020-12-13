package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

var shutdownChan = make(chan error)

func main() {
	ctx := context.Background()

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ServeApp(gctx.Done())
	})

	g.Go(func() error {
		return NotifySignal(gctx.Done())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("errgroup err %s\n", err)
	}
	if err := <-shutdownChan; err != nil {
		fmt.Printf("shutdown err %s\n", err)
	}
}

// HelloHandler 返回hello
type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello")
}

// ServeApp 服务启动
func ServeApp(stop <-chan struct{}) error {
	addr := "0.0.0.0:8080"
	handler := &HelloHandler{}
	return serve(addr, handler, stop)
}

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	svr := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop
		fmt.Println("serve rcv stop and svr is about to shutdown")
		shutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		shutdownChan <- svr.Shutdown(shutCtx)
		fmt.Println("svr truely shutdown")
	}()
	return svr.ListenAndServe()
}

// NotifySignal 监听信号
func NotifySignal(stop <-chan struct{}) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	defer func() {
		signal.Stop(sigint)
		close(sigint)
	}()

	select {
	case <-stop:
		fmt.Println("NotifySignal rcv stop")
	case <-sigint:
		fmt.Println("sigint notified")
		return errors.New("sigint notified")
	}
	return nil
}
