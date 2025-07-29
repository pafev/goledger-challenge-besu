package smartContractApp

import (
	"math/big"

	"goledger-challenge-besu/internal/domain/smart-contract"
)

type SmartContractService struct {
	// here, its possible to consume other repositories or perform cache functions on top
	repository *smartContractDomain.SmartContractRepository
}

func NewService(repository *smartContractDomain.SmartContractRepository) *SmartContractService {
	return &SmartContractService{repository}
}

func (r *SmartContractService) GetValue() (*big.Int, error) {
	// for multiple requests, a cache system or a service method could be implemented
	// that calls a repository of more than one type of contract
	value, err := r.repository.GetValue()
	if err != nil {
		return new(big.Int), err
	}
	return value, nil
}
func (r *SmartContractService) SetValue(value *big.Int, privateKey string) error {
	// here it would be possible to validate the value before passing it to the repository
	err := r.repository.SetValue(value, privateKey)
	return err
}
func (r *SmartContractService) CheckValue(value *big.Int) (bool, error) {
	// For multiple requests, a cache system could be implemented
	isEqual, err := r.repository.CheckValue(value)
	if err != nil {
		return false, nil
	}
	return isEqual, nil
}
func (r *SmartContractService) SyncValue() error {
	err := r.repository.SyncValue()
	return err
}
