package httpConfig

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"goledger-challenge-besu/configs/besu"
	"goledger-challenge-besu/configs/db"
	"goledger-challenge-besu/internal/app/smart-contract"
	"goledger-challenge-besu/internal/domain/smart-contract"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type HTTP struct {
	*gin.Engine
	Host           string
	Port           string
	Address        string
	AllowedOrigins string
}

func (r *HTTP) Route(ctx *context.Context, db *dbConfig.DB, ethClient *besuConfig.EthClient) error {
	// (DI) Dependency Injection and Middlewares
	smartContractRepo, err := smartContractDomain.NewRepository(ctx, db, ethClient)
	if err != nil {
		return err
	}
	smartContractService := smartContractApp.NewService(smartContractRepo)
	smartContractHandler := smartContractApp.NewHandler(smartContractService)

	// Routes
	v1 := r.Group("/api/v1")
	{
		smartContract := v1.Group("/smart-contract")
		{
			smartContract.GET("", smartContractHandler.GetValue)
			smartContract.GET("/check-value/:value", smartContractHandler.CheckValue)
			smartContract.POST("/set-value", smartContractHandler.SetValue)
		}
	}
	return nil
}

func (r *HTTP) Serve() error {
	return r.Run(r.Address)
}

func New() (*HTTP, error) {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	host := os.Getenv("APP_DOMAIN")
	port := os.Getenv("APP_PORT")
	allowedOrigins := "*"
	address := fmt.Sprintf("%s:%s", host, port)

	ginConfig := cors.DefaultConfig()
	ginConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	router := gin.New()

	// Global Middlewares
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	return &HTTP{
		router,
		host,
		port,
		address,
		allowedOrigins,
	}, nil
}
