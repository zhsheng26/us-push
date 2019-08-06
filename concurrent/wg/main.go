package wg

import (
	"fmt"
	"sync"
	"time"
)

//
// @author zhangsheng
// @date 2019/7/31
//

var limit = make(chan int, 3)

type job func(it string, wg *sync.WaitGroup)

func main() {
	var work = []job{use, use, use, use, use, use}
	var wg sync.WaitGroup
	//只能同时启动三个协程
	for _, w := range work {
		wg.Add(1)
		go func() {
			limit <- 1
			w("hhh", &wg)
			<-limit
		}()
	}
	wg.Wait()
}

func use(it string, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	fmt.Println(it)
	wg.Done()
}
