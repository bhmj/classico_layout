package app

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bhmj/classico_layout/internal/pkg/config"
	"github.com/bhmj/classico_layout/internal/pkg/html"
	"github.com/bhmj/classico_layout/internal/pkg/log"
	"github.com/bhmj/classico_layout/internal/pkg/service"
)

var errSignalReceived = errors.New("signal received")

// App represents runnable application
type App interface {
	Run()
}

type app struct {
	cfg *config.Config
	log log.Logger
}

// New creates application
func New(cfg *config.Config, log log.Logger) App {
	return &app{cfg: cfg, log: log}
}

// Run runs application
func (s *app) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())

	s.launch(ctx, cancel)
}

func (s *app) launch(ctx context.Context, cancel context.CancelFunc) {
	var wg sync.WaitGroup

	srv := service.NewService(s.cfg)

	srv.GenerateMatrix()

	html.WriteHTML(srv.GetMatrix(), "matrix.html", "Base matrix")

	wg.Add(1)
	// Main calculation
	go func() {
		defer wg.Done()
		srv.Run(ctx)
	}()

	// SIGTERM handler run
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <-ctx.Done():
			return
		case signal := <-ch:
			cancel()
			s.log.L().Info(fmt.Errorf("%w: %s", errSignalReceived, signal))
		}
	}()

	wg.Wait()
}
