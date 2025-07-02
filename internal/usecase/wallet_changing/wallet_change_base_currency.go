package walletchanging

import (
	"travelWallet/internal/usecase"
)

type WalletChangeBaseCurrencyUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletChangeBaseCurrencyUsecase(repo usecase.WalletRepository) *WalletChangeBaseCurrencyUsecase {
	return &WalletChangeBaseCurrencyUsecase{repo: repo}
}

func (w *WalletChangeBaseCurrencyUsecase) ChangeBaseCurrency(walletId string, newCurrency string) error {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return err
	}

	wallet.ChangeBaseCurrency(newCurrency)

	err = w.repo.Save(wallet)
	if err != nil {
		return err
	}

	return nil
}
