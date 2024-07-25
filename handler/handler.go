package handler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	WorkerTimer    int = 10
	SuccessTaskMsg     = "SUCCESS task, has been work"
	ErrorTaskMsg       = "FAILED task, not has been work"
)

type (
	AppTaskHandler struct {
		Ctx           context.Context
		Wg            *sync.WaitGroup
		SendTaskCh    chan Task
		WarningTaskCh chan Task
		SuccessTaskCh chan Task
		store         *AppStore
		ticker        *time.Ticker
	}

	Task struct {
		Id            int
		CreateTime    string // время создания
		ExecutionTime string // время выполнения
		ResultRunTask []byte
	}
)

func InitAppTaskHandler(ctx context.Context, store *AppStore) *AppTaskHandler {
	ticker := time.NewTicker(time.Second * 2)
	return &AppTaskHandler{
		Ctx:           ctx,
		Wg:            &sync.WaitGroup{},
		SendTaskCh:    make(chan Task),
		WarningTaskCh: make(chan Task),
		SuccessTaskCh: make(chan Task),
		store:         store,
		ticker:        ticker,
	}
}

// сonsume channel with tasks
func (app *AppTaskHandler) Recv(stopCh chan bool) {
	defer app.Wg.Done()

	for {
		select {
		case <-stopCh:
			close(stopCh)
			return
		case task, ok := <-app.SendTaskCh:
			if ok {
				// разделение тасков
				if CheckCreateTime(task.CreateTime) {
					task.ResultRunTask = []byte(SuccessTaskMsg)
					app.SuccessTaskCh <- task
				} else {
					task.ResultRunTask = []byte(ErrorTaskMsg)
					app.WarningTaskCh <- task
				}
			} else {
				close(app.SuccessTaskCh)
				close(app.WarningTaskCh)
				// stopCh <- true
				// close(stopCh)
				// app.ticker.Stop()
				return
			}
		}
	}
}

// обрабатываем отдельно Success и Failed
func (app *AppTaskHandler) OutputSuccessData() {
	defer app.Wg.Done()

	for task := range app.SuccessTaskCh {
		fmt.Println("SUCCESS:", task)
		// app.store.LoadSuccess(task)
	}

}

// обрабатываем отдельно Failed Tasks и Success Tasks
func (app *AppTaskHandler) OutputFailedData() {
	defer app.Wg.Done()

	for task := range app.WarningTaskCh {
		fmt.Println("FAILED:", task)
		// app.store.LoadFail(task)
	}
}

func (app *AppTaskHandler) PrintData(d chan bool) {
	defer app.Wg.Done()

	for {
		select {
		case <-d:
			return
		case _, ok := <-app.ticker.C:
			if ok {
				app.store.mutex.RLock()
				for id, v := range app.store.FailedData {
					fmt.Println("FAILED:", "ID:", id, v)
				}
				for id, v := range app.store.SuccessData {
					fmt.Println("SUCCESS:", "ID:", id, v)
				}
				app.store.mutex.RUnlock()
				app.store.CleanMap()
				fmt.Println("---")
			} else {
				return
			}
		}
	}
}

// проверка состояния выполнения задачи
func CheckCreateTime(createTime string) bool {
	var res bool
	_, err := time.Parse(time.RFC3339, createTime)
	if err != nil {
		res = false
		return res
	}
	res = true
	return res
}
