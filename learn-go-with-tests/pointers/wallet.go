package main

import (
	"errors"
	"fmt"
)

// Bitcoin 比特币
type Bitcoin int

var InsufficientFundsError = errors.New("insufficient funds")

// Wallet 钱包
type Wallet struct {
	bitcoin Bitcoin
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

// Deposit 存款
func (w *Wallet) Deposit(amount Bitcoin) {
	w.bitcoin += amount
}

// Balance 查看余额
func (w *Wallet) Balance() Bitcoin {
	return w.bitcoin
}

// Withdraw 取款, 余额不足时返回 InsufficientFundsError 错误
func (w *Wallet) Withdraw(amount Bitcoin) error {
	if w.bitcoin < amount {
		return InsufficientFundsError
	}
	w.bitcoin -= amount
	return nil
}
