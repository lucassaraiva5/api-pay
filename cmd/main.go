package main

import (
	"lucassaraiva5/api-pay/internal/app"
	"lucassaraiva5/api-pay/internal/infra/logger"
	"lucassaraiva5/api-pay/internal/infra/variables"
)

func main() {
	logger.Init(&logger.Option{
		ServiceName:    variables.ServiceName(),
		ServiceVersion: variables.ServiceVersion(),
		Environment:    variables.Environment(),
		LogLevel:       variables.LogLevel(),
	})

	defer logger.Sync()

	application := app.Instance()
	application.Start(false)
}
