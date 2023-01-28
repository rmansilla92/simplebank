package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rmansilla92/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, FromAccount Account, ToAccount Account) Transfer {
	arg := CreateTransferParams {
		FromAccountID: FromAccount.ID,
		ToAccountID: ToAccount.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	expectedTransfer := createRandomTransfer(t, fromAccount, toAccount)

	transfer, err := testQueries.GetTransfer(context.Background(), expectedTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, expectedTransfer.ID, transfer.ID)
	require.Equal(t, expectedTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, expectedTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, expectedTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, expectedTransfer.CreatedAt, transfer.CreatedAt, time.Second)

}

func TestUpdateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	expectedTransfer := createRandomTransfer(t, fromAccount, toAccount)

	arg := UpdateTransferParams {
		ID: expectedTransfer.ID,
		FromAccountID: fromAccount.ID, 
		ToAccountID: toAccount.ID,
		Amount: util.RandomMoney(),
	}

	updatedTransfer, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)

	require.Equal(t, expectedTransfer.ID, updatedTransfer.ID)
	require.Equal(t, expectedTransfer.FromAccountID, updatedTransfer.FromAccountID)
	require.Equal(t, expectedTransfer.ToAccountID, updatedTransfer.ToAccountID)
	require.Equal(t, arg.Amount, updatedTransfer.Amount)
	require.WithinDuration(t, expectedTransfer.CreatedAt, updatedTransfer.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount, toAccount)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		fromAccount := createRandomAccount(t)
		toAccount := createRandomAccount(t)
		createRandomTransfer(t, fromAccount, toAccount)
	}

	arg := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg) 
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}