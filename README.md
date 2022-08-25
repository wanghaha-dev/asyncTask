# Async Task

![](asyncTask.jpg)

安装
```go
go get github.com/wanghaha-dev/asyncTask
```

使用
```go
package main

import (
	"context"
	"fmt"
	"github.com/wanghaha-dev/asyncTask/asyncTask"
	"sync"
)

func main() {
	ctx := context.Background()
	task, err := asyncTask.NewTask(ctx, asyncTask.Config{
		Addr:      "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// take task
	go func() {
		asyncTask.Each(1000, func() {
			data, err := task.TakeNormalTask()
			if err != nil {
				panic(err)
			}
			fmt.Println("data:", data)
		})
		wg.Done()
	}()

	// put task
	asyncTask.Each(1000, func() {
		err = task.PutNormalTask("ddd", asyncTask.Map{
			"name": "wanghaha",
		})
		if err != nil {
			panic(err)
		}
	})

	wg.Wait()
	fmt.Println("finish.")
}
```