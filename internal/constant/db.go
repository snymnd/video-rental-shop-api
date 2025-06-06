package constant

type PaymentStatus string

const (
	PAYMENT_PENDING PaymentStatus = "pending"
	PAYMENT_EXPIRED PaymentStatus = "expired"
	PAYMENT_SUCCESS PaymentStatus = "success"
)

type PaymentMethod string

const (
	PAYMENT_METHOD_TRANSFER PaymentMethod = "transfer"
	PAYMENT_METHOD_CASH     PaymentMethod = "cash"
	PAYMENT_METHOD_DEBIT    PaymentMethod = "debit"
)

var PAYMENT_METHODS_MAP map[PaymentMethod]bool = map[PaymentMethod]bool{
	PAYMENT_METHOD_TRANSFER: true,
	PAYMENT_METHOD_CASH:     true,
	PAYMENT_METHOD_DEBIT:    true,
}

type RentalStatus string

const (
	RENTAL_PENDING  RentalStatus = "pending"
	RENTAL_FAILED   RentalStatus = "failed"
	RENTAL_RENTED   RentalStatus = "rented"
	RENTAL_RETURNED RentalStatus = "returned"
)
