package walletreading

import "travelWallet/internal/usecase"

type WalletGetBaseCurrencyUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletGetBaseCurrencyUsecase(repo usecase.WalletRepository) *WalletGetBaseCurrencyUsecase {
	return &WalletGetBaseCurrencyUsecase{repo: repo}
}

func (w *WalletGetBaseCurrencyUsecase) GetBaseCurrency(walletId string) (string, error) {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return "", err
	}

	return wallet.GetBaseCurrencyCode(), nil
}
