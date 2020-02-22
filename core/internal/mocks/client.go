// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	assets "nulink/core/assets"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	decimal "github.com/shopspring/decimal"

	eth "nulink/core/eth"

	ethereum "github.com/ethereum/go-ethereum"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetAggregatorPrice provides a mock function with given fields: address, precision
func (_m *Client) GetAggregatorPrice(address common.Address, precision int32) (decimal.Decimal, error) {
	ret := _m.Called(address, precision)

	var r0 decimal.Decimal
	if rf, ok := ret.Get(0).(func(common.Address, int32) decimal.Decimal); ok {
		r0 = rf(address, precision)
	} else {
		r0 = ret.Get(0).(decimal.Decimal)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, int32) error); ok {
		r1 = rf(address, precision)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAggregatorRound provides a mock function with given fields: address
func (_m *Client) GetAggregatorRound(address common.Address) (*big.Int, error) {
	ret := _m.Called(address)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address) *big.Int); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockByNumber provides a mock function with given fields: hex
func (_m *Client) GetBlockByNumber(hex string) (eth.BlockHeader, error) {
	ret := _m.Called(hex)

	var r0 eth.BlockHeader
	if rf, ok := ret.Get(0).(func(string) eth.BlockHeader); ok {
		r0 = rf(hex)
	} else {
		r0 = ret.Get(0).(eth.BlockHeader)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChainID provides a mock function with given fields:
func (_m *Client) GetChainID() (*big.Int, error) {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetERC20Balance provides a mock function with given fields: address, contractAddress
func (_m *Client) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	ret := _m.Called(address, contractAddress)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address, common.Address) *big.Int); ok {
		r0 = rf(address, contractAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address) error); ok {
		r1 = rf(address, contractAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEthBalance provides a mock function with given fields: address
func (_m *Client) GetEthBalance(address common.Address) (*assets.Eth, error) {
	ret := _m.Called(address)

	var r0 *assets.Eth
	if rf, ok := ret.Get(0).(func(common.Address) *assets.Eth); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Eth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLogs provides a mock function with given fields: q
func (_m *Client) GetLogs(q ethereum.FilterQuery) ([]eth.Log, error) {
	ret := _m.Called(q)

	var r0 []eth.Log
	if rf, ok := ret.Get(0).(func(ethereum.FilterQuery) []eth.Log); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]eth.Log)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ethereum.FilterQuery) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNonce provides a mock function with given fields: address
func (_m *Client) GetNonce(address common.Address) (uint64, error) {
	ret := _m.Called(address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(common.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxReceipt provides a mock function with given fields: hash
func (_m *Client) GetTxReceipt(hash common.Hash) (*eth.TxReceipt, error) {
	ret := _m.Called(hash)

	var r0 *eth.TxReceipt
	if rf, ok := ret.Get(0).(func(common.Hash) *eth.TxReceipt); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*eth.TxReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendRawTx provides a mock function with given fields: hex
func (_m *Client) SendRawTx(hex string) (common.Hash, error) {
	ret := _m.Called(hex)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(string) common.Hash); ok {
		r0 = rf(hex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeToLogs provides a mock function with given fields: channel, q
func (_m *Client) SubscribeToLogs(channel chan<- eth.Log, q ethereum.FilterQuery) (eth.Subscription, error) {
	ret := _m.Called(channel, q)

	var r0 eth.Subscription
	if rf, ok := ret.Get(0).(func(chan<- eth.Log, ethereum.FilterQuery) eth.Subscription); ok {
		r0 = rf(channel, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(eth.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(chan<- eth.Log, ethereum.FilterQuery) error); ok {
		r1 = rf(channel, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeToNewHeads provides a mock function with given fields: channel
func (_m *Client) SubscribeToNewHeads(channel chan<- eth.BlockHeader) (eth.Subscription, error) {
	ret := _m.Called(channel)

	var r0 eth.Subscription
	if rf, ok := ret.Get(0).(func(chan<- eth.BlockHeader) eth.Subscription); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(eth.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(chan<- eth.BlockHeader) error); ok {
		r1 = rf(channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
