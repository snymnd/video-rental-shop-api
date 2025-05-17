package dto

import "time"

type (
	CreateVideoReq struct {
		Title             string  `json:"title" binding:"required"`
		Overview          string  `json:"overview,omitempty"`
		RentPrice         float64 `json:"rent_price" binding:"required,gte=0"`
		Format            string  `json:"format" binding:"eq=dvd|eq=bluray|eq=digital|eq=vhs"`
		TotalStock        int     `json:"total_stock" binding:"required,gte=0"`
		CoverPath         *string `json:"cover_path,omitempty"`
		ProductionCompany *string `json:"production_company,omitempty"`
	}

	CreateVideoRes struct {
		ID                string  `json:"id" `
		Title             string  `json:"title"`
		Overview          string  `json:"overview,omitempty"`
		Format            string  `json:"format"`
		RentPrice         float64 `json:"rent_price"`
		TotalStock        int     `json:"total_stock"`
		AvailableStock    int     `json:"available_stock"`
		CoverPath         *string `json:"cover_path,omitempty"`
		ProductionCompany *string `json:"production_company,omitempty"`
	}

	VideoRes struct {
		ID                string     `json:"id" `
		Title             string     `json:"title"`
		Overview          string     `json:"overview,omitempty"`
		Format            string     `json:"format"`
		RentPrice         float64    `json:"rent_price"`
		TotalStock        int        `json:"total_stock"`
		AvailableStock    int        `json:"available_stock"`
		CoverPath         *string    `json:"cover_path,omitempty"`
		ProductionCompany *string    `json:"production_company,omitempty"`
		CreatedAt         *time.Time `json:"created_at,omitempty"`
		UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	}

	GetVideosRes []VideoRes
)
