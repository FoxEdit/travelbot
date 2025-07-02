package walletchanging

import (
	"travelWallet/internal/usecase"
)

type WalletCreateUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletCreateUsecase(repo usecase.WalletRepository) *WalletCreateUsecase {
	return &WalletCreateUsecase{repo: repo}
}

func (w *WalletCreateUsecase) Create(walletId string, baseCurrency string) error {
	err := w.repo.CreateByID(walletId, baseCurrency)
	if err != nil {
		return err
	}

	return nil
}
