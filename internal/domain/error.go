package domain

import (
	"errors"
)

var (
	ErrInternal              = errors.New("Internal Error")
	ErrDataNotFound          = errors.New("Data not Found")
	ErrConflictingData       = errors.New("Data Conflicts with Existing Data in Database")
	ErrUnauthorized          = errors.New("User is Unauthorized to Access the Resource")
	ErrForbidden             = errors.New("User is Forbidden to Access the Resource")
	ErrAbiNotFound           = errors.New("Abi Contract Not Found")
	ErrInvalidAbi            = errors.New("Invalid Abi Contract")
	ErrInvalidContractAddr   = errors.New("Invalid Contract Address")
	ErrBoundContractCall     = errors.New("Error Calling Contract (BoundContract)")
	ErrBoundContractTransact = errors.New("Invalid Executing Transaction in Contract (BoundContract)")
	ErrInvalidChain          = errors.New("Error Getting Chain from Eth Client")
	ErrInvalidMined          = errors.New("Internal Error Minering")
	ErrInvalidSQL            = errors.New("Internal Error Invalid SQL")
)
