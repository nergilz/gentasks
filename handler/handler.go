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
		ctx           context.Context
		Wg            *sync.WaitGroup
		SendTaskCh    chan Task
		WarningTaskCh chan Task
		SuccessTaskCh chan Task
		Store         *AppStore
		done          chan int
	}

	Task struct {
		Id            int
		CreateTime    string
		ExecutionTime string
		WorkResult    string
		State         int
	}
)

func InitAppHandler(ctx context.Context, store *AppStore) *AppHandler {
	return &AppHandler{
		ctx:           ctx,
		Wg:            &sync.WaitGroup{},
		SendTaskCh:    make(chan Task, 256),
		WarningTaskCh: make(chan Task, 256),
		SuccessTaskCh: make(chan Task, 256),
		Store:         store,
		done:          make(chan int),
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
		case <-app.ctx.Done():
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
				app.done <- 1
				return
			}
		}
	}
}

func (app *AppHandler) LoadSuccess() {
	defer app.Wg.Done()

	for task := range app.SuccessTaskCh {
		app.Store.LoadSuccess(task)
	}
}

func (app *AppHandler) LoadFailed() {
	defer app.Wg.Done()

	for task := range app.WarningTaskCh {
		app.Store.LoadFail(task)
	}
}

func (app *AppHandler) Output(sendC, recvC *uint32) {
	defer app.Wg.Done()

	ticker := time.NewTicker(time.Second * OutpuTicker)
	defer ticker.Stop()

	for {
		select {
		case <-app.ctx.Done():
			return
		case <-app.done:
			app.printSeparateTaskFromStore(recvC)
			log.Println("send task:", *sendC, "recive task:", *recvC)
			fmt.Println(" send worker done")
			return
		case <-ticker.C:
			app.printSeparateTaskFromStore(recvC)
			log.Println("send task:", *sendC, "recive task:", *recvC)
		}
	}
}

func (app *AppHandler) printSeparateTaskFromStore(recvC *uint32) {
	app.Store.mutex.RLock()

	for id, task := range app.Store.FailedData {
		fmt.Println("ID:", id, task)
		atomic.AddUint32(recvC, 1)
	}

	for id, task := range app.Store.SuccessData {
		fmt.Println("ID:", id, task)
		atomic.AddUint32(recvC, 1)
	}

	app.Store.mutex.RUnlock()
	app.Store.ClearMap()
}
