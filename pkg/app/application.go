package app

import (
	"time"

	"os"

	"fmt"

	"github.com/kudryashov-sv/pgspy/pkg/config"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Application struct {
	Log    *zap.Logger
	Config config.Configurer
}

func Builder(l *zap.Logger, c config.Configurer) *Application {
	return &Application{
		Log:    l.With(zap.String("app", "pgspy")),
		Config: c,
	}
}

func (a *Application) Run() error {
	a.Log.Info("application started", zap.Reflect("config", a.Config))

	u := a.Config.PostgresURL()
	cs, err := pq.ParseURL(u)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}
	a.Log.Debug("create listener", zap.String("connection", cs))
	l := pq.NewListener(cs, time.Second, time.Second*10, func(e pq.ListenerEventType, err error) {
		a.Log.Debug("connect event", zap.Reflect("event", e), zap.Error(err))
	})

	name := a.Config.NotifyName()
	defer l.Unlisten(name)

	err = l.Listen(name)
	if err != nil {
		return errors.Wrap(err, "listen")
	}

	return a.eventLoop(l.NotificationChannel())
}

func (a *Application) eventLoop(ch <-chan *pq.Notification) error {
	fd, err := os.OpenFile(a.Config.OutPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open output file")
	}
	defer fd.Close()

	a.Log.Debug("run event loop")
	for n := range ch {
		if n == nil {
			continue
		}
		_, err := fmt.Fprintln(fd, n.Extra)
		if err != nil {
			return errors.Wrap(err, "write output")
		}
	}
	a.Log.Debug("exit event loop")
	return nil
}
