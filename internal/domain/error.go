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
	ErrInvalidChain          = errors.New("Invalid Chain from Eth Client")
	ErrBoundContractCall     = errors.New("Error Calling Contract (BoundContract)")
	ErrBoundContractTransact = errors.New("Error Executing Transaction in Contract (BoundContract)")
	ErrInvalidSQL            = errors.New("Invalid SQL Query")
)
