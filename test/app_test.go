package test

import (
	"context"
	"gentasks/generator"
	"gentasks/handler"
	"os"
	"os/signal"
	"testing"
)

func TestMain(t *testing.T) {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	store := handler.InitStore()
	app := handler.InitAppHandler(ctx, store)
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

	if sendC != recvC {
		t.Errorf("send task count: %d and recive task count: %d must by equal", sendC, recvC)
	}
}
