package smartContractDomain

import (
	"time"
)

type SmartContractDB struct {
	SmartContractId uint64
	Address         string
	Value           [4]uint64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
