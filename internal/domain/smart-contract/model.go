package smartContractDomain

import (
	"math/big"
	"time"
)

type SmartContractDB struct {
	SmartContractId uint64
	Address         string
	Value           big.Int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
