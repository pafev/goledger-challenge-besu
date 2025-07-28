package smartContractDomain

type SmartContractRepository struct{}

func (r *SmartContractRepository) GetValue() ([]any, error) {
	return []any{0}, nil
}
func (r *SmartContractRepository) SetValue() error {
	return nil
}
func (r *SmartContractRepository) CheckValue() (bool, error) {
	return false, nil
}
func (r *SmartContractRepository) SyncValue() error {
	return nil
}
