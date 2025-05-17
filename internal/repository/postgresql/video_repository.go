package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"

	"github.com/jackc/pgx/v5/pgtype"
)

type VideoRepository struct {
	conn *sql.DB
}

func NewVideoRepository(conn *sql.DB) *VideoRepository {
	return &VideoRepository{conn}
}

func (vr *VideoRepository) Create(ctx context.Context, video *entity.Video) error {
	query := `insert into videos(title, overview, format, rent_price, production_company, cover_path, total_stock, available_stock, genre_ids) 
				values 
				($1, $2, $3, $4, $5, $6, $7, $8, $9)
				returning id, created_at, updated_at, deleted_at`

	if err := vr.conn.QueryRowContext(ctx, query, video.Title, video.Overview, video.Format, video.RentPrice, video.ProductionCompany, video.CoverPath, video.TotalStock, video.AvailableStock, video.GenreIDs).
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
	query := `select id, title, overview, format, rent_price, production_company, cover_path, total_stock, available_stock, genre_ids
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

	m := pgtype.NewMap()
	var video entity.Video
	for rows.Next() {
		if err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.Overview,
			&video.Format,
			&video.RentPrice,
			&video.ProductionCompany,
			&video.CoverPath,
			&video.TotalStock,
			&video.AvailableStock,
			m.SQLScanner(&video.GenreIDs),
		); err != nil {
			return videos, customerrors.NewError(
				"failed to fetch video data",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return videos, customerrors.NewError(
			"failed to fetch video data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return videos, nil
}

func (vr *VideoRepository) FetchMultipleVideos(ctx context.Context, videosID []int) (videos entity.Videos, err error) {
	query := `select id, title, overview, format, production_company, rent_price, cover_path, total_stock, available_stock, genre_ids
				from videos
				where id in (`

	args := make([]any, len(videosID))
	for i, id := range videosID {
		query += fmt.Sprintf("$%d,", i+1)
		args[i] = id
	}
	query = query[:len(query)-1]
	query += ")"

	var conn DBinf = vr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
		query += " for update"
	}

	rows, rowsErr := conn.QueryContext(ctx, query, args...)
	if rowsErr != nil {
		return videos, customerrors.NewError(
			"failed to fetch videos",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	m := pgtype.NewMap()
	var video entity.Video
	for rows.Next() {
		if err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.Overview,
			&video.Format,
			&video.ProductionCompany,
			&video.RentPrice,
			&video.CoverPath,
			&video.TotalStock,
			&video.AvailableStock,
			m.SQLScanner(&video.GenreIDs),
		); err != nil {
			return videos, customerrors.NewError(
				"failed to scan video data",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return videos, customerrors.NewError(
			"failed to fetch videos data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return videos, nil
}

func (vr *VideoRepository) RentMultipleVideos(ctx context.Context, videosID []int) error {
	query := `update videos 
				set available_stock	= available_stock - 1
				where id in (`

	args := make([]any, len(videosID))
	for i, id := range videosID {
		query += fmt.Sprintf("$%d,", i+1)
		args[i] = id
	}
	query = query[:len(query)-1]
	query += ")"

	var conn DBinf = vr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	rows, rowsErr := conn.QueryContext(ctx, query, args...)
	if rowsErr != nil {
		return customerrors.NewError(
			"failed to update video available stock data",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return customerrors.NewError(
			"failed to update video available stock data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}
