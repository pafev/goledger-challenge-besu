package smartContractApp

import (
	"math/big"

	"goledger-challenge-besu/internal/domain/smart-contract"
)

type SmartContractService struct {
	repositoryDB   *smartContractDomain.SmartContractRepositoryDB
	repositoryBesu *smartContractDomain.SmartContractRepositoryBesu
}

func NewService(
	repositoryDB *smartContractDomain.SmartContractRepositoryDB,
	repositoryBesu *smartContractDomain.SmartContractRepositoryBesu) *SmartContractService {
	return &SmartContractService{repositoryDB, repositoryBesu}
}

func (r *SmartContractService) GetValue() (*big.Int, error) {
	// for multiple requests, a cache system or a service method could be implemented
	// that calls a repository of more than one type of contract
	value, err := r.repositoryBesu.GetValue()
	if err != nil {
		return new(big.Int), err
	}
	// if the repositoryDB singleton has never been updated
	if r.repositoryDB.SmartContract.UpdatedAt.IsZero() {
		r.repositoryDB.SmartContract.Value = value.Uint64()
	}
	return value, nil
}

func (r *SmartContractService) SetValue(value *big.Int, privateKey string) error {
	// here it would be possible to validate the value before passing it to the repository
	err := r.repositoryBesu.SetValue(value, privateKey)
	if err != nil {
		return err
	}
	// automatically updates the singleton of the repositoryDB
	if value.Cmp(new(big.Int).SetUint64(r.repositoryDB.SmartContract.Value)) != 0 {
		r.repositoryDB.SmartContract.Value = value.Uint64()
	}
	return nil
}

func (r *SmartContractService) CheckValue(value *big.Int) (bool, error) {
	// for multiple requests, a cache system could be implemented
	isEqual, err := r.repositoryBesu.CheckValue(value)
	if err != nil {
		return false, nil
	}
	return isEqual, nil
}

func (r *SmartContractService) SyncValue() error {
	// if the repositoryDB singleton has never been updated
	if r.repositoryDB.SmartContract.UpdatedAt.IsZero() {
		value, err := r.repositoryBesu.GetValue()
		if err != nil {
			return err
		}
		r.repositoryDB.SmartContract.Value = value.Uint64()
	}
	err := r.repositoryDB.SyncValue()
	return err
}
