package walletchanging

import (
	"travelWallet/internal/domain"
	"travelWallet/internal/usecase"
)

type WalletWithdrawUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletWithdrawUsecase(repo usecase.WalletRepository) *WalletWithdrawUsecase {
	return &WalletWithdrawUsecase{repo: repo}
}

func (w *WalletWithdrawUsecase) Withdraw(walletId string, amount domain.Currency) error {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return err
	}

	err = wallet.Withdraw(amount)
	if err != nil {
		return err
	}

	err = w.repo.Save(wallet)
	if err != nil {
		return err
	}

	return nil
}
