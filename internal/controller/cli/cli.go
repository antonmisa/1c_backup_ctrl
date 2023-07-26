package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	e "github.com/antonmisa/1cctl_cli/internal/common/clierror"
	"github.com/antonmisa/1cctl_cli/internal/entity"
	"github.com/antonmisa/1cctl_cli/internal/usecase"
)

const (
	_defaultBackupTimeout    = 60
	_defaultOperationTimeout = 5000
)

type Ctrl1CCLI struct {
	ctx context.Context
	c   usecase.Ctrl
}

func New(ctx context.Context, c usecase.Ctrl) *Ctrl1CCLI {
	return &Ctrl1CCLI{
		ctx: ctx,
		c:   c,
	}
}

func (cc *Ctrl1CCLI) Backup(clusterName string, infobase string,
	clusterAdmin string, clusterPwd string,
	infobaseAdmin string, infobasePwd string,
	lockCode string, outputPath string) (re error) {

	ctx, cancel := context.WithTimeout(cc.ctx, _defaultOperationTimeout*time.Second)
	defer cancel()

	clusterCred := entity.Credentials{
		Name: clusterAdmin,
		Pwd:  clusterPwd,
	}

	infobaseCred := entity.Credentials{
		Name: infobaseAdmin,
		Pwd:  infobasePwd,
	}

	// Check cluster exists
	cl, err := cc.c.ClusterByName(ctx, clusterName)

	if err != nil {
		re = fmt.Errorf("cli - Process - cc.c.ClusterByName: %w", err)
		return
	}

	// Check infobase exists in cluster
	ib, err := cc.c.InfobaseByName(ctx, cl, infobase, clusterCred)

	if err != nil {
		re = fmt.Errorf("cli - Process - cc.c.InfobaseByName: %w", err)
		return
	}

	defer func() {
		// UnBlock all sessions in infobase, always
		c, cncl := context.WithTimeout(context.TODO(), _defaultOperationTimeout*time.Second)
		defer cncl()

		err = cc.c.EnableSessions(c, cl, ib, clusterCred, infobaseCred, lockCode)
		if err != nil {
			re = fmt.Errorf("cli - Process - cc.c.EnableSessions: %w", err)
		}
	}()

	// Block all new sessions in infobase
	err = cc.c.DisableSessions(ctx, cl, ib, clusterCred, infobaseCred, lockCode)

	if err != nil {
		re = fmt.Errorf("cli - Process - cc.c.DisableSessions: %w", err)
		return
	}

	// Get all sessions in infobase
	sessions, err := cc.c.Sessions(ctx, cl, ib, clusterCred)

	if err != nil {
		re = fmt.Errorf("cli - Process - cc.c.Sessions: %w", err)
		return
	}

	// Drop all sessions in infobase
	_ = cc.c.DeleteSessions(ctx, cl, sessions, clusterCred) //nolint:errcheck // do not need errors

	// Get all connections in infobase
	connections, err := cc.c.Connections(ctx, cl, ib, clusterCred)

	if err != nil {
		re = fmt.Errorf("cli - Process - cc.c.Connections: %w", err)
		return
	}

	// Drop all connections in infobase
	_ = cc.c.DeleteConnections(ctx, cl, connections, clusterCred) //nolint:errcheck // do not need errors

	// Run backup as long operation
	cx, cancel := context.WithTimeout(cc.ctx, _defaultBackupTimeout*time.Minute)

	defer cancel()

	p, err := cc.c.RunBackup(cx, cl, ib, infobaseCred, lockCode, outputPath)

	if err != nil {
		re = fmt.Errorf("cli - Process - cc.c.RunBackup: %w", err)
		return
	}

	// Check backup exists
	_, err = os.Stat(p)

	if os.IsNotExist(err) {
		re = e.WithText{
			Txt: fmt.Sprintf("backup file does not exist at: %s", p),
		}
		return
	}

	return
}
