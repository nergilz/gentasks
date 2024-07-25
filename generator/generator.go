package generator

import (
	"context"
	"gentasks/handler"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TaskDuration = 100
)

type Generator struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	SendCh chan handler.Task
}

func InitGenerator(ctx context.Context, sendCh chan handler.Task, wg *sync.WaitGroup) *Generator {
	return &Generator{
		ctx:    ctx,
		wg:     wg,
		SendCh: sendCh,
	}
}

// генератор создает канал и когда вызываеш reciv передаешь его
func (g *Generator) Send(sendC *uint32) {
	defer g.wg.Done()
	defer close(g.SendCh)

	timer := time.NewTimer(time.Second * time.Duration(handler.WorkerTimer))
	defer timer.Stop()

	for {
		select {
		case <-g.ctx.Done():
			log.Println(" ctx send done")
			return
		case <-timer.C:
			log.Println(" timer stop")
			return
		case <-time.After(time.Millisecond * TaskDuration):
			t := time.Now()
			newTask := handler.Task{
				Id:         int(t.UnixMicro()),
				CreateTime: t.Format(time.RFC3339),
				State:      1,
			}
			// условие появления ошибочных тасков
			if t.Nanosecond()%2 > 0 {
				newTask.State = 0
			}

			atomic.AddUint32(sendC, 1)
			g.SendCh <- newTask
		}
	}
}
