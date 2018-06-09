package main

import (
	"os"

	"github.com/kudryashov-sv/pgspy/pkg/app"
	"github.com/kudryashov-sv/pgspy/pkg/config"
	"github.com/kudryashov-sv/pgspy/pkg/log"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func main() {
	di := dig.New()
	di.Provide(log.Builder)
	di.Provide(config.Builder)
	di.Provide(app.Builder)

	err := di.Invoke(func(a *app.Application) error {
		err := a.Run()
		if err != nil {
			a.Log.Error("error exit", zap.Error(err))
		}
		return err
	})
	if err != nil {
		os.Exit(1)
	}
}
