package constant

type PaymentStatus string

const (
	PAYMENT_PENDING PaymentStatus = "pending"
	PAYMENT_EXPIRED PaymentStatus = "expired"
	PAYMENT_SUCCESS PaymentStatus = "success"
)
