package pipe

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"testing"

	"github.com/antonmisa/1cctl_cli/internal/entity"
	"github.com/antonmisa/1cctl_cli/pkg/pipe/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type FakeReadCloser struct {
	body []byte
	pos  int
}

func NewFakeSession() *FakeReadCloser {
	text := `session : 1111-3434-5656 
				infobase: 3333-4444
				connection: 3-4-5-6
				process: 1-1-1-1
				user-name: Тестовый пользователь
				host: test-ic
				app-id: 1cv8

				session : 2222-3434-5656 
				infobase: 3333-4444
				connection: 3-4-5-7
				process: 1-1-1-2
				user-name: Тестовый пользователь 1
				host: test-ic-1
				app-id: 1cv8
				
				session : 3333-3434-5656 
				infobase: 1111-4444
				connection: 1-4-5-7
				process: 1-2-1-2
				user-name: неизвестный пользователь
				host: test-ic-2
				app-id: 1cv8`

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeInfobase() *FakeReadCloser {
	text := `infobase : 1212-3434-5656 
			 name: test
			 descr:

			 infobase : 1111-2222-3333 
			  name: test_ib
			 descr: "test desc" `

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeCluster() *FakeReadCloser {
	text := `cluster : 1212-3434-5656 
			 host: localhost 
			 port: 1234 
			 name: "test"

			 cluster : 1111-2222-3333 
			  host: localhost.tnx.ru    
			 port: 1545 
			  name: "test cluster" `

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func (t *FakeReadCloser) Read(p []byte) (n int, err error) {
	if t.pos >= len(t.body) {
		return 0, io.EOF
	}

	t.pos = copy(p, t.body[t.pos:])

	if t.pos >= len(t.body) {
		return t.pos, io.EOF
	} else {
		return t.pos, nil
	}
}

func (t *FakeReadCloser) Close() error {
	return nil
}

var _ io.ReadCloser = (*FakeReadCloser)(nil)

func TestGetClusters(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		cls           []entity.Cluster
		respError     string
		pipeMockError error
	}{
		{
			name: "Success",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeCluster(),
			cls: []entity.Cluster{
				{
					ID:   "1212-3434-5656",
					Host: "localhost",
					Port: "1234",
					Name: "\"test\"",
				},
				{
					ID:   "1111-2222-3333",
					Host: "localhost.tnx.ru",
					Port: "1545",
					Name: "\"test cluster\"",
				},
			},
		},
		{
			name:          "Error no command",
			ctx:           context.Background(),
			cs:            "localhost:1545",
			cmd:           &exec.Cmd{},
			stdout:        NewFakeCluster(),
			cls:           make([]entity.Cluster, 0),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Once()

			ctrl := New(pipeMock, tc.cs)

			cls, err := ctrl.GetClusters(tc.ctx)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, cls)
				require.ElementsMatch(t, cls, tc.cls)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, cls)
			}
		})
	}
}

func TestGetInfobases(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cl            entity.Cluster
		cred          entity.Credentials
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		ibs           []entity.Infobase
		respError     string
		pipeMockError error
	}{
		{
			name: "Success wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			cred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeInfobase(),
			ibs: []entity.Infobase{
				{
					ID:   "1212-3434-5656",
					Name: "test",
					Desc: "",
				},
				{
					ID:   "1111-2222-3333",
					Name: "test_ib",
					Desc: "\"test desc\"",
				},
			},
		},
		{
			name: "Success w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeInfobase(),
			ibs: []entity.Infobase{
				{
					ID:   "1212-3434-5656",
					Name: "test",
					Desc: "",
				},
				{
					ID:   "1111-2222-3333",
					Name: "test_ib",
					Desc: "\"test desc\"",
				},
			},
		},
		{
			name:          "Error no command",
			ctx:           context.Background(),
			cs:            "localhost:1545",
			cmd:           &exec.Cmd{},
			stdout:        NewFakeInfobase(),
			ibs:           make([]entity.Infobase, 0),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Once()

			ctrl := New(pipeMock, tc.cs)

			ibs, err := ctrl.GetInfobases(tc.ctx, tc.cl, tc.cred)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, ibs)
				require.ElementsMatch(t, ibs, tc.ibs)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, ibs)
			}
		})
	}
}

