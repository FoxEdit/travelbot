package app

import (
	"travelWallet/internal/usecase"
	walletchanging "travelWallet/internal/usecase/wallet_changing"
	walletreading "travelWallet/internal/usecase/wallet_reading"
)

type Application struct {
	Wallets  WalletServices
	Exchange ExchangeServices
}

type WalletServices struct {
	GetBalance      *walletreading.WalletGetBalanceUsecase
	GetForeign      *walletreading.WalletGetCurrenciesUsecase
	IsExists        *walletreading.WalletIsExistsUsecase
	GetBaseCurrency *walletreading.WalletGetBaseCurrencyUsecase

	AddForeign       *walletchanging.WalletAddCurrencyUsecase
	ChangeBase       *walletchanging.WalletChangeBaseCurrencyUsecase
	Create           *walletchanging.WalletCreateUsecase
	Delete           *walletchanging.WalletDeleteUsecase
	Deposit          *walletchanging.WalletDepositUsecase
	RemoveAllForeign *walletchanging.WalletRemoveAllCurrenciesUsecase
	RemoveForeign    *walletchanging.WalletRemoveCurrencyUsecase
	Withdraw         *walletchanging.WalletWithdrawUsecase
}

type ExchangeServices struct {
	ConvertCurrency *usecase.CurrencyConverterUsecase
}

func NewApplication(

	walletRepo usecase.WalletRepository,
	currencyConverter usecase.CurrencyConverter,
) *Application {
	getBalanceUC := walletreading.NewWalletGetBalanceUsecase(walletRepo)
	getForeignUC := walletreading.NewWalletGetCurrenciesUsecase(walletRepo)
	isExistsUC := walletreading.NewWalletIsExistsUsecase(walletRepo)
	getBaseCurrencyUC := walletreading.NewWalletGetBaseCurrencyUsecase(walletRepo)

	addForeignUC := walletchanging.NewWalletAddCurrencyUsecase(walletRepo)
	changeBaseUC := walletchanging.NewWalletChangeBaseCurrencyUsecase(walletRepo)
	createWalletUC := walletchanging.NewWalletCreateUsecase(walletRepo)
	deleteWalletUC := walletchanging.NewWalletDeleteUsecase(walletRepo)
	depositUC := walletchanging.NewWalletDepositUsecase(walletRepo)
	removeAllForeign := walletchanging.NewWalletRemoveAllCurrenciesUsecase(walletRepo)
	removeForeignUC := walletchanging.NewWalletRemoveCurrencyUsecase(walletRepo)
	withdrawUC := walletchanging.NewWalletWithdrawUsecase(walletRepo)

	// Этот use case может зависеть сразу от нескольких адаптеров!
	convertCurrencyUC := usecase.NewCurrencyConverterService(currencyConverter)

	return &Application{
		Wallets: WalletServices{
			GetBalance:      getBalanceUC,
			GetForeign:      getForeignUC,
			IsExists:        isExistsUC,
			GetBaseCurrency: getBaseCurrencyUC,

			AddForeign:       addForeignUC,
			ChangeBase:       changeBaseUC,
			Create:           createWalletUC,
			Delete:           deleteWalletUC,
			Deposit:          depositUC,
			RemoveAllForeign: removeAllForeign,
			RemoveForeign:    removeForeignUC,
			Withdraw:         withdrawUC,
		},
		Exchange: ExchangeServices{
			ConvertCurrency: convertCurrencyUC,
		},
	}
}
