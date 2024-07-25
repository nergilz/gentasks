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
	}

	Task struct {
		Id            int
		CreateTime    string // время создания
		ExecutionTime string // время выполнения
		ResultRunTask []byte
	}
)

func InitAppTaskHandler(ctx context.Context) *AppTaskHandler {
	return &AppTaskHandler{
		Ctx:           ctx,
		Wg:            &sync.WaitGroup{},
		SendTaskCh:    make(chan Task),
		WarningTaskCh: make(chan Task),
		SuccessTaskCh: make(chan Task),
	}
}

// сonsume channel with tasks
func (app *AppTaskHandler) Recv(d chan struct{}) {
	defer app.Wg.Done()

	for {
		select {
		case <-d:
			return
		// case <-app.Ctx.Done():
		// 	close(app.SuccessTaskCh)
		// 	close(app.WarningTaskCh)
		// 	fmt.Println(" ctx done")
		// 	return
		case task, ok := <-app.SendTaskCh:
			if ok {
				if CheckCreateTime(task.CreateTime) {
					task.ResultRunTask = []byte(SuccessTaskMsg)
					// fmt.Println("SUCCESS:", task)
					app.SuccessTaskCh <- task
				} else {
					task.ResultRunTask = []byte(ErrorTaskMsg)
					// fmt.Println("FAILED:", task)
					app.WarningTaskCh <- task
				}
			} else {
				close(app.SuccessTaskCh)
				close(app.WarningTaskCh)
				return
			}
		}
	}
}

// обрабатываем отдельно Success Tasks и Failed Tasks
func (app *AppTaskHandler) OutputSuccessData() {
	defer app.Wg.Done()

	for task := range app.SuccessTaskCh { // !!!
		fmt.Println("SUCCESS:", task)
	}

	// for {
	// 	select {
	// 	case <-app.Ctx.Done():
	// 		return
	// 	case task, ok := <-app.SuccessTaskCh:
	// 		if ok {
	// 			fmt.Println("SUCCESS:", task)
	// 			// app.StoreTask.LoadSuccess(task)
	// 		} else {
	// 			// TODO done ch
	// 			return
	// 		}
	// 	}
	// }
}

// обрабатываем отдельно Failed Tasks и Success Tasks
func (app *AppTaskHandler) OutputFailedData() {
	defer app.Wg.Done()

	for task := range app.WarningTaskCh { // !!!
		fmt.Println("FAILED:", task)
	}

	// for {
	// 	select {
	// 	case <-app.Ctx.Done():
	// 		return
	// 	case task, ok := <-app.WarningTaskCh:
	// 		if ok {
	// 			fmt.Println("FAILED:", task)
	// 			// app.StoreTask.LoadFail(task)
	// 		} else {
	// 			// TODO done ch
	// 			return
	// 		}
	// 	}
	// }
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