func TestGetSessions(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cl            entity.Cluster
		ib            entity.Infobase
		cred          entity.Credentials
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		res           []entity.Session
		respError     string
		pipeMockError error
	}{
		{
			name: "Success w empty ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:   entity.Infobase{},
			cred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "Тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "Тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
				{
					ID:           "3333-3434-5656",
					InfobaseID:   "1111-4444",
					ConnectionID: "1-4-5-7",
					ProcessID:    "1-2-1-2",
					UserName:     "неизвестный пользователь",
					Host:         "test-ic-2",
					AppID:        "1cv8",
				},
			},
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			cred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "Тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "Тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
			},
		},
		{
			name: "Success w empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "Тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "Тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
				{
					ID:           "3333-3434-5656",
					InfobaseID:   "1111-4444",
					ConnectionID: "1-4-5-7",
					ProcessID:    "1-2-1-2",
					UserName:     "неизвестный пользователь",
					Host:         "test-ic-2",
					AppID:        "1cv8",
				},
			},
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "Тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "Тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
			},
		},

		{
			name:          "Error no command",
			ctx:           context.Background(),
			cs:            "localhost:1545",
			cl:            entity.Cluster{},
			ib:            entity.Infobase{},
			cred:          entity.Credentials{},
			cmd:           &exec.Cmd{},
			stdout:        NewFakeSession(),
			res:           make([]entity.Session, 0),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Once()

			ctrl := New(pipeMock, tc.cs)

			res, err := ctrl.GetSessions(tc.ctx, tc.cl, tc.ib, tc.cred)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, res)
				require.ElementsMatch(t, res, tc.res)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, res)
			}
		})
	}
}

func TestDisableSessions(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cl            entity.Cluster
		ib            entity.Infobase
		clCred        entity.Credentials
		ibCred        entity.Credentials
		code          string
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		respError     string
		pipeMockError error
	}{
		{
			name: "Error infobase empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:     entity.Infobase{},
			clCred: entity.Credentials{},
			ibCred: entity.Credentials{},
			code:   "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout:        NewFakeSession(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			ibCred: entity.Credentials{},
			code:   "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code: "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout:        NewFakeSession(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code: "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ib: entity.Infobase{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			ibCred:        entity.Credentials{},
			code:          "12345",
			cmd:           &exec.Cmd{},
			stdout:        NewFakeSession(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock, tc.cs)

			err := ctrl.DisableSessions(tc.ctx, tc.cl, tc.ib, tc.clCred, tc.ibCred, tc.code)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestEnableSessions(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cl            entity.Cluster
		ib            entity.Infobase
		clCred        entity.Credentials
		ibCred        entity.Credentials
		code          string
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		respError     string
		pipeMockError error
	}{
		{
			name: "Error infobase empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:     entity.Infobase{},
			clCred: entity.Credentials{},
			ibCred: entity.Credentials{},
			code:   "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout:        NewFakeSession(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			ibCred: entity.Credentials{},
			code:   "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code: "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout:        NewFakeSession(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code: "12345",
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ib: entity.Infobase{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			ibCred:        entity.Credentials{},
			code:          "12345",
			cmd:           &exec.Cmd{},
			stdout:        NewFakeSession(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock, tc.cs)

			err := ctrl.EnableSessions(tc.ctx, tc.cl, tc.ib, tc.clCred, tc.ibCred, tc.code)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cl            entity.Cluster
		s             entity.Session
		clCred        entity.Credentials
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		respError     string
		pipeMockError error
	}{
		{
			name: "Error session is empty empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s:      entity.Session{},
			clCred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout:        NewFakeSession(),
			respError:     ": session is empty",
			pipeMockError: ErrSessionIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Session{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Session{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout:        NewFakeSession(),
			respError:     ": session is empty",
			pipeMockError: ErrSessionIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Session{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			s: entity.Session{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			cmd:           &exec.Cmd{},
			stdout:        NewFakeSession(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock, tc.cs)

			err := ctrl.DeleteSession(tc.ctx, tc.cl, tc.s, tc.clCred)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestDeleteSessions(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		cs            string
		cl            entity.Cluster
		ss            []entity.Session
		clCred        entity.Credentials
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		respError     string
		pipeMockError error
	}{
		{
			name: "Success session is empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss:     make([]entity.Session, 0),
			clCred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Session{
				{
					ID: "3333-4444",
				},
				{
					ID: "5555-4444",
				},
			},
			clCred: entity.Credentials{},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Success empty sessions w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Session{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Success w sessions w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Session{
				{
					ID: "3333-4444",
				},
				{
					ID: "5555-4444",
				},
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			cmd: &exec.Cmd{
				Path: "C:\\Windows\\System32\\ping.exe",
			},
			stdout: NewFakeSession(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ss: []entity.Session{
				{
					ID: "12",
				},
			},
			clCred:        entity.Credentials{},
			cmd:           &exec.Cmd{},
			stdout:        NewFakeSession(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipeMock := mocks.NewInterface(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(tc.cmd, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock, tc.cs)

			err := ctrl.DeleteSessions(tc.ctx, tc.cl, tc.ss, tc.clCred)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}
