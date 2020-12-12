package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	g, gctx := errgroup.WithContext(ctx)

	svr := &http.Server{
		Addr: ":8080",
	}
	g.Go(func() error {
		return StartServer(svr)
	})

	g.Go(func() error {
		return NotifySignal(gctx, svr)
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

// StartServer 服务启动
func StartServer(svr *http.Server) error {
	if err := svr.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	fmt.Println("svr ErrServerClosed")
	return nil
}

// NotifySignal 监听信号
func NotifySignal(gctx context.Context, svr *http.Server) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-gctx.Done():
		fmt.Println("gctx cancel")
	case <-sigint:
		fmt.Println("sigint shutdown")
		shutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		svr.Shutdown(shutCtx)
	}
	return nil
}
