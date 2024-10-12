package storage

import (
	"astrologerService/entity"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

func NewStorage(conn *pgxpool.Pool) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) SaveData(ctx context.Context, metadata entity.Metadata) error {
	query := `INSERT INTO metadata (date, explanation, hdurl, media_type, service_version, title, url)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.conn.Exec(
		ctx,
		query,
		metadata.Date,
		metadata.Explanation,
		metadata.HdURL,
		metadata.MediaType,
		metadata.ServiceVersion,
		metadata.Title,
		metadata.URL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetMetaDataBases(ctx context.Context) ([]entity.Metadata, error) {
	query := "SELECT date, explanation, hdurl, media_type, service_version, title, url FROM  metadata"

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var dataBase []entity.Metadata

	for rows.Next() {
		var data entity.Metadata
		err = rows.Scan(
			&data.Date,
			&data.Explanation,
			&data.HdURL,
			&data.MediaType,
			&data.ServiceVersion,
			&data.Title,
			&data.URL,
		)

		if err != nil {
			return nil, err
		}

		dataBase = append(dataBase, data)
	}
	return dataBase, nil
}

func (s *Storage) GetMetaData(ctx context.Context, date string) (entity.Metadata, error) {
	query := "SELECT  date, explanation, hdurl, media_type, service_version, title, url FROM  metadata WHERE date=$1"

	var data entity.Metadata

	err := s.conn.QueryRow(ctx, query, date).Scan(
		&data.Date,
		&data.Explanation,
		&data.HdURL,
		&data.MediaType,
		&data.ServiceVersion,
		&data.Title,
		&data.URL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Metadata{}, entity.ErrNotFound
		}

		return entity.Metadata{}, err
	}

	return data, nil
}
