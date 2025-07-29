package smartContractDomain

import (
	"context"
	"os"
	"time"

	"goledger-challenge-besu/configs/db"
	"goledger-challenge-besu/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
)

type SmartContractRepositoryDB struct {
	ctx           *context.Context
	db            *dbConfig.DB
	SmartContract *SmartContractDB // 'singleton' instantiated here because in this scope the contract is a fixed instance
}

// NewRepositoryDB initializes a new instance of SmartContractRepositoryDB.
// Parameters:
//   - ctx: The context for database operations.
//   - db: The database configuration to use.
//
// Returns:
//   - A pointer to SmartContractRepositoryDB if successful.
//   - An error if the provided smart contract address is invalid.
func NewRepositoryDB(ctx *context.Context, db *dbConfig.DB) (*SmartContractRepositoryDB, error) {
	contractHexAddress := os.Getenv("SMART_CONTRACT_ADDR")
	if !common.IsHexAddress(contractHexAddress) {
		return nil, domain.ErrInvalidContractAddr
	}

	return &SmartContractRepositoryDB{
		ctx: ctx,
		db:  db,
		SmartContract: &SmartContractDB{
			Address:   contractHexAddress,
			Value:     *new(uint64),
			CreatedAt: *new(time.Time),
			UpdatedAt: *new(time.Time),
		},
	}, nil
}

// SyncValue synchronizes the smart contract's data with the database.
// It checks whether the smart contract's address already exists in the database.
// If it does, the value is updated; otherwise, a new entry is created.
// Returns:
//   - An error if any database operation fails or if there is conflicting data.
func (r *SmartContractRepositoryDB) SyncValue() error {
	alreadyExists := true

	getQuery := r.db.QueryBuilder.Select("smart_contract_id").From("smart_contracts").Where(sq.Eq{"address": r.SmartContract.Address}).Limit(1)
	sql, args, err := getQuery.ToSql()
	if err != nil {
		return domain.ErrInvalidSQL
	}
	err = r.db.QueryRow(*r.ctx, sql, args...).Scan(
		&r.SmartContract.SmartContractId,
	)
	if err == pgx.ErrNoRows {
		alreadyExists = false
	} else if err != nil {
		return domain.ErrInternal
	}

	if alreadyExists {
		query := r.db.QueryBuilder.Update("smart_contracts").Set("value", r.SmartContract.Value).Where(sq.Eq{"smart_contract_id": r.SmartContract.SmartContractId}).Suffix("RETURNING *")
		sql, args, err = query.ToSql()
		if err != nil {
			return domain.ErrInvalidSQL
		}
	} else {
		query := r.db.QueryBuilder.Insert("smart_contracts").Columns("address", "value").Values(r.SmartContract.Address, r.SmartContract.Value).Suffix("RETURNING *")
		sql, args, err = query.ToSql()
		if err != nil {
			return domain.ErrInvalidSQL
		}
	}

	err = r.db.QueryRow(*r.ctx, sql, args...).Scan(
		&r.SmartContract.SmartContractId,
		&r.SmartContract.Address,
		&r.SmartContract.Value,
		&r.SmartContract.CreatedAt,
		&r.SmartContract.UpdatedAt,
	)
	if err != nil {
		if errCode := r.db.ErrorCode(err); errCode == "23505" {
			return domain.ErrConflictingData
		}
		return domain.ErrInternal
	}
	return nil
}
