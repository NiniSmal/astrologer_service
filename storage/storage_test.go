package storage

import (
	"astrologerService/dbtest"
	"astrologerService/entity"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage_SaveData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conn := dbtest.DBConnection(t)
	ns := NewStorage(conn)

	createdAt := time.Now().UTC().Round(time.Millisecond)
	data := entity.Metadata{
		Date:           createdAt.String(),
		Explanation:    uuid.NewString(),
		HdURL:          uuid.NewString(),
		MediaType:      uuid.NewString(),
		ServiceVersion: uuid.NewString(),
		Title:          uuid.NewString(),
		URL:            uuid.NewString(),
	}
	err := ns.SaveData(ctx, data)
	require.NoError(t, err)

	dataDB, err := ns.GetMetaData(ctx, createdAt.String())
	require.NoError(t, err)
	require.Equal(t, data, dataDB)
}

func TestStorage_GetMetaData_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	conn := dbtest.DBConnection(t)
	ns := NewStorage(conn)
	_, err := ns.GetMetaData(ctx, "123")
	require.Error(t, err)
}

func TestStorage_GetMetaDataBases(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conn := dbtest.DBConnection(t)
	ns := NewStorage(conn)

	createdAt := time.Now().UTC().Round(time.Millisecond)
	data := entity.Metadata{
		Date:           createdAt.String(),
		Explanation:    uuid.NewString(),
		HdURL:          uuid.NewString(),
		MediaType:      uuid.NewString(),
		ServiceVersion: uuid.NewString(),
		Title:          uuid.NewString(),
		URL:            uuid.NewString(),
	}
	err := ns.SaveData(ctx, data)
	require.NoError(t, err)

	dataBases, err := ns.GetMetaDataBases(ctx)
	require.NoError(t, err)
	require.Contains(t, dataBases, data)
}
