package smartContractApp

import (
	"math/big"

	"goledger-challenge-besu/internal/domain/smart-contract"
)

type SmartContractService struct {
	// aqui eh possivel consumir outros repositories ou executar funcoes de cache por cima
	repository *smartContractDomain.SmartContractRepository
}

func NewService(repository *smartContractDomain.SmartContractRepository) *SmartContractService {
	return &SmartContractService{repository}
}

func (r *SmartContractService) GetValue() (*big.Int, error) {
	// para multiplas requisicoes, poderia ser implementado um sistema de cache
	// ou um metodo de service que chamasse repository de mais de um tipo de contrato
	value, err := r.repository.GetValue()
	if err != nil {
		return new(big.Int), err
	}
	return value, nil
}
func (r *SmartContractService) SetValue() error {
	return nil
}
func (r *SmartContractService) CheckValue() (bool, error) {
	return false, nil
}
func (r *SmartContractService) SyncValue() error {
	return nil
}
