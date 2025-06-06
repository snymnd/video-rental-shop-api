package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"

	"github.com/gin-gonic/gin"
)

type (
	PaymentUsecase interface {
		PayRentals(ctx context.Context, paymentID int, paymentMethod constant.PaymentMethod) error
	}

	PaymentController struct {
		puc PaymentUsecase
	}
)

func NewPaymentController(puc PaymentUsecase) *PaymentController {
	return &PaymentController{puc}
}

func (pc *PaymentController) PayRentals(ctx *gin.Context) {
	paymentMethodParam := strings.ToLower(ctx.Param("method"))

	if _, isExist := constant.PAYMENT_METHODS_MAP[constant.PaymentMethod(paymentMethodParam)]; !isExist {
		ctx.Error(customerrors.NewError(
			"invalid payment method",
			fmt.Errorf("method: '%s' is not a valid payment method", paymentMethodParam),
			customerrors.InvalidAction,
		))
		return
	}
	paymentMethod := constant.PaymentMethod(paymentMethodParam)

	paymentIdParam := ctx.Param("id")
	paymentId, err := strconv.Atoi(paymentIdParam)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"payment id must be a number",
			err,
			customerrors.InvalidAction,
		))
		return
	}

	if err := pc.puc.PayRentals(ctx, paymentId, paymentMethod); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    fmt.Sprintf("Payment id: %d is successfully paid", paymentId),
	})
}
