package smartContractDomain

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"
	"time"

	"goledger-challenge-besu/configs/besu"
	dbConfig "goledger-challenge-besu/configs/db"
	"goledger-challenge-besu/internal/domain"

	"github.com/core-coin/uint256"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type SmartContractRepository struct {
	db            *dbConfig.DB
	ctx           *context.Context
	abi           *abi.ABI
	boundContract *bind.BoundContract
	address       common.Address
	client        *besuConfig.EthClient
	smartContract *SmartContractDB // instanciado aqui pois nesse escopo o contrato eh uma instancia fixa
}

func NewRepository(ctx *context.Context, db *dbConfig.DB, client *besuConfig.EthClient) (*SmartContractRepository, error) {
	data, err := os.ReadFile(os.Getenv("SMART_CONTRACT_ABI_PATH"))
	if err != nil {
		return nil, err
	}
	var contractDefinition map[string]any
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
		ctx:           ctx,
		abi:           &abi,
		boundContract: boundContract,
		address:       contractAddress,
		client:        client,
		smartContract: &SmartContractDB{
			SmartContractId: 1,
			Address:         contractHexAddress,
			Value:           uint256.Int{0},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}, nil
}

func (r *SmartContractRepository) GetValue() (*big.Int, error) {
	caller := bind.CallOpts{
		Pending: false,
		Context: *r.ctx,
	}
	var output []any
	err := r.boundContract.Call(&caller, &output, "get")
	if err != nil {
		return new(big.Int), domain.ErrInternal
	}
	result := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	return result, nil
}
func (r *SmartContractRepository) SetValue(value *big.Int, privateKey string) error {
	chainId, err := r.client.ChainID(*r.ctx)
	if err != nil {
		return domain.ErrInternal
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return domain.ErrInternal
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainId)
	if err != nil {
		return domain.ErrInternal
	}

	tx, err := r.boundContract.Transact(auth, "set", value)
	if err != nil {
		return domain.ErrInternal
	}

	tx.Hash().Hex()

	_, err = bind.WaitMined(
		context.Background(),
		r.client,
		tx,
	)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}
func (r *SmartContractRepository) CheckValue(value *big.Int) (bool, error) {
	correctValue, err := r.GetValue()
	if err != nil {
		return false, err
	}
	return correctValue.Cmp(value) == 0, nil
}
func (r *SmartContractRepository) SyncValue() error {
	return nil
}
