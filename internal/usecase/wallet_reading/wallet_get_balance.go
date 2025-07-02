package walletreading

import (
	"travelWallet/internal/usecase"

	"github.com/shopspring/decimal"
)

type WalletGetBalanceUsecase struct {
	repo usecase.WalletRepository
}

func NewWalletGetBalanceUsecase(repo usecase.WalletRepository) *WalletGetBalanceUsecase {
	return &WalletGetBalanceUsecase{repo: repo}
}

func (w *WalletGetBalanceUsecase) GetBalance(walletId string) (decimal.Decimal, error) {
	wallet, err := w.repo.FindByID(walletId)
	if err != nil {
		return decimal.Zero, err
	}

	return wallet.GetRawBalance(), nil
}
