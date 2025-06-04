package dto

import "time"

type (
	CreateVideoReq struct {
		Title             string  `json:"title" binding:"required"`
		Overview          string  `json:"overview,omitempty"`
		RentPrice         float64 `json:"rent_price" binding:"required,gte=0"`
		Format            string  `json:"format" binding:"eq=dvd|eq=bluray|eq=digital|eq=vhs"`
		TotalStock        int     `json:"total_stock" binding:"required,gte=0"`
		GenreIDs          []int   `json:"genre_ids"`
		CoverPath         *string `json:"cover_path,omitempty"`
		ProductionCompany *string `json:"production_company,omitempty"`
	}

	CreateVideoRes struct {
		ID                int     `json:"id" `
		Title             string  `json:"title"`
		Overview          string  `json:"overview,omitempty"`
		Format            string  `json:"format"`
		GenreIDs          []int   `json:"genre_ids"`
		RentPrice         float64 `json:"rent_price"`
		TotalStock        int     `json:"total_stock"`
		AvailableStock    int     `json:"available_stock"`
		CoverPath         *string `json:"cover_path,omitempty"`
		ProductionCompany *string `json:"production_company,omitempty"`
	}

	VideoRes struct {
		ID                int        `json:"id" `
		Title             string     `json:"title"`
		Overview          string     `json:"overview,omitempty"`
		Format            string     `json:"format"`
		RentPrice         float64    `json:"rent_price"`
		GenreIDs          []int      `json:"genre_ids"`
		TotalStock        int        `json:"total_stock"`
		AvailableStock    int        `json:"available_stock"`
		CoverPath         *string    `json:"cover_path,omitempty"`
		ProductionCompany *string    `json:"production_company,omitempty"`
		CreatedAt         *time.Time `json:"created_at,omitempty"`
		UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	}

	GetVideosQuery struct {
		OrderSort string `form:"order_sort,omitempty"`
		// TODO: validate number on genre ids
		Genres  []int    `form:"[]genres,omitempty" binding:"omitempty,dive"`
		Title   string   `form:"title,omitempty"`
		OrderBy []string `form:"[]order_by,omitempty"`
		PaginationQuery
	}
)
