package walletchanging

import (
	"travelWallet/internal/usecase"
)

type WalletRemoveAllCurrenciesUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletRemoveAllCurrenciesUsecase(repo usecase.WalletRepository) *WalletRemoveAllCurrenciesUsecase {
	return &WalletRemoveAllCurrenciesUsecase{repo: repo}
}

func (w *WalletRemoveAllCurrenciesUsecase) RemoveAllCurrencies(walletId string) error {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return err
	}

	wallet.RemoveAllForeignCurrencies()

	err = w.repo.Save(wallet)
	if err != nil {
		return err
	}

	return nil
}
