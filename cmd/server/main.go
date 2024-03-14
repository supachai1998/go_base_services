package main

import (
	"context"

	"go_base/logger"
	"go_base/server"
)

func main() {
	err := server.Run(context.Background())
	if err != nil {
		logger.L().Fatal(err)
	}
}
