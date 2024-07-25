package main

import (
	"context"
	"gentasks/generator"
	"gentasks/handler"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("run app worker")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	store := handler.InitStore()
	app := handler.InitAppHandler(ctx, store, cancel)
	gen := generator.InitGenerator(ctx, app.SendTaskCh, app.Wg)

	var sendC uint32
	var recvC uint32

	app.Wg.Add(1)
	go gen.Send(&sendC)

	app.Wg.Add(1)
	go app.Recv()

	app.Wg.Add(1)
	go app.LoadSuccess()

	app.Wg.Add(1)
	go app.LoadFailed()

	app.Wg.Add(1)
	go app.Output(&sendC, &recvC)

	app.Wg.Wait()
	log.Println("worker stoped")
}
