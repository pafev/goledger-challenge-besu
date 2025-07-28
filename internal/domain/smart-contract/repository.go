package smartContractDomain

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"goledger-challenge-besu/configs/besu"
	dbConfig "goledger-challenge-besu/configs/db"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type SmartContractRepository struct {
	db            *dbConfig.DB
	abi           *abi.ABI
	boundContract *bind.BoundContract
	address       common.Address
	smartContract *SmartContractDB // instanciado aqui pois nesse escopo o contrato eh uma instancia fixa
}

func NewRepository(db *dbConfig.DB, client *besuConfig.EthClient) (*SmartContractRepository, error) {
	data, err := os.ReadFile(os.Getenv("SMART_CONTRACT_ABI_PATH"))
	if err != nil {
		return nil, err
	}
	var contractDefinition map[any]any
	json.Unmarshal(data, &contractDefinition)
	data, err = json.Marshal(contractDefinition["abi"])
	if err != nil {
		return nil, err
	}
	abi, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	contractHexAddress := os.Getenv("SMART_CONTRACT_ADDR")
	if !common.IsHexAddress(contractHexAddress) {
		return nil, errors.New("Invalid contract addresss")
	}
	var contractAddress = common.HexToAddress(contractHexAddress)

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	return &SmartContractRepository{
		db:            db,
		abi:           &abi,
		boundContract: boundContract,
		address:       contractAddress,
		smartContract: &SmartContractDB{
			SmartContractId: 1,
			Address:         contractHexAddress,
			Value:           [4]uint64{0},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}, nil
}

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
