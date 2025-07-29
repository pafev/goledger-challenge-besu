package smartContractDomain

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"time"

	"goledger-challenge-besu/configs/besu"
	dbConfig "goledger-challenge-besu/configs/db"
	"goledger-challenge-besu/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
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

	return &SmartContractRepository{
		db:            db,
		ctx:           ctx,
		abi:           &abi,
		boundContract: boundContract,
		address:       contractAddress,
		client:        client,
		smartContract: &SmartContractDB{
			Address:   contractHexAddress,
			Value:     *new(uint64),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
		return new(big.Int), domain.ErrBoundContractCall
	}
	result := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	return result, nil
}

func (r *SmartContractRepository) SetValue(value *big.Int, privateKey string) error {
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

	// updatedAt sera atualizado mediante trigger no db com o metodo de sync
	r.smartContract.Value = value.Uint64()

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
	alreadyExists := true

	getQuery := r.db.QueryBuilder.Select("smart_contract_id").From("smart_contracts").Where(sq.Eq{"address": r.smartContract.Address}).Limit(1)
	sql, args, err := getQuery.ToSql()
	if err != nil {
		return domain.ErrInvalidSQL
	}
	err = r.db.QueryRow(*r.ctx, sql, args...).Scan(
		&r.smartContract.SmartContractId,
	)
	if err == pgx.ErrNoRows {
		alreadyExists = false
	} else if err != nil {
		return domain.ErrInternal
	}

	if alreadyExists {
		query := r.db.QueryBuilder.Update("smart_contracts").Set("value", r.smartContract.Value).Where(sq.Eq{"smart_contract_id": r.smartContract.SmartContractId}).Suffix("RETURNING *")
		sql, args, err = query.ToSql()
		if err != nil {
			return domain.ErrInvalidSQL
		}
	} else {
		query := r.db.QueryBuilder.Insert("smart_contracts").Columns("address", "value").Values(r.smartContract.Address, r.smartContract.Value).Suffix("RETURNING *")
		sql, args, err = query.ToSql()
		if err != nil {
			return domain.ErrInvalidSQL
		}
	}

	err = r.db.QueryRow(*r.ctx, sql, args...).Scan(
		&r.smartContract.SmartContractId,
		&r.smartContract.Address,
		&r.smartContract.Value,
		&r.smartContract.CreatedAt,
		&r.smartContract.UpdatedAt,
	)
	if err != nil {
		if errCode := r.db.ErrorCode(err); errCode == "23505" {
			return domain.ErrConflictingData
		}
		return domain.ErrInternal
	}
	return nil
}
