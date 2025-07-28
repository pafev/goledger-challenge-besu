package smartContractDomain

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type SmartContractDB struct {
	SmartContractId uint64
	Address         string
	Value           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type SmartContractBesu struct {
	Abi           *abi.ABI
	BoundContract *bind.BoundContract
	Address       common.Address
}
