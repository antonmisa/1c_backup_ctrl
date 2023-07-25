package usecase

import (
	"context"
	"fmt"
	"path"
	"time"

	e "github.com/antonmisa/1cctl_cli/internal/common/clierror"
	"github.com/antonmisa/1cctl_cli/internal/entity"
)

// CtrlUseCase -.
type CtrlUseCase struct {
	pipe   CtrlPipe
	backup CtrlBackup
}

var _ Ctrl = (*CtrlUseCase)(nil)

// New -.
func New(p CtrlPipe, b CtrlBackup) *CtrlUseCase {
	return &CtrlUseCase{
		pipe:   p,
		backup: b,
	}
}

// ClusterByName - getting cluster by name -.
func (uc *CtrlUseCase) ClusterByName(ctx context.Context, clusterName string) (entity.Cluster, error) {
	clusters, err := uc.pipe.GetClusters(ctx)
	if err != nil {
		return entity.Cluster{}, fmt.Errorf("CtrlUseCase - Clusters - uc.pipe.GetClusters: %w", err)
	}

	for i := range clusters {
		if fmt.Sprintf("%s:%s", clusters[i].Host, clusters[i].Port) == clusterName {
			return clusters[i], nil
		}
	}

	return entity.Cluster{}, e.WithText{
		Txt: fmt.Sprintf("cluster with name %s not found", clusterName)}
}

// Infobases - getting infobases list for cluster.
func (uc *CtrlUseCase) InfobaseByName(ctx context.Context, cluster entity.Cluster, infobaseName string, clusterCred entity.Credentials) (entity.Infobase, error) {
	infobases, err := uc.pipe.GetInfobases(ctx, cluster, clusterCred)
	if err != nil {
		return entity.Infobase{}, fmt.Errorf("CtrlUseCase - Infobases - uc.pipe.GetInfobases: %w", err)
	}

	for i := range infobases {
		if infobases[i].Name == infobaseName {
			return infobases[i], nil
		}
	}

	return entity.Infobase{}, e.WithText{
		Txt: fmt.Sprintf("infobase with name %s not found", infobaseName)}
}

// Sessions - getting sessions list for cluster.
func (uc *CtrlUseCase) Sessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Session, error) {
	sessions, err := uc.pipe.GetSessions(ctx, cluster, infobase, clusterCred)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Sessions - uc.pipe.GetSessions: %w", err)
	}

	return sessions, nil
}

// Disable new sessions for current infobase -.
func (uc *CtrlUseCase) DisableSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred, infobaseCred entity.Credentials, code string) error {
	err := uc.pipe.DisableSessions(ctx, cluster, infobase, clusterCred, infobaseCred, code)
	if err != nil {
		return fmt.Errorf("CtrlUseCase - DisableSessions - uc.pipe.DisableSessions: %w", err)
	}

	return nil
}

// Enable new sessions for current infobase -.
func (uc *CtrlUseCase) EnableSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred, infobaseCred entity.Credentials, code string) error {
	err := uc.pipe.EnableSessions(ctx, cluster, infobase, clusterCred, infobaseCred, code)
	if err != nil {
		return fmt.Errorf("CtrlUseCase - EnableSessions - uc.pipe.EnableSessions: %w", err)
	}

	return nil
}

// Delete sessions -.
func (uc *CtrlUseCase) DeleteSessions(ctx context.Context, cluster entity.Cluster, sessions []entity.Session, clusterCred entity.Credentials) error {
	err := uc.pipe.DeleteSessions(ctx, cluster, sessions, clusterCred)
	if err != nil {
		return fmt.Errorf("CtrlUseCase - DeleteSessions - uc.pipe.DeleteSessions: %w", err)
	}

	return nil
}

// Connections - getting connections list for cluster.
func (uc *CtrlUseCase) Connections(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Connection, error) {
	connections, err := uc.pipe.GetConnections(ctx, cluster, infobase, clusterCred)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Connections - uc.pipe.GetConnections: %w", err)
	}

	return connections, nil
}

// Delete connections -.
func (uc *CtrlUseCase) DeleteConnections(ctx context.Context, cluster entity.Cluster, connections []entity.Connection, clusterCred entity.Credentials) error {
	err := uc.pipe.DeleteConnections(ctx, cluster, connections, clusterCred)
	if err != nil {
		return fmt.Errorf("CtrlUseCase - DeleteConnections - uc.pipe.DeleteConnections: %w", err)
	}

	return nil
}

// Backup -.
func (uc *CtrlUseCase) RunBackup(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, lockCode, outputPath string) (string, error) {
	fullPath := path.Join(outputPath, fmt.Sprintf("%s_%s.dt", time.Now().Format("02_01_2006_15_04_05"), infobase.Name))

	err := uc.backup.RunBackup(ctx, cluster, infobase, clusterCred, lockCode, fullPath)
	if err != nil {
		return "", fmt.Errorf("CtrlUseCase - RunBackup - uc.backup.RunBackup: %w", err)
	}

	return fullPath, nil
}
