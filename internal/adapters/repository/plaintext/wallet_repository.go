package plaintext

import (
	"encoding/json"
	"os"
	"path/filepath"
	"travelWallet/internal/domain"
	"travelWallet/internal/usecase"
)

var _ usecase.WalletRepository = (*FileWalletRepository)(nil)

type FileWalletRepository struct {
	folderPath string
}

func NewWalletRepository() (*FileWalletRepository, error) {
	walletsDir := "wallets"

	err := os.Mkdir(walletsDir, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &FileWalletRepository{
		folderPath: filepath.Join(currentDir, walletsDir) + "/",
	}, nil
}

// CreateByID implements usecase.WalletRepository.
func (f *FileWalletRepository) CreateByID(id string, baseCurrency string) error {
	file, err := os.Create(f.folderPath + id + ".txt")
	if err != nil {
		return err
	}

	wallet, err := domain.NewWallet(id, baseCurrency)
	if err != nil {
		return err
	}

	jsonWallet, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonWallet)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID implements usecase.WalletRepository.
func (f *FileWalletRepository) DeleteByID(id string) error {
	panic("unimplemented")
}

// FindByID implements usecase.WalletRepository.
func (f *FileWalletRepository) FindByID(id string) (*domain.Wallet, error) {
	fileData, err := os.ReadFile(f.folderPath + id + ".txt")
	if err != nil {
		return nil, err
	}

	wallet := domain.Wallet{}

	err = json.Unmarshal(fileData, &wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// Save implements usecase.WalletRepository.
func (f *FileWalletRepository) Save(wallet *domain.Wallet) error {
	file, err := os.Create(f.folderPath + wallet.GetID() + ".txt")
	if err != nil {
		return err
	}

	jsonWallet, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonWallet)
	if err != nil {
		return err
	}

	return nil
}
