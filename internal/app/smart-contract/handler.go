package smartContractApp

type SmartContractHandler struct{}

func (r *SmartContractHandler) GetValue() ([]any, error) {
	return []any{0}, nil
}
func (r *SmartContractHandler) SetValue() error {
	return nil
}
func (r *SmartContractHandler) CheckValue() (bool, error) {
	return false, nil
}
func (r *SmartContractHandler) SyncValue() error {
	return nil
}
