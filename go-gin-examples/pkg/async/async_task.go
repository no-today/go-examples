package async

import (
	"cathub.me/go-gin-examples/pkg/setting"
	"github.com/gammazero/workerpool"
	"sync"
)

var workerPool *workerpool.WorkerPool
var _onceWorkerPool sync.Once

func GetWorkerPool() *workerpool.WorkerPool {
	_onceWorkerPool.Do(func() {
		workerPool = workerpool.New(setting.Async.WorkerPoolSize)
	})
	return workerPool
}

func Submit(task func()) {
	GetWorkerPool().Submit(task)
}
