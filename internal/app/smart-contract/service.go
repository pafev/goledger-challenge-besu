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
func (r *SmartContractService) SetValue(value *big.Int, privateKey string) error {
	// seria possivel colocar validacoes sobre o valor, antes de passar ao repository
	err := r.repository.SetValue(value, privateKey)
	return err
}
func (r *SmartContractService) CheckValue(value *big.Int) (bool, error) {
	// para multiplas requisicoes, poderia ser implementado um sistema de cache
	isEqual, err := r.repository.CheckValue(value)
	if err != nil {
		return false, nil
	}
	return isEqual, nil
}
func (r *SmartContractService) SyncValue() error {
	return nil
}
