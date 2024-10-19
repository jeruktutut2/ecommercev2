package main

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/setups"
	"backend-golang/commons/utils"
	"context"
	"os"
	"os/signal"
)

func main() {
	postgresUtil := utils.NewPostgresConnection()
	defer postgresUtil.Close()

	redisUtil := utils.NewRedisConnection()
	defer redisUtil.Close()

	validate := setups.SetValidator()
	// bcryptHelper := helpers.NewBcryptHelper()
	uuidHelper := helpers.NewUuidHelper()
	redisHelper := helpers.NewRedisHelper()

	e := setups.SetEcho(postgresUtil, redisUtil, validate, uuidHelper, redisHelper)
	setups.StartEcho(e)
	defer setups.StopEcho(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
