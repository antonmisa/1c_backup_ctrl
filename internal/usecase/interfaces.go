// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/antonmisa/1cctl_cli/internal/entity"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --all

type (
	// Ctrl -.
	Ctrl interface {
		ClusterByName(ctx context.Context, clusterName string) (entity.Cluster, error)
		InfobaseByName(ctx context.Context, cluster entity.Cluster, infobaseName string, clusterCred entity.Credentials) (entity.Infobase, error)

		Sessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Session, error)
		DisableSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error
		EnableSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error
		DeleteSessions(ctx context.Context, cluster entity.Cluster, sessions []entity.Session, clusterCred entity.Credentials) error

		Connections(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Connection, error)
		DeleteConnections(ctx context.Context, cluster entity.Cluster, connections []entity.Connection, clusterCred entity.Credentials) error

		RunBackup(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, infobaseCred entity.Credentials, lockCode string, outputPath string) (string, error)
	}

	// CtrlPipe -.
	CtrlPipe interface {
		GetClusters(ctx context.Context) ([]entity.Cluster, error)
		GetInfobases(ctx context.Context, cluster entity.Cluster, clusterCred entity.Credentials) ([]entity.Infobase, error)
		GetSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Session, error)
		GetConnections(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Connection, error)

		DisableSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error
		EnableSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error

		DeleteSession(ctx context.Context, cluster entity.Cluster, session entity.Session, clusterCred entity.Credentials) error
		DeleteSessions(ctx context.Context, cluster entity.Cluster, sessions []entity.Session, clusterCred entity.Credentials) error

		DeleteConnection(ctx context.Context, cluster entity.Cluster, connection entity.Connection, clusterCred entity.Credentials) error
		DeleteConnections(ctx context.Context, cluster entity.Cluster, connections []entity.Connection, clusterCred entity.Credentials) error
	}

	// CtrlBackup -.
	CtrlBackup interface {
		RunBackup(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, infobaseCred entity.Credentials, lockCode string, outputPath string) error
	}
)
