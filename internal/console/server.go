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
	authPb "github.com/kodinggo/user-service-gb1/pb/auth"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	authConn, err := grpc.NewClient("localhost:4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal(err)
	}
	auth := authPb.NewJWTValidatorClient(authConn)
	authMiddleware := helper.NewJWTMiddleware(auth)

	handlerHttp.NewPaymentHandler(routeGroup, paymentUsecase, authMiddleware.ValidateJWT)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		err := e.Start(":3400")
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
