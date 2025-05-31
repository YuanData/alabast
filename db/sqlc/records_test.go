package db

import (
	"alabast/util"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomRecord(t *testing.T, account User) Record {
	str := util.RandomString(8)
	quoted, err := json.Marshal(str)
	require.NoError(t, err)

	arg := CreateRecordParams{
		Username: account.Username,
		Content:  json.RawMessage(quoted),
	}

	entry, err := testQueries.CreateRecord(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.Username, entry.Username)
	require.Equal(t, arg.Content, entry.Content)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateRecord(t *testing.T) {
	account := createRandomUser(t)
	createRandomRecord(t, account)
}

func TestGetRecord(t *testing.T) {
	account := createRandomUser(t)
	entry1 := createRandomRecord(t, account)
	entry2, err := testQueries.GetRecord(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Username, entry2.Username)
	require.Equal(t, entry1.Content, entry2.Content)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListRecords(t *testing.T) {
	account := createRandomUser(t)
	for i := 0; i < 2; i++ {
		createRandomRecord(t, account)
	}

	arg := ListRecordsParams{
		Username: account.Username,
		Limit:    1,
		Offset:   1,
	}

	entries, err := testQueries.ListRecords(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 1)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.Username, entry.Username)
	}
}
