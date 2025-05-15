package dto

type (
	CreateVideoReq struct {
		Title             string `json:"title" binding:"required"`
		Overview          string `json:"overview,omitempty"`
		Format            string `json:"format" binding:"eq=dvd|eq=bluray|eq=digital|eq=vhs"`
		TotalStock        int    `json:"total_stock" binding:"gte=0"`
		CoverPath         string `json:"cover_path,omitempty"`
		ProductionCompany string `json:"production_company,omitempty"`
	}

	CreateVideoRes struct {
		ID                string `json:"id" `
		Title             string `json:"title"`
		Overview          string `json:"overview,omitempty"`
		Format            string `json:"format"`
		TotalStock        int    `json:"total_stock"`
		AvailableStock    int    `json:"available_stock"`
		CoverPath         string `json:"cover_path,omitempty"`
		ProductionCompany string `json:"production_company,omitempty"`
	}

	VideosRes struct {
		ID                string `json:"id" `
		Title             string `json:"title"`
		Overview          string `json:"overview,omitempty"`
		Format            string `json:"format"`
		TotalStock        int    `json:"total_stock"`
		AvailableStock    int    `json:"available_stock"`
		CoverPath         string `json:"cover_path,omitempty"`
		ProductionCompany string `json:"production_company,omitempty"`
	}

	GetVideosRes []VideosRes
)
