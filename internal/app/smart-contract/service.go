package smartContractApp

import smartContractDomain "goledger-challenge-besu/internal/domain/smart-contract"

type SmartContractService struct {
	// aqui eh possivel consumir outros repositories ou executar funcoes de cache por cima
	repository *smartContractDomain.SmartContractRepository
}

func NewService(repository *smartContractDomain.SmartContractRepository) *SmartContractService {
	return &SmartContractService{repository}
}

func (r *SmartContractService) GetValue() ([]any, error) {
	// para multiplas requisicoes, poderia ser implementado um sistema de cache
	value, err := r.repository.GetValue()
	if err != nil {
		return nil, err
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
