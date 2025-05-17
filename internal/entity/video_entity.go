package entity

import "time"

type (
	Video struct {
		ID                string
		Title             string
		Overview          string
		Format            string
		TotalStock        int
		AvailableStock    int
		RentPrice         float64
		CoverPath         *string
		ProductionCompany *string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		DeletedAt         *time.Time
	}

	Videos []Video
)
