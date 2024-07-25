package handler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WorkerTimer    = 10
	OutpuTicker    = 3
	SuccessTaskMsg = "SUCCESS, task has been work"
	FailTaskMsg    = "FAILED, task not has been work"
)

type (
	AppHandler struct {
		Ctx           context.Context
		cancel        context.CancelFunc
		Wg            *sync.WaitGroup
		SendTaskCh    chan Task
		WarningTaskCh chan Task
		SuccessTaskCh chan Task
		store         *AppStore
		ticker        *time.Ticker
		done          chan bool
	}

	Task struct {
		Id            int
		CreateTime    string
		ExecutionTime string
		WorkResult    string
		State         int
	}
)

func InitAppHandler(ctx context.Context, store *AppStore, cancel context.CancelFunc) *AppHandler {
	return &AppHandler{
		Ctx:           ctx,
		Wg:            &sync.WaitGroup{},
		SendTaskCh:    make(chan Task, 256),
		WarningTaskCh: make(chan Task, 256),
		SuccessTaskCh: make(chan Task, 256),
		store:         store,
		cancel:        cancel,
		ticker:        time.NewTicker(time.Second * OutpuTicker),
		done:          make(chan bool, 1),
	}
}

// сonsume channel with tasks
func (app *AppHandler) Recv() {
	defer app.Wg.Done()
	defer close(app.SuccessTaskCh)
	defer close(app.WarningTaskCh)
	defer close(app.done)

	for {
		select {
		case <-app.Ctx.Done():
			return
		case task, ok := <-app.SendTaskCh:
			if ok {
				switch task.State {
				case 0:
					task.WorkResult = FailTaskMsg
					task.ExecutionTime = time.Now().Format(time.RFC3339)
					app.WarningTaskCh <- task
				case 1:
					task.WorkResult = SuccessTaskMsg
					task.ExecutionTime = time.Now().Format(time.RFC3339)
					app.SuccessTaskCh <- task
				}
			} else {
				app.done <- true
				return
			}
		}
	}
}

func (app *AppHandler) LoadSuccess() {
	defer app.Wg.Done()

	for task := range app.SuccessTaskCh {
		app.store.LoadSuccess(task)
	}
}

func (app *AppHandler) LoadFailed() {
	defer app.Wg.Done()

	for task := range app.WarningTaskCh {
		app.store.LoadFail(task)
	}
}

func (app *AppHandler) Output(sendC, recvC *uint32) {
	defer app.Wg.Done()

	for {
		select {
		case <-app.Ctx.Done():
			return
		case <-app.done:
			fmt.Println(" send worker done")
			return
		case _, ok := <-app.ticker.C:
			if ok {
				app.printSeparateTaskFromStore(recvC)
				log.Println("send task:", *sendC, "recive task", *recvC)
			} else {
				return
			}
		}
	}
}

func (app *AppHandler) printSeparateTaskFromStore(recvC *uint32) {
	app.store.mutex.RLock()

	for id, task := range app.store.FailedData {
		fmt.Println("ID:", id, task)
		atomic.AddUint32(recvC, 1)
	}

	for id, task := range app.store.SuccessData {
		fmt.Println("ID:", id, task)
		atomic.AddUint32(recvC, 1)
	}

	app.store.mutex.RUnlock()
	app.store.ClearMap()
}
