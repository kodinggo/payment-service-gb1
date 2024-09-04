package handler

import (
	"net/http"
	"strconv"

	"github.com/kodinggo/payment-service-gb1/internal/model"
	"github.com/kodinggo/payment-service-gb1/internal/usecase"
	"github.com/kodinggo/payment-service-gb1/internal/utils"
	"github.com/kodinggo/user-service-gb1/middleware-example/helper"

	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type paymentHandler struct {
	PaymentUsecase model.IPaymentUsecase
}

func NewPaymentHandler(e *echo.Group, us model.IPaymentUsecase, auth echo.MiddlewareFunc) {
	handlers := &paymentHandler{
		PaymentUsecase: us,
	}

	payments := e.Group("/payments")

	payments.Use(auth)

	payments.GET("", handlers.GetPayments)
	payments.GET("/:id", handlers.GetPayment)
	payments.POST("", handlers.CreatePayment)
	payments.PUT("/:id", handlers.UpdatePayment)
	payments.DELETE("/:id", handlers.DeletePayment)
}

func (p *paymentHandler) GetPayments(c echo.Context) error {

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

		if err == usecase.ErrNotFound {
			return c.JSON(http.StatusNotFound, response{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

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

	var in model.CreatePaymentInput

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err := c.Validate(&in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user := helper.GetUserSession(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, response{
			Status:  http.StatusBadRequest,
			Message: "Error getting user : unauthorized",
		})
	}

	if user.Role.Name != "admin" {
		return c.JSON(http.StatusForbidden, response{
			Status:  http.StatusForbidden,
			Message: "sorry you dont have permission",
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

	err = c.Validate(&in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user, err := utils.UserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, response{
			Status:  http.StatusForbidden,
			Message: "sorry you dont have permission",
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

	id := c.Param("id")
	parseId, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := utils.UserClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, response{
			Status:  http.StatusForbidden,
			Message: "sorry you dont have permission",
		})
	}

	if err := p.PaymentUsecase.Delete(c.Request().Context(), int64(parseId)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response{
		Status:  http.StatusOK,
		Message: "success",
	})
}
