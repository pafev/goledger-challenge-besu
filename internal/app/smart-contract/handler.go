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
	ctx.JSON(http.StatusOK, value.String())
}

type setValueRequest struct {
	Value string `json:"value" binding:"required,omitempty" example:"0"`
}

func (r *SmartContractHandler) SetValue(ctx *gin.Context) {
	var req setValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		ctx.JSON(http.StatusBadRequest, "Invalid value param")
	}
	err := r.service.SetValue(value)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, value.String())
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
func (r *SmartContractHandler) SyncValue() error {
	return nil
}
