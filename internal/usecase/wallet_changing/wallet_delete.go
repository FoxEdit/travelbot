package walletchanging

import (
	"travelWallet/internal/usecase"
)

type WalletDeleteUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletDeleteUsecase(repo usecase.WalletRepository) *WalletDeleteUsecase {
	return &WalletDeleteUsecase{repo: repo}
}

func (w *WalletDeleteUsecase) Delete(walletId string, baseCurrency string) error {
	err := w.repo.DeleteByID(walletId)
	if err != nil {
		return err
	}

	return nil
}
