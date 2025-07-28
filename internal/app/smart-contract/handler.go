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
func (r *SmartContractHandler) SetValue() error {
	return nil
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
	}
	ctx.JSON(http.StatusOK, isEqual)
}
func (r *SmartContractHandler) SyncValue() error {
	return nil
}
