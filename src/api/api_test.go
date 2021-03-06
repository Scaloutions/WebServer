package api

/*
	TODO:
	getQuote
	cancelSetBuy
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const (
	TEST_URL                     = "http://localhost:8082/api/test"
	TEST_ACCOUNT_TRANSACTION_URL = "http://localhost:8082/api/accounttransaction"
	TEST_SYSTEM_EVENT_URL        = "http://localhost:8082/api/systemevent"
)

func activateHttpmock(url string) {

	httpmock.Activate()

	httpmock.RegisterResponder(
		"POST",
		url,
		httpmock.NewStringResponder(200, "ok"))
}

func activateMockAuditServer() {
	activateHttpmock(TEST_SYSTEM_EVENT_URL)
	activateHttpmock(TEST_ACCOUNT_TRANSACTION_URL)
}

func initializeAccountForTesting(amount float64) *Account {

	account := InitializeAccount("test123")
	transactionNum := 1
	Add(&account, amount, transactionNum)

	return &account
}

func buyStockForTesting(account *Account) {
	amount := float64(64)
	stock := "S"
	transactionNum := 2
	stockNum := float64(4)
	buyHelper(account, amount, stock, stockNum, transactionNum)
}

func commitBuyForTesting(account *Account) {
	buyStockForTesting(account)
	CommitBuy(account, 3)
}

func sellStockForTesting(account *Account) {
	amount := float64(64)
	stock := "S"
	transactionNum := 5
	stockNum := float64(4)
	sellHelper(account, stock, amount, transactionNum, stockNum)
}

func TestAdd(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	amount := 100.01
	account := initializeAccountForTesting(amount)
	assert.Equal(t, amount, account.Available)
	assert.Equal(t, amount, account.Balance)
}

func TestBuyWithoutQS(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	amount := 64.00
	account := initializeAccountForTesting(amount)
	assert.Equal(t, amount, account.Balance)
	assert.Equal(t, amount, account.Available)
	targetAmount := float64(0)
	stock := "S"
	stockNum := float64(10)
	transactionNum := 2
	buyHelper(account, amount, stock, stockNum, transactionNum)
	assert.Equal(t, targetAmount, account.Available)
	assert.Equal(t, amount, account.Balance) // 0 only when buy operation is committed
	assert.False(t, account.hasStock(stock, stockNum))
	// has stock only after buy is committed

}

func TestSell(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	commitBuyForTesting(account)
	actualBalance := float64(36)
	assert.True(t, account.hasStock("S", float64(4)))
	assert.Equal(t, actualBalance, account.Balance)
	sellHelper(account, "S", float64(64), 4, float64(4))
	assert.False(t, account.hasStock("S", float64(4))) // stock on hold
	assert.Equal(t, actualBalance, account.Balance)

}

func TestCommitBuy(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	buyStockForTesting(account)
	assert.Equal(t, float64(100), account.Balance)
	assert.Equal(t, float64(36), account.Available)

	CommitBuy(account, 3)
	assert.Equal(t, float64(36), account.Balance)
	assert.True(t, account.hasStock("S", float64(4)))

}

func TestCanCelBuy(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	buyStockForTesting(account)
	assert.Equal(t, float64(100), account.Balance)
	assert.Equal(t, float64(36), account.Available)

	CancelBuy(account, 5)
	assert.Equal(t, float64(100), account.Available)

}

func TestCommitSell(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	commitBuyForTesting(account)
	assert.Equal(t, float64(36), account.Balance)
	assert.Equal(t, float64(36), account.Available)
	assert.True(t, account.hasStock("S", float64(4)))

	sellStockForTesting(account)
	assert.False(t, account.hasStock("S", float64(4)))
	assert.Equal(t, float64(36), account.Balance)

	CommitSell(account, 7)
	assert.False(t, account.hasStock("S", float64(4)))
	assert.Equal(t, float64(100), account.Available)
	assert.Equal(t, float64(100), account.Balance)

}

func TestCancelSell(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	commitBuyForTesting(account)
	assert.Equal(t, float64(36), account.Balance)
	assert.Equal(t, float64(36), account.Available)
	assert.True(t, account.hasStock("S", float64(4)))

	sellStockForTesting(account)
	assert.False(t, account.hasStock("S", float64(4)))
	assert.Equal(t, float64(36), account.Balance)

	CancelSell(account, 8)
	assert.True(t, account.hasStock("S", float64(4)))
	assert.Equal(t, float64(36), account.Available)
	assert.Equal(t, float64(36), account.Balance)

}

func TestSetBuyAmount(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	assert.Equal(t, float64(100), account.Balance)
	assert.Equal(t, float64(100), account.Available)

	SetBuyAmount(account, "S", float64(64), 8)
	assert.Equal(t, float64(100), account.Balance)
	assert.Equal(t, float64(36), account.Available)

}

// func TestSetBuyTrigger(t *testing.T) {

// 	activateMockAuditServer()
// 	defer httpmock.DeactivateAndReset()

// 	account := initializeAccountForTesting(100)
// 	assert.Equal(t, float64(100), account.Balance)
// 	assert.Equal(t, float64(100), account.Available)

// 	setBuyAmount(account, "S", float64(64), 8)
// 	assert.Equal(t, float64(100), account.Balance)
// 	assert.Equal(t, float64(36), account.Available)

// 	setBuyTrigger(account, "S", )

// }
