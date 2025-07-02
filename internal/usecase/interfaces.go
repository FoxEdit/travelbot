package usecase

import "travelWallet/internal/domain"

// ===== wallet interfaces =====

type WalletRepository interface {
	CreateByID(id string, baseCurrency string) error
	DeleteByID(id string) error

	Save(wallet *domain.Wallet) error
	FindByID(id string) (*domain.Wallet, error)
}

// ===== wallet interfaces =====
// ===== converter interfaces =====

type CurrencyConverter interface {
	ConvertToMany(currency domain.Currency, anothers ...string) ([]domain.Currency, error)
	ConvertToAnother(currency domain.Currency, another string) (domain.Currency, error)
}
