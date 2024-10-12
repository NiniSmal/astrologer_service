package dbtest

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func DBConnection(t *testing.T) *pgxpool.Pool {
	t.Helper()

	connString := os.Getenv("POSTGRES_CONN_STRING")
	if connString == "" {
		connString = "postgres://postgres:dev@localhost:8091/postgres?sslmode=disable"
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, connString)
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close()
		require.NoError(t, err)
	})
	conn.Ping(ctx)
	require.NoError(t, err)

	return conn
}
