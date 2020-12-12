package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	startTime := time.Now()
	ctx := context.Background()
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		t := time.NewTimer(1 * time.Second)
		select {
		case <-t.C:
		case <-gctx.Done():
		}
		return errors.New("sleep 1")
	})
	g.Go(func() error {
		t := time.NewTimer(10 * time.Second)
		select {
		case <-t.C:
		case <-gctx.Done():
		}
		return errors.New("sleep 10")
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("g err %s", err)
	}
	costTime := time.Since(startTime)
	fmt.Printf("time cost %s", costTime)
	return
}
