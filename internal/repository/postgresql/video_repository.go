package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func (vr *VideoRepository) FetchAll(ctx context.Context, params entity.GetVideosParams) (entity.GetVideosReturn, error) {
	baseQuery := `select id, title, overview, format, rent_price, production_company, cover_path, total_stock, available_stock, genre_ids
				from videos where deleted_at is null`
	totalRowQuery := `select count(id) as total_row from videos where deleted_at is null`

	// Filters
	paramsQuery := ""
	filterArgs := []any{}
	argIdx := 1
	if params.Title != "" {
		paramsQuery += fmt.Sprintf(" and title ilike $%d", argIdx)
		filterArgs = append(filterArgs, "%"+params.Title+"%")
		argIdx++
	}

	if len(params.GenreIDs) > 0 {
		paramsQuery += fmt.Sprintf(" and genre_ids && $%d", argIdx)
		filterArgs = append(filterArgs, params.GenreIDs)
		argIdx++
	}

	// Order & Sort
	sortArgs := []any{}
	direction := "asc"
	params.SortOrder = strings.ToLower(params.SortOrder)
	if params.SortOrder != "" && (params.SortOrder == "asc" || params.SortOrder == "desc") {
		direction = params.SortOrder
	}

	sortQuery := " order by id"
	// Validate against allowed columns
	allowedColumns := map[string]bool{"id": true, "title": true, "created_at": true, "total_stock": true, "available_stock": true, "rent_price": true}
	if len(params.OrderBy) > 0 {
		sortQuery = " order by "
		for _, column := range params.OrderBy {
			column = strings.ToLower(column)
			if allowedColumns[column] {
				sortQuery += fmt.Sprintf("%s %s,", column, direction)
			}
		}
		// remove extra "," from sortQuery
		sortQuery = strings.TrimRight(sortQuery, ",")
	}

	if params.Limit > 0 {
		sortQuery += fmt.Sprintf(" limit $%d", argIdx)
		sortArgs = append(sortArgs, params.Limit)
		argIdx++

		if params.Page > 0 {
			offset := (params.Page - 1) * params.Limit
			sortQuery += fmt.Sprintf(" offset $%d", argIdx)
			sortArgs = append(sortArgs, offset)
			argIdx++
		}
	}

	baseQuery += paramsQuery + sortQuery
	totalRowQuery += paramsQuery
	rows, rowsErr := vr.conn.QueryContext(ctx, baseQuery, append(filterArgs, sortArgs...)...)
	if rowsErr != nil {
		return entity.GetVideosReturn{}, customerrors.NewError(
			"failed to get videos data",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	var totalRow int
	if err := vr.conn.QueryRowContext(ctx, totalRowQuery, filterArgs...).Scan(&totalRow); err != nil {
		return entity.GetVideosReturn{}, customerrors.NewError(
			"failed to get total videos data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	m := pgtype.NewMap()
	videos := entity.Videos{}
	var video entity.Video
	for rows.Next() {
		if err := rows.Scan(
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
			return entity.GetVideosReturn{}, customerrors.NewError(
				"failed to fetch video data",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		videos = append(videos, video)
	}
	if err := rows.Err(); err != nil {
		return entity.GetVideosReturn{}, customerrors.NewError(
			"failed to fetch video data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	pageInfo := entity.PageInfo{
		Page:    params.Page,
		Limit:   params.Limit,
		OrderBy: params.OrderBy,
		Filters: []entity.PageFilter{
			{
				Field: "title",
				Value: params.Title,
			},
			{
				Field: "genre_ids",
				Value: params.GenreIDs,
			},
		},
		OrderSort: params.SortOrder,
		TotalRow:  totalRow,
	}

	videosReturn := entity.GetVideosReturn{
		PageInfo: pageInfo,
		Entries:  videos,
	}

	return videosReturn, nil
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
