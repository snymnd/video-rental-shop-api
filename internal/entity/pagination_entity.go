package entity

type (
	PageFilter struct {
		Field string
		Value any
	}

	PageInfo struct {
		Page      int
		Limit     int
		OrderBy   []string
		OrderSort string
		Filters   []PageFilter
		TotalRow  int
	}

	PaginationQuery struct {
		Limit int
		Page  int
	}

	PaginatadResponse[T any] struct {
		Entries []T
		PageInfo
	}
)
