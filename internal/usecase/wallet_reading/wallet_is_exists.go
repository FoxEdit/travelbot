package walletreading

import (
	"errors"
	"os"
	"travelWallet/internal/usecase"
)

type WalletIsExistsUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletIsExistsUsecase(repo usecase.WalletRepository) *WalletIsExistsUsecase {
	return &WalletIsExistsUsecase{repo: repo}
}

func (w *WalletIsExistsUsecase) IsExists(walletId string) (bool, error) {
	_, err := w.repo.FindByID(walletId)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}
