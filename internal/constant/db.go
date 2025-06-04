package constant

type PaymentStatus string

const (
	PAYMENT_PENDING PaymentStatus = "pending"
	PAYMENT_EXPIRED PaymentStatus = "expired"
	PAYMENT_SUCCESS PaymentStatus = "success"
)

type RentalStatus string

const (
	RENTAL_PENDING  RentalStatus = "pending"
	RENTAL_FAILED   RentalStatus = "failed"
	RENTAL_RENTED   RentalStatus = "rented"
	RENTAL_RETURNED RentalStatus = "returned"
)
