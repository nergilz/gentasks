package generator

import (
	"gentasks/handler"
	"log"
	"sync"
	"time"
)

type GeneratorTask struct {
}

// run tasks and write to channel with warning tasks
func Send(c chan handler.Task, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(time.Second * time.Duration(handler.WorkerTimer))

	for {
		select {
		// case <-ctx.Done():
		// 	close(c)
		// 	log.Println(" ctx done")
		// 	return
		case <-timer.C:
			close(c)
			log.Println(" timer click")
			return
		default:
			time.Sleep(time.Millisecond * 100)

			t := time.Now()
			create := t.Format(time.RFC3339)

			newTask := handler.Task{
				Id:         int(t.Unix()),
				CreateTime: create,
			}
			// условие появления ошибочных тасков
			if t.Nanosecond()%2 > 0 {
				newTask.CreateTime = "bad create time"
			}
			c <- newTask // !!!
		}
	}
}

func TestSend(c chan handler.Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)

		t := time.Now()
		create := t.Format(time.RFC3339)

		newTask := handler.Task{
			Id:         i,
			CreateTime: create,
		}

		if i%2 > 0 {
			newTask.CreateTime = "bad create time"
		}

		c <- newTask
	}

	close(c)
}
