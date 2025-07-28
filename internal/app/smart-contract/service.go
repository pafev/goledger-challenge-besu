package smartContractApp

type SmartContractService struct{}

func (r *SmartContractService) GetValue() ([]any, error) {
	return []any{0}, nil
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
