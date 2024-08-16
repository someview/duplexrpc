package client

import (
	"context"

	netpoll "github.com/cloudwego/netpoll"

	"github.com/panjf2000/ants/v2"
)

func init() {
	SetAntsRunner(-1)
	//SetRunner(func(task func()) error {
	//	go task()
	//	return nil
	//})
}

// Runner is a function that runs the task.
var runTask Runner

type Runner func(task func()) error

// SetRunner sets the runner for the server to run handle-task.
// By default, the runner is a goroutine.
// runner can resolve the problem of goroutine leak and stack expansion.
func SetRunner(r Runner) {
	runTask = r
	netpoll.SetRunner(func(ctx context.Context, f func()) {
		_ = runTask(f)
	})
}

// SetAntsRunner uses the ants callCtxPool to SetRunner.
// it equals to SetRunner(antsPool.Submit).
func SetAntsRunner(size int) {
	pool, _ := ants.NewPool(size)
	SetRunner(pool.Submit)
}
