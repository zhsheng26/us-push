package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//
// @author zhangsheng
// @date 2019/7/31
//
func worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		default:
			time.Sleep(time.Second)
			fmt.Println("hello")
		case <-ctx.Done():
			return
		}
	}
}
func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//2s 后自动取消
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}
	wg.Wait()
}
