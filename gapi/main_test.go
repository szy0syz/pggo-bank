package gapi

import (
	"github.com/stretchr/testify/require"
	db "github.com/szy0syz/pggo-bank/db/sqlc"
	"github.com/szy0syz/pggo-bank/util"
	"github.com/szy0syz/pggo-bank/worker"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
