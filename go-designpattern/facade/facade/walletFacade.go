package facade

import "fmt"

type WalletFacade struct {
	account      *Account
	wallet       *Wallet
	securityCode *SecurityCode
	notification *Notification
	ledger       *Ledger
}

func NewWalletFacade(accountID string, code int) *WalletFacade {
	fmt.Println("Starting create account")
	walletFacade := &WalletFacade{
		account:      newAccount(accountID),
		wallet:       newWallet(),
		securityCode: newSecurityCode(code),
		notification: &Notification{},
		ledger:       &Ledger{},
	}
	fmt.Println("Account created")
	return walletFacade
}

func (w *WalletFacade) AddMoneyToWallet(accountID string, securityCode, amount int) error {
	fmt.Println("Starting add money to wallet")
	if err := w.account.checkAccount(accountID); err != nil {
		return err
	}
	if err := w.securityCode.checkCode(securityCode); err != nil {
		return err
	}
	w.wallet.creditBalance(amount)
	w.notification.sendWalletCreditNotification()
	w.ledger.makeEntry(accountID, "credit", amount)
	return nil
}

func (w *WalletFacade) DeductMoneyFromWallet(accountID string, securityCode, amount int) error {
	fmt.Println("Starting debit money from wallet")
	if err := w.account.checkAccount(accountID); err != nil {
		return err
	}
	if err := w.securityCode.checkCode(securityCode); err != nil {
		return err
	}
	w.notification.sendWalletDebitNotification()
	w.ledger.makeEntry(accountID, "credit", amount)
	return nil
}
