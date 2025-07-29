package smartContractApp

import (
	"goledger-challenge-besu/internal/domain"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SmartContractHandler handles HTTP requests related to smart contract operations.
type SmartContractHandler struct {
	// The service layer for interacting with smart contracts.
	service *SmartContractService
}

// NewHandler initializes a new SmartContractHandler.
// Parameters:
//   - service: The SmartContractService used for business logic.
//
// Returns:
//   - A pointer to a newly created SmartContractHandler.
func NewHandler(service *SmartContractService) *SmartContractHandler {
	return &SmartContractHandler{service}
}

// GetValue retrieves the current value stored in the smart contract.
// HTTP Method: GET
// URL: /smart-contract
// Responses:
//   - 200: The current value in the smart contract.
//   - 500: Internal server error if retrieval fails.
func (r *SmartContractHandler) GetValue(ctx *gin.Context) {
	value, err := r.service.GetValue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, value)
}

type setValueRequest struct {
	Value      big.Int `json:"value" binding:"required,omitempty" example:"0"`
	PrivateKey string  `json:"privateKey" binding:"required,omitempty" example:"ef321a27ac482e12c1d1"`
}

// SetValue updates the value stored in the smart contract.
// HTTP Method: POST
// URL: /smart-contract/set-value
// Request Body:
//   - value (int): The new value to set in the contract.
//   - privateKey (string): The private key for authorization.
//
// Responses:
//   - 200: Success message upon updating the value.
//   - 400: Bad request if input validation fails.
//   - 401: Unauthorized if the private key is invalid.
//   - 500: Internal server error if the update fails.
func (r *SmartContractHandler) SetValue(ctx *gin.Context) {
	var req setValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := r.service.SetValue(&req.Value, req.PrivateKey)
	if err != nil {
		if err == domain.ErrUnauthorized {
			ctx.JSON(http.StatusUnauthorized, err.Error())
		} else {
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, "New Value Defined Successfully")
}

// CheckValue verifies if a given value matches the one in the smart contract.
// HTTP Method: GET
// URL: /smart-contract/check-value/:value
// Path Parameters:
//   - value (string): The value to compare with the contract's stored value.
//
// Responses:
//   - 200: True or false indicating if the value matches.
//   - 400: Bad request if the input value is invalid.
//   - 500: Internal server error if the verification fails.
func (r *SmartContractHandler) CheckValue(ctx *gin.Context) {
	valueStr := ctx.Param("value")
	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		ctx.JSON(http.StatusBadRequest, "Invalid param value")
	}
	isEqual, err := r.service.CheckValue(value)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, isEqual)
}

// SyncValue synchronizes the smart contract's value with an external data source.
// HTTP Method: POST
// URL: /smart-contract/sync
// Responses:
//   - 200: Success message upon synchronization.
//   - 500: Internal server error if synchronization fails.
func (r *SmartContractHandler) SyncValue(ctx *gin.Context) {
	err := r.service.SyncValue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "Sync Successfully")
}
