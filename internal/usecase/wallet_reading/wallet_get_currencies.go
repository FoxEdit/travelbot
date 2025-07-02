package walletreading

import "travelWallet/internal/usecase"

type WalletGetCurrenciesUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletGetCurrenciesUsecase(repo usecase.WalletRepository) *WalletGetCurrenciesUsecase {
	return &WalletGetCurrenciesUsecase{repo: repo}
}

func (w *WalletGetCurrenciesUsecase) GetCurrencies(walletId string) ([]string, error) {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return []string{}, err
	}

	return wallet.GetAllForeignCurrencies(), nil
}
