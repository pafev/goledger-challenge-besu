package smartContractDomain

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
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
	abiPath := os.Getenv("SMART_CONTRACT_ABI_PATH")
	data, err := os.ReadFile(abiPath)
	if err != nil {
		slog.Error("Error reading contract abi file", "path", abiPath, "error", err.Error())
		return nil, err
	}
	var contractDefinition struct {
		Abi any `json:"abi"`
	}
	json.Unmarshal(data, &contractDefinition)
	data, err = json.Marshal(contractDefinition.Abi)
	if err != nil {
		slog.Error("Error getting contract abi structure from abi file", "error", err.Error())
		return nil, err
	}
	abi, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		slog.Error("Error converting contract abi json to string", "error", err.Error())
		return nil, err
	}

	contractHexAddress := os.Getenv("SMART_CONTRACT_ADDR")
	if !common.IsHexAddress(contractHexAddress) {
		slog.Error("Error reading contract addres", "error", "Contract address in invalid format")
		return nil, errors.New("Invalid contract address")
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
		slog.Error("Error calling contract (bound contract)", "options", caller, "error", err.Error())
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
		slog.Error("Error getting chain from eth client", "error", err.Error())
		return domain.ErrInvalidChain
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		slog.Error("Error converting private key hex format to ECDSA format", "error", err.Error())
		return domain.ErrUnauthorized
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainId)
	if err != nil {
		slog.Error("Error getting auth opts to transact bound contract", "error", err.Error())
		return domain.ErrUnauthorized
	}

	tx, err := r.boundContract.Transact(auth, "set", value)
	if err != nil {
		slog.Error("Error executing transaction in contract (bound contract)", "options", auth, "error", err.Error())
		return domain.ErrBoundContractTransact
	}

	tx.Hash().Hex()

	if _, err = bind.WaitMined(*r.ctx, r.client, tx); err != nil {
		slog.Error("Error waiting to be mined", "error", err.Error())
		return domain.ErrBoundContractTransact
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
		slog.Error("Error contract value in SmartContractRepositoryBesu.CheckValue", "error", err.Error())
		return false, err
	}
	return correctValue.Cmp(value) == 0, nil
}
