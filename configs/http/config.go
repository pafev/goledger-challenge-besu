package httpConfig

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"goledger-challenge-besu/configs/besu"
	"goledger-challenge-besu/configs/db"
	"goledger-challenge-besu/internal/app/smart-contract"
	"goledger-challenge-besu/internal/domain/smart-contract"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
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
	// (DI) Dependency Injection
	smartContractRepoBesu, err := smartContractDomain.NewRepositoryBesu(ctx, ethClient)
	if err != nil {
		slog.Error("Error building SmartContractRepositoryBesu", "error", err)
		return err
	}
	smartContractRepoDB, err := smartContractDomain.NewRepositoryDB(ctx, db)
	if err != nil {
		slog.Error("Error building SmartContractRepositoryDB", "error", err)
		return err
	}
	smartContractService := smartContractApp.NewService(smartContractRepoDB, smartContractRepoBesu)
	smartContractHandler := smartContractApp.NewHandler(smartContractService)

	// Routes and Middlewares (for specifics groups or routes)
	v1 := r.Group("/api/v1")
	{
		smartContract := v1.Group("/smart-contract")
		{
			smartContract.GET("", smartContractHandler.GetValue)
			smartContract.GET("/check-value/:value", smartContractHandler.CheckValue)
			smartContract.POST("/set-value", smartContractHandler.SetValue)
			smartContract.POST("/sync", smartContractHandler.SyncValue)
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
	router.Use(timeout.New(
		timeout.WithTimeout(12*time.Second),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, "Request Timed Out")
		}),
	))
	// ...it would be possible, for example, to add middleware to strip slashes

	return &HTTP{
		router,
		host,
		port,
		address,
		allowedOrigins,
	}, nil
}
