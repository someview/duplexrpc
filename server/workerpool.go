package server

import "time"

type WorkerPool interface {
	Submit(task func())
	StopAndWaitFor(deadline time.Duration)
	Stop()
	StopAndWait()
}
