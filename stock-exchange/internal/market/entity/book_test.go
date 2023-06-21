package entity

import (
	"sync"
	"testing"

	testify "github.com/stretchr/testify/assert"
)

func TestBuyAsset(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 10)
	investor1.AddAssetPosition(investorAssetPosition)

	wg := sync.WaitGroup{}
	orderChanIn := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChanIn, orderChanOut, &wg)
	go book.Trade()

	// add sell order
	wg.Add(1)
	order1 := NewOrder("1", investor1, asset1, 5, 5, "SELL")
	orderChanIn <- order1

	// add buy order
	order2 := NewOrder("2", investor2, asset1, 5, 5, "BUY")
	orderChanIn <- order2
	wg.Wait()

	assert := testify.New(t)
	assert.Equal("CLOSED", order1.Status, "Order 1 should be closed")
	assert.Equal(0, order1.PendingShares, "Order 1 should have 0 PendingShares")
	assert.Equal("CLOSED", order2.Status, "Order 2 should be closed")
	assert.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(5, investorAssetPosition.Shares, "Investor 1 should have 5 shares of asset 1")
	assert.Equal(5, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 5 shares of asset 1")
}

func TestBuyAssetWithDifferentAssents(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)
	asset2 := NewAsset("asset2", "Asset 2", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 10)
	investor1.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewInvestorAssetPosition("asset2", 10)
	investor2.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChanIn := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChanIn, orderChanOut, &wg)
	go book.Trade()

	order1 := NewOrder("1", investor1, asset1, 5, 5, "SELL")
	orderChanIn <- order1

	order2 := NewOrder("2", investor2, asset2, 5, 5, "BUY")
	orderChanIn <- order2

	assert := testify.New(t)
	assert.Equal("OPEN", order1.Status, "Order 1 should be closed")
	assert.Equal(5, order1.PendingShares, "Order 1 should have 5 PendingShares")
	assert.Equal("OPEN", order2.Status, "Order 2 should be closed")
	assert.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}

func TestBuyPartialAsset(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")
	investor3 := NewInvestor("3", "Investor 3")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 3)
	investor1.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewInvestorAssetPosition("asset1", 5)
	investor3.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChanIn := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChanIn, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("1", investor2, asset1, 5, 5.0, "BUY")
	orderChanIn <- order2

	// investidor 1 quer vender 3 shares
	order1 := NewOrder("2", investor1, asset1, 3, 5.0, "SELL")
	orderChanIn <- order1

	assert := testify.New(t)
	go func() {
		for range orderChanOut {
		}
	}()

	wg.Wait()

	// assert := assert.New(t)
	assert.Equal("CLOSED", order1.Status, "Order 1 should be closed")
	assert.Equal(0, order1.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	assert.Equal(2, order2.PendingShares, "Order 2 should have 2 PendingShares")

	assert.Equal(0, investorAssetPosition.Shares, "Investor 1 should have 0 shares of asset 1")
	assert.Equal(3, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder("3", investor3, asset1, 2, 5.0, "SELL")
	orderChanIn <- order3
	wg.Wait()

	assert.Equal("CLOSED", order3.Status, "Order 3 should be closed")
	assert.Equal(0, order3.PendingShares, "Order 3 should have 0 PendingShares")

	assert.Equal("CLOSED", order2.Status, "Order 2 should be CLOSED")
	assert.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(2, len(book.Transactions), "Should have 2 transactions")
	assert.Equal(15.0, book.Transactions[0].Total, "Transactions should have price 15")
	assert.Equal(10.0, book.Transactions[1].Total, "Transactions should have price 10")
}

func TestBuyWithDifferentPrice(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")
	investor3 := NewInvestor("3", "Investor 3")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 3)
	investor1.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewInvestorAssetPosition("asset1", 5)
	investor3.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChanIn := make(chan *Order)

	orderChanOut := make(chan *Order)

	book := NewBook(orderChanIn, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("2", investor2, asset1, 5, 5.0, "BUY")
	orderChanIn <- order2

	// investidor 1 quer vender 3 shares
	order1 := NewOrder("1", investor1, asset1, 3, 4.0, "SELL")
	orderChanIn <- order1

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := testify.New(t)
	assert.Equal("CLOSED", order1.Status, "Order 1 should be closed")
	assert.Equal(0, order1.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	assert.Equal(2, order2.PendingShares, "Order 2 should have 2 PendingShares")

	assert.Equal(0, investorAssetPosition.Shares, "Investor 1 should have 0 shares of asset 1")
	assert.Equal(3, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder("3", investor3, asset1, 3, 4.5, "SELL")
	orderChanIn <- order3

	wg.Wait()

	assert.Equal("OPEN", order3.Status, "Order 3 should be open")
	assert.Equal(1, order3.PendingShares, "Order 3 should have 1 PendingShares")

	assert.Equal("CLOSED", order2.Status, "Order 2 should be CLOSED")
	assert.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	// assert.Equal(2, len(book.Transactions), "Should have 2 transactions")
	// assert.Equal(15.0, float64(book.Transactions[0].Total), "Transactions should have price 15")
	// assert.Equal(10.0, float64(book.Transactions[1].Total), "Transactions should have price 10")
}

func TestNoMatch(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 3)
	investor1.AddAssetPosition(investorAssetPosition)

	wg := sync.WaitGroup{}
	orderChanIn := make(chan *Order)

	orderChanOut := make(chan *Order)

	book := NewBook(orderChanIn, orderChanOut, &wg)
	go book.Trade()

	wg.Add(0)
	// investidor 1 quer vender 3 shares
	order1 := NewOrder("1", investor1, asset1, 3, 6.0, "SELL")
	orderChanIn <- order1

	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("2", investor2, asset1, 5, 5.0, "BUY")
	orderChanIn <- order2

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := testify.New(t)
	assert.Equal("OPEN", order1.Status, "Order 1 should be closed")
	assert.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	assert.Equal(3, order1.PendingShares, "Order 1 should have 3 PendingShares")
	assert.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}
