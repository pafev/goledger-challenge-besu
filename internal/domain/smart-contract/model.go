package smartContractDomain

import (
	"time"
)

type SmartContractDB struct {
	SmartContractId uint64
	Address         string
	Value           uint64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
