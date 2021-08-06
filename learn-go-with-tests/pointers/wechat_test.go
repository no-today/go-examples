package main

import "testing"

func TestWallet(t *testing.T) {
	t.Run("存款", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(10)

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("取款", func(t *testing.T) {
		wallet := Wallet{20}
		err := wallet.Withdraw(10)

		assertBalance(t, wallet, Bitcoin(10))
		assertNotError(t, err)
	})

	t.Run("取款:余额不足", func(t *testing.T) {
		wallet := Wallet{20}
		err := wallet.Withdraw(100)

		assertBalance(t, wallet, wallet.bitcoin)
		assertError(t, err, InsufficientFundsError.Error())
	})
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertError(t *testing.T, err error, want string) {
	if err == nil {
		t.Fatal("wanted an error but didn't get one")
	}

	if err.Error() != want {
		t.Errorf("got %s, want %s", err.Error(), want)
	}
}

func assertNotError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("no errors expected")
	}
}
