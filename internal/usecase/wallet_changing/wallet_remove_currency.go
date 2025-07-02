package walletchanging

import (
	"travelWallet/internal/usecase"
)

type WalletRemoveCurrencyUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletRemoveCurrencyUsecase(repo usecase.WalletRepository) *WalletRemoveCurrencyUsecase {
	return &WalletRemoveCurrencyUsecase{repo: repo}
}

func (w *WalletRemoveCurrencyUsecase) RemoveCurrency(walletId string, currencyToRemove string) error {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return err
	}

	wallet.RemoveForeignCurrency(currencyToRemove)

	err = w.repo.Save(wallet)
	if err != nil {
		return err
	}

	return nil
}
