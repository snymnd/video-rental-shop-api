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
		GenreIDs          []int
		RentPrice         float64
		CoverPath         *string
		ProductionCompany *string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		DeletedAt         *time.Time
	}

	GetVideosParams struct {
		SortOrder string
		GenreIDs  []int
		Title     string
		OrderBy   []string
		PaginationQuery
	}

	Videos []Video

	GetVideosReturn struct {
		PageInfo PageInfo
		Entries  Videos
	}
)
