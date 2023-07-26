// nolint
package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/antonmisa/1cctl_cli/internal/entity"
	"github.com/antonmisa/1cctl_cli/internal/usecase"
	"github.com/antonmisa/1cctl_cli/internal/usecase/mocks"
)

func TestClusterByName(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		clConn        string
		clName        string
		cls           []entity.Cluster
		respError     string
		pipeMockError error
	}{
		{
			name:   "Success",
			ctx:    context.Background(),
			clConn: "localhost:1545",
			clName: "localhost:1234",
			cls: []entity.Cluster{
				{
					ID:   "1",
					Host: "localhost",
					Port: "1234",
					Name: "test",
				},
			},
		},
		{
			name:   "Not found",
			ctx:    context.Background(),
			clConn: "localhost:1545",
			clName: "test",
			cls: []entity.Cluster{
				{
					ID:   "1",
					Host: "localhost",
					Port: "1234",
					Name: "bad",
				},
			},
			respError: "cluster with name test not found",
		},
		{
			name:          "Pipe error",
			ctx:           context.Background(),
			clConn:        "localhost:1545",
			cls:           make([]entity.Cluster, 0),
			respError:     ": unexpected error",
			pipeMockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrlPipeMock := mocks.NewCtrlPipe(t)
			ctrlBackupMock := mocks.NewCtrlBackup(t)

			ctrlPipeMock.On("GetClusters", mock.MatchedBy(func(ctx context.Context) bool { return true })).
				Return(tc.cls, tc.pipeMockError).
				Once()

			ctrl := usecase.New(ctrlPipeMock, ctrlBackupMock)

			cl, err := ctrl.ClusterByName(tc.ctx, tc.clName)

			if tc.respError == "" {
				require.NoError(t, err)

				require.Equal(t, fmt.Sprintf("%s:%s", cl.Host, cl.Port), tc.clName)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Equal(t, cl, entity.Cluster{})
			}
		})
	}
}

func TestInfobaseByName(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cl            entity.Cluster
		ibName        string
		clCred        entity.Credentials
		ibs           []entity.Infobase
		respError     string
		pipeMockError error
	}{
		{
			name:   "Success",
			ctx:    context.Background(),
			cl:     entity.Cluster{ID: "123"},
			ibName: "TEST",
			clCred: entity.Credentials{},
			ibs: []entity.Infobase{
				{
					ID:   "1",
					Name: "TEST",
				},
			},
		},
		{
			name:   "Not found",
			ctx:    context.Background(),
			cl:     entity.Cluster{ID: "123"},
			ibName: "test",
			clCred: entity.Credentials{},
			ibs: []entity.Infobase{
				{
					ID:   "1",
					Name: "bad",
				},
			},
			respError: "infobase with name test not found",
		},
		{
			name:          "Pipe error",
			ctx:           context.Background(),
			cl:            entity.Cluster{ID: "123"},
			ibs:           make([]entity.Infobase, 0),
			clCred:        entity.Credentials{},
			respError:     ": unexpected error",
			pipeMockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrlPipeMock := mocks.NewCtrlPipe(t)
			ctrlBackupMock := mocks.NewCtrlBackup(t)

			ctrlPipeMock.On("GetInfobases",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("entity.Cluster"),
				mock.AnythingOfType("entity.Credentials")).
				Return(tc.ibs, tc.pipeMockError).
				Once()

			ctrl := usecase.New(ctrlPipeMock, ctrlBackupMock)

			ib, err := ctrl.InfobaseByName(tc.ctx, tc.cl, tc.ibName, tc.clCred)

			if tc.respError == "" {
				require.NoError(t, err)

				require.Equal(t, ib.Name, tc.ibName)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Equal(t, ib, entity.Infobase{})
			}
		})
	}
}

func TestSessions(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cl            entity.Cluster
		ib            entity.Infobase
		cred          entity.Credentials
		ss            []entity.Session
		respError     string
		pipeMockError error
	}{
		{
			name: "Success all",
			ctx:  context.Background(),
			cl:   entity.Cluster{ID: "123"},
			ib:   entity.Infobase{},
			cred: entity.Credentials{},
			ss: []entity.Session{
				{
					ID: "1",
				},
			},
		},
		{
			name: "Success with infobase",
			ctx:  context.Background(),
			cl:   entity.Cluster{ID: "123"},
			ib:   entity.Infobase{ID: "123"},
			cred: entity.Credentials{},
			ss: []entity.Session{
				{
					ID: "1",
				},
			},
		},
		{
			name:          "Pipe error",
			ctx:           context.Background(),
			cl:            entity.Cluster{ID: "123"},
			ib:            entity.Infobase{ID: "123"},
			cred:          entity.Credentials{},
			ss:            make([]entity.Session, 0),
			respError:     ": unexpected error",
			pipeMockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrlPipeMock := mocks.NewCtrlPipe(t)
			ctrlBackupMock := mocks.NewCtrlBackup(t)

			ctrlPipeMock.On("GetSessions",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("entity.Cluster"),
				mock.AnythingOfType("entity.Infobase"),
				mock.AnythingOfType("entity.Credentials")).
				Return(tc.ss, tc.pipeMockError).
				Once()

			ctrl := usecase.New(ctrlPipeMock, ctrlBackupMock)

			ss, err := ctrl.Sessions(tc.ctx, tc.cl, tc.ib, tc.cred)

			if tc.respError == "" {
				require.NoError(t, err)

				require.NotEmpty(t, ss)
				require.Equal(t, ss, tc.ss)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, ss)
			}
		})
	}
}
