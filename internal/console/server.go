package console

import (
	"sync"

	"github.com/go-playground/validator"
	db "github.com/kodinggo/payment-service-gb1/db"
	"github.com/kodinggo/payment-service-gb1/internal/config"
	handlerHttp "github.com/kodinggo/payment-service-gb1/internal/delivery/http"
	"github.com/kodinggo/payment-service-gb1/internal/helper"
	"github.com/kodinggo/payment-service-gb1/internal/repository"
	"github.com/kodinggo/payment-service-gb1/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCMd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start the server",
	Long:  "Start the server",
	Run:   httpServer,
}

func httpServer(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	mysql := db.NewMysql()
	defer mysql.Close()

	redis := db.NewRedis()

	e := echo.New()
	e.Validator = &helper.CustomValidator{
		Validator: validator.New(),
	}

	paymentRepo := repository.NewPaymentRepository(mysql, redis)

	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo)

	// e := echo.New()

	routeGroup := e.Group("/api/v1")

	handlerHttp.NewPaymentHandler(routeGroup, paymentUsecase)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		err := e.Start(":3200")
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			logrus.Error(err.Error())
		}
	}
}
