package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonmisa/1cctl_cli/config"
	"github.com/antonmisa/1cctl_cli/internal/controller/cli"
	"github.com/antonmisa/1cctl_cli/internal/usecase"
	ucbackup "github.com/antonmisa/1cctl_cli/internal/usecase/backup"
	ucpipe "github.com/antonmisa/1cctl_cli/internal/usecase/pipe"
	"github.com/antonmisa/1cctl_cli/pkg/logger"
	"github.com/antonmisa/1cctl_cli/pkg/pipe"
)

var (
	ErrEmptyClusterOrInfobase = errors.New("app - RunCLI - empty cluster or infobase")
	ErrEmptyClusterConnection = errors.New("app - RunCLI - empty cluster connection string")
)

func Run(cfg *config.Config, clusterConnection, clusterName, clusterAdmin, clusterPwd, infobase, infobaseUser, infobasePwd, outputPath string) {

	l, err := logger.New(cfg.Log.Path, cfg.Log.Level)
	if err != nil {
		l.Fatal(fmt.Errorf("app - RunCLI - logger.New: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-interrupt
		l.Info("app - RunCLI - signal: " + s.String())
		cancel()
	}()

	if clusterConnection == "" {
		l.Fatal(ErrEmptyClusterConnection) //nolint:goerr13 // high level error
	}

	if clusterName == "" && infobase == "" {
		l.Fatal(ErrEmptyClusterOrInfobase) //nolint:goerr13 // high level error
	}

	p, err := pipe.New(cfg.App.PathToRAC)
	if err != nil {
		l.Fatal(fmt.Errorf("app - RunCLI - pipe.New: %w", err))
	}

	ctrlPipe := ucpipe.New(p, clusterConnection)

	ctrlBackup, err := ucbackup.New(cfg.App.PathTo1C)

	if err != nil {
		l.Fatal(fmt.Errorf("app - RunCLI - ucbackup.New: %w", err))
	}

	ucCtrl := usecase.New(ctrlPipe, ctrlBackup)

	ctrl := cli.New(ctx, ucCtrl)

	now := time.Now()

	l.Info("app - RunCLI - start")

	err = ctrl.Backup(clusterName, infobase,
		clusterAdmin, clusterPwd,
		infobaseUser, infobasePwd,
		cfg.App.LockCode, outputPath)

	if err != nil {
		l.Fatal(err)
	}

	l.Info("app - RunCLI - succefully end, time taken: %s", time.Since(now).String())
}
