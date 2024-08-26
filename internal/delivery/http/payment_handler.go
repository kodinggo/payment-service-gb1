package handler

import (
	"net/http"
	"strconv"

	"github.com/tubagusmf/payment-service-gb1/internal/model"

	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type paymentHandler struct {
	PaymentUsecase model.IPaymentUsecase
}

func NewPaymentHandler(e *echo.Group, us model.IPaymentUsecase) {
	handlers := &paymentHandler{
		PaymentUsecase: us,
	}

	payments := e.Group("/payments")

	// payments.Use(echojwt.WithConfig(jwtConfig()))

	payments.GET("", handlers.GetPayments)
	payments.GET("/:id", handlers.GetPayment)
	payments.POST("", handlers.CreatePayment)
	payments.PUT("/:id", handlers.UpdatePayment)
	payments.DELETE("/:id", handlers.DeletePayment)
}

func (p *paymentHandler) GetPayments(c echo.Context) error {
	// claims := claimsSession(c)
	// if claims == nil {
	// 	return c.JSON(http.StatusUnauthorized, response{
	// 		Status:  http.StatusForbidden,
	// 		Message: "Forbidden",
	// 	})
	// }

	// if claims.Role != "admin" {
	// 	return c.JSON(http.StatusForbidden, response{
	// 		Status:  http.StatusForbidden,
	// 		Message: "Forbidden",
	// 	})
	// }

	reqLimit := c.QueryParam("limit")
	reqOffset := c.QueryParam("offset")

	var limit, offset int32
	if reqLimit == "" {
		limit = 10
	}
	if reqOffset == "" {
		offset = 0
	}

	payments, err := p.PaymentUsecase.FindAll(c.Request().Context(), model.PaymentFilter{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    payments,
	})
}

func (p *paymentHandler) GetPayment(c echo.Context) error {
	// claims := claimsSession(c)
	// if claims == nil {
	// 	return c.JSON(http.StatusUnauthorized, response{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "Unauthorized",
	// 	})
	// }

	id := c.Param("id")
	parseId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	payment, err := p.PaymentUsecase.FindById(c.Request().Context(), int64(parseId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    payment,
	})
}

func (p *paymentHandler) CreatePayment(c echo.Context) error {
	// claims := claimsSession(c)
	// if claims == nil {
	// 	return c.JSON(http.StatusUnauthorized, response{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "Unauthorized",
	// 	})
	// }

	var in model.CreatePaymentInput

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := p.PaymentUsecase.Create(c.Request().Context(), in); err != nil {
		return c.JSON(http.StatusInternalServerError, response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response{
		Status:  http.StatusCreated,
		Message: "success",
	})
}

func (p *paymentHandler) UpdatePayment(c echo.Context) error {
	// claims := claimsSession(c)
	// if claims == nil {
	// 	return c.JSON(http.StatusUnauthorized, response{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "Unauthorized",
	// 	})
	// }

	paymentId := c.Param("id")
	parseId, err := strconv.Atoi(paymentId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var in model.UpdatePaymentInput

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if parseId > 0 {
		in.Id = int64(parseId)
	}

	if err := p.PaymentUsecase.Update(c.Request().Context(), in); err != nil {
		return c.JSON(http.StatusInternalServerError, response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response{
		Status:  http.StatusOK,
		Message: "success",
	})
}

func (p *paymentHandler) DeletePayment(c echo.Context) error {
	// claims := claimsSession(c)
	// if claims == nil {
	// 	return c.JSON(http.StatusUnauthorized, response{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "Unauthorized",
	// 	})
	// }

	id := c.Param("id")
	parseId, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := p.PaymentUsecase.Delete(c.Request().Context(), int64(parseId)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response{
		Status:  http.StatusOK,
		Message: "success",
	})
}
