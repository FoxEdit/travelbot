package walletchanging

import (
	"travelWallet/internal/usecase"
)

type WalletAddCurrencyUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletAddCurrencyUsecase(repo usecase.WalletRepository) *WalletAddCurrencyUsecase {
	return &WalletAddCurrencyUsecase{repo: repo}
}

func (w *WalletAddCurrencyUsecase) AddCurrency(walletId string, newCurrency string) error {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return err
	}

	wallet.AddForeignCurrency(newCurrency)

	err = w.repo.Save(wallet)
	if err != nil {
		return err
	}

	return nil
}
