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
	Ctx    context.Context
	Wg     *sync.WaitGroup
	SendCh chan handler.Task
}

func InitGenerator(ctx context.Context, sendCh chan handler.Task, wg *sync.WaitGroup) *Generator {
	return &Generator{
		Ctx:    ctx,
		Wg:     wg,
		SendCh: sendCh,
	}
}

func (g *Generator) Send(sendC *uint32) {
	defer g.Wg.Done()
	defer close(g.SendCh)

	timer := time.NewTimer(time.Second * time.Duration(handler.WorkerTimer))

	for {
		select {
		case <-g.Ctx.Done():
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
