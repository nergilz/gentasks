package main

import (
	"context"
	"gentasks/generator"
	"gentasks/handler"
	"log"
	"time"
)

func main() {
	log.Println("run app worker")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*9)
	defer cancel()

	app := handler.InitAppTaskHandler(ctx)

	app.Wg.Add(1)
	go generator.Send(ctx, app.SendTaskCh, app.Wg)

	var d chan struct{}

	app.Wg.Add(1)
	go app.Recv(d)

	app.Wg.Add(1)
	go app.OutputSuccessData()

	app.Wg.Add(1)
	go app.OutputFailedData()

	app.Wg.Wait()

	log.Println("worker stoped")

	time.Sleep(time.Second * 5)
	log.Println("timer 5 sec")
}

// //----------------оригинальный код из задачи---------------------------------------------------------------

/*
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func _main() {
	taskCreturer := func(a chan Ttype) {

		go func() {
			for {
				time.Sleep(time.Millisecond * 100)

				tnow := time.Now()
				ft := tnow.Format(time.RFC3339)

				if tnow.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}

				res := Ttype{
					cT: ft,
					id: int(tnow.Unix()),
				}

				a <- res // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)

		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = task_worker(t)
			go tasksorter(t)
		}
		close(superChan)
	}()

	result := map[int]Ttype{}
	err := []error{}

	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range undoneTasks {
			go func() {
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
*/
