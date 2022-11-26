package db

import (
	"context"
	"simplebank/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	require.NotZero(t, entry2.ID)
	require.NotZero(t, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     6,
		Offset:    4,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 6)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, args.AccountID)
	}
}
