package smartContractApp

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SmartContractHandler struct {
	// aqui sao feitos os tratamentos http por cima do service
	service *SmartContractService
}

func NewHandler(service *SmartContractService) *SmartContractHandler {
	return &SmartContractHandler{service}
}

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

func (r *SmartContractHandler) SetValue(ctx *gin.Context) {
	var req setValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := r.service.SetValue(&req.Value, req.PrivateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "New Value Defined Succesfully")
}

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

func (r *SmartContractHandler) SyncValue(ctx *gin.Context) {
	err := r.service.SyncValue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "Sync Succesfully")
}
