package models

const (
	PAYMENT_METHODS_KEY = "payment_methods"
	CASH                = "cash"
	CC                  = "stripe_cc"

	ORDER_METHODS_KEY = "order_methods"
	DELIVERY          = "delivery"
	PICKUP            = "pickup"

	ALL_STATUSES_KEY      = "status"
	DELIVERY_STATUSES_KEY = "delivery_status"
	PICKUP_STATUSES_KEY   = "pickup_status"

	PENDING     = "PENDING"
	APPROVED    = "APPROVED"
	REJECTED    = "REJECTED"
	IN_PROGRESS = "IN-PROGRESS"
	CANCELED    = "CANCELED"
	EN_ROUTE    = "EN-ROUTE"
	COMPLETED   = "COMPLETED"

	FAILED_PAYMENT_HOLD    = "FAILED-PAYMENT-HOLD"
	PAYMENT_ON_HOLD        = "PAYMENT-ON-HOLD"
	FAILED_PAYMENT_CAPTURE = "FAILED-PAYMENT-CAPTURE"
	PAYMENT_CAPTURED       = "PAYMENT-CAPTURED"
)

var (
	StripeSK = ""

	OrderMethod    = []string{DELIVERY, PICKUP}
	PaymentMethods = []string{CASH, CC}
	OrderStatuses  = []string{
		CANCELED,
		IN_PROGRESS,
		COMPLETED,
		REJECTED,
		EN_ROUTE,
		APPROVED,
		PENDING,
	}
	DeliveryOrderStatuses = OrderStatuses
	PickupOrderStatuses   = []string{
		CANCELED,
		IN_PROGRESS,
		COMPLETED,
		REJECTED,
		APPROVED,
		PENDING,
	}

	STATUS_TREE = map[string]interface{}{
		DELIVERY: map[string]interface{}{
			PENDING: map[string]interface{}{
				APPROVED: map[string]interface{}{
					CANCELED: []string{},
					IN_PROGRESS: map[string][]string{
						EN_ROUTE: []string{
							CANCELED,
							COMPLETED,
						},
						CANCELED: []string{},
					},
				},
				REJECTED: []string{},
			},
		},
		PICKUP: map[string]interface{}{
			PENDING: map[string]interface{}{
				APPROVED: map[string][]string{
					CANCELED: []string{},
					IN_PROGRESS: []string{
						CANCELED,
						COMPLETED,
					},
				},
				REJECTED: []string{},
			},
		},
	}

	ALLOWED_STATUS_PATH = map[string]map[string][]string{
		DELIVERY: map[string][]string{
			CANCELED:  []string{},
			COMPLETED: []string{},
			REJECTED:  []string{},
			IN_PROGRESS: []string{
				EN_ROUTE,
				CANCELED,
			},
			EN_ROUTE: []string{
				COMPLETED,
				CANCELED,
			},
			APPROVED: []string{
				IN_PROGRESS,
				CANCELED,
			},
			PENDING: []string{
				APPROVED,
				REJECTED,
			},
		},
		PICKUP: map[string][]string{
			CANCELED:  []string{},
			COMPLETED: []string{},
			REJECTED:  []string{},
			IN_PROGRESS: []string{
				COMPLETED,
				CANCELED,
			},
			APPROVED: []string{
				IN_PROGRESS,
				CANCELED,
			},
			PENDING: []string{
				APPROVED,
				REJECTED,
			},
		},
	}

	// helper map for all constants to allow for kw lookup
	CONST_MAP = map[string]map[string]string{
		PAYMENT_METHODS_KEY: map[string]string{
			CASH: CASH,
			CC:   CC,
		},
		ORDER_METHODS_KEY: map[string]string{
			DELIVERY: DELIVERY,
			PICKUP:   PICKUP,
		},
		ALL_STATUSES_KEY: map[string]string{
			CANCELED:    CANCELED,
			IN_PROGRESS: IN_PROGRESS,
			COMPLETED:   COMPLETED,
			REJECTED:    REJECTED,
			EN_ROUTE:    EN_ROUTE,
			APPROVED:    APPROVED,
			PENDING:     PENDING,
		},
		DELIVERY_STATUSES_KEY: map[string]string{
			CANCELED:    CANCELED,
			IN_PROGRESS: IN_PROGRESS,
			COMPLETED:   COMPLETED,
			REJECTED:    REJECTED,
			EN_ROUTE:    EN_ROUTE,
			APPROVED:    APPROVED,
			PENDING:     PENDING,
		},
		PICKUP_STATUSES_KEY: map[string]string{
			CANCELED:    CANCELED,
			IN_PROGRESS: IN_PROGRESS,
			COMPLETED:   COMPLETED,
			REJECTED:    REJECTED,
			APPROVED:    APPROVED,
			PENDING:     PENDING,
		},
	}

	OrderStatusMsgs = map[string]string{
		CANCELED:    "We regret to inform you that your order has been canceled.",
		IN_PROGRESS: "Your order is in the works!",
		COMPLETED:   "Your order has been completed! Thank you for shopping with us and we hope to see you again soon!",
		REJECTED:    "We regret to inform you that your order could not be accepted at this time.",
		EN_ROUTE:    "Your order is on it's way!",
		APPROVED:    "Thank you for shopping with us! We would be glad to fulfill your order.",
		PENDING:     "Your order is currently being reviwed by the store.",
	}
)
