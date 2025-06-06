package usecase

import (
	"context"
	"fmt"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
)

type (
	PaymentRepository interface {
		UpdatePayment(ctx context.Context, paymentParams entity.UpdatePaymentParams) error
		GetPayment(ctx context.Context, paymentID int) (entity.Payment, error)
	}

	PaymentRentalRepository interface {
		UpdatesRentalStatusByPaymentID(ctx context.Context, paymentID int, status constant.RentalStatus) error
	}

	PaymentTxRepository interface {
		WithTx(ctx context.Context, tFunc func(ctx context.Context) error) error
	}

	PaymentUsecase struct {
		pr  PaymentRepository
		rr  PaymentRentalRepository
		txr PaymentTxRepository
	}
)

func NewPaymentUsecase(pr PaymentRepository, rr PaymentRentalRepository, txr PaymentTxRepository) *PaymentUsecase {
	return &PaymentUsecase{pr, rr, txr}
}

func (puc PaymentUsecase) PayRentals(ctx context.Context, paymentID int, paymentMethod constant.PaymentMethod) error {
	if err := puc.txr.WithTx(ctx, func(txCtx context.Context) error {
		// get payment
		payment, getPaymentErr := puc.pr.GetPayment(txCtx, paymentID)
		if getPaymentErr != nil {
			return getPaymentErr
		}

		// validate payment status
		if payment.Status == constant.PAYMENT_EXPIRED {
			return customerrors.NewError(
				"payment is already expired",
				fmt.Errorf("payment with id %d is already expired", paymentID),
				customerrors.InvalidAction,
			)
		}
		if payment.Status != constant.PAYMENT_PENDING {
			return customerrors.NewError(
				"payment is not found",
				fmt.Errorf("payment id: %d status is %s not %s", paymentID, payment.Status, constant.PAYMENT_PENDING),
				customerrors.InvalidAction,
			)
		}

		succesStatus := constant.PAYMENT_SUCCESS
		updatePaymentParams := entity.UpdatePaymentParams{
			ID:     paymentID,
			Method: &paymentMethod,
			Status: &succesStatus,
		}

		// update status of a payment to success
		if err := puc.pr.UpdatePayment(txCtx, updatePaymentParams); err != nil {
			return err
		}

		// updates every rental status to rented with specified payment id
		if err := puc.rr.UpdatesRentalStatusByPaymentID(txCtx, paymentID, constant.RENTAL_RENTED); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
