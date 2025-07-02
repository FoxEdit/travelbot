package walletchanging

import (
	"travelWallet/internal/domain"
	"travelWallet/internal/usecase"
)

type WalletDepositUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletDepositUsecase(repo usecase.WalletRepository) *WalletDepositUsecase {
	return &WalletDepositUsecase{repo: repo}
}

func (w *WalletDepositUsecase) Deposit(walletId string, amount domain.Currency) error {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return err
	}

	err = wallet.Deposit(amount)
	if err != nil {
		return err
	}

	err = w.repo.Save(wallet)
	if err != nil {
		return err
	}

	return nil
}
