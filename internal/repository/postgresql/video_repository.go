package postgresql

import (
	"context"
	"database/sql"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
)

type VideoRepository struct {
	conn *sql.DB
}

func NewVideoRepository(conn *sql.DB) *VideoRepository {
	return &VideoRepository{conn}
}

func (vr *VideoRepository) Create(ctx context.Context, video *entity.Video) error {
	query := `insert into videos(title, overview, format, production_company, cover_path, total_stock, available_stock) 
				values 
				($1, $2, $3, $4, $5, $6, $7)
				returning id, created_at, updated_at, deleted_at`

	if err := vr.conn.QueryRowContext(ctx, query, video.Title, video.Overview, video.Format, video.ProductionCompany, video.CoverPath, video.TotalStock, video.AvailableStock).
		Scan(
			&video.ID,
			&video.CreatedAt,
			&video.UpdatedAt,
			&video.DeletedAt,
		); err != nil {
		return customerrors.NewError(
			"failed to create video data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (vr *VideoRepository) FetchAll(ctx context.Context) (videos entity.Videos, err error) {
	query := `select id, title, overview, format, production_company, cover_path, total_stock, available_stock
				from videos`

	rows, rowsErr := vr.conn.QueryContext(ctx, query)
	if rowsErr != nil {
		return videos, customerrors.NewError(
			"failed to create video data",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	var video entity.Video
	for rows.Next() {
		if err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.Overview,
			&video.Format,
			&video.ProductionCompany,
			&video.CoverPath,
			&video.TotalStock,
			&video.AvailableStock,
		); err != nil {
			return videos, customerrors.NewError(
				"failed to create video data",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return videos, customerrors.NewError(
			"failed to create video data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return videos, nil
}
