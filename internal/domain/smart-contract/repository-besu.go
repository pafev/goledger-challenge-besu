package smartContractDomain

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strings"

	"goledger-challenge-besu/configs/besu"
	"goledger-challenge-besu/internal/domain"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type SmartContractRepositoryBesu struct {
	ctx           *context.Context
	abi           *abi.ABI
	boundContract *bind.BoundContract
	address       common.Address
	client        *besuConfig.EthClient
}

// NewRepositoryBesu initializes a new instance of SmartContractRepositoryBesu.
// Parameters:
//   - ctx: The context for contract operations.
//   - client: The Ethereum client configuration.
//
// Returns:
//   - A pointer to SmartContractRepositoryBesu if successful.
//   - An error if there is an issue with the ABI or contract address.
func NewRepositoryBesu(ctx *context.Context, client *besuConfig.EthClient) (*SmartContractRepositoryBesu, error) {
	data, err := os.ReadFile(os.Getenv("SMART_CONTRACT_ABI_PATH"))
	if err != nil {
		return nil, domain.ErrAbiNotFound
	}
	var contractDefinition map[string]any
	json.Unmarshal(data, &contractDefinition)
	data, err = json.Marshal(contractDefinition["abi"])
	if err != nil {
		return nil, domain.ErrInvalidAbi
	}
	abi, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		return nil, domain.ErrInvalidAbi
	}

	contractHexAddress := os.Getenv("SMART_CONTRACT_ADDR")
	if !common.IsHexAddress(contractHexAddress) {
		return nil, domain.ErrInvalidContractAddr
	}
	var contractAddress = common.HexToAddress(contractHexAddress)

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	return &SmartContractRepositoryBesu{
		ctx:           ctx,
		abi:           &abi,
		boundContract: boundContract,
		address:       contractAddress,
		client:        client,
	}, nil
}

// GetValue retrieves the current value stored in the smart contract.
// Returns:
//   - A pointer to a big.Int containing the value.
//   - An error if the call to the bound contract fails.
func (r *SmartContractRepositoryBesu) GetValue() (*big.Int, error) {
	caller := bind.CallOpts{
		Pending: false,
		Context: *r.ctx,
	}
	var output []any
	err := r.boundContract.Call(&caller, &output, "get")
	if err != nil {
		return new(big.Int), domain.ErrBoundContractCall
	}
	result := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	return result, nil
}

// SetValue sets a new value in the smart contract.
// Parameters:
//   - value: A pointer to a big.Int containing the value to set.
//   - privateKey: A string representing the private key for transaction authorization.
//
// Returns:
//   - An error if the chain ID retrieval, private key parsing, or transaction execution fails.
func (r *SmartContractRepositoryBesu) SetValue(value *big.Int, privateKey string) error {
	chainId, err := r.client.ChainID(*r.ctx)
	if err != nil {
		return domain.ErrInvalidChain
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return domain.ErrUnauthorized
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainId)
	if err != nil {
		return domain.ErrUnauthorized
	}

	tx, err := r.boundContract.Transact(auth, "set", value)
	if err != nil {
		return domain.ErrBoundContractTransact
	}

	tx.Hash().Hex()

	_, err = bind.WaitMined(
		context.Background(),
		r.client,
		tx,
	)
	if err != nil {
		return domain.ErrInvalidMined
	}

	return nil
}

// CheckValue verifies if the given value matches the value stored in the smart contract.
// Parameters:
//   - value: A pointer to a big.Int containing the value to check.
//
// Returns:
//   - A boolean indicating if the values match.
//   - An error if retrieving the current value from the smart contract fails.
func (r *SmartContractRepositoryBesu) CheckValue(value *big.Int) (bool, error) {
	correctValue, err := r.GetValue()
	if err != nil {
		return false, err
	}
	return correctValue.Cmp(value) == 0, nil
}
