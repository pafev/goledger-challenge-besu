package smartContractDomain

import (
	"time"

	"github.com/core-coin/uint256"
)

type SmartContractDB struct {
	SmartContractId uint64
	Address         string
	Value           uint256.Int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
