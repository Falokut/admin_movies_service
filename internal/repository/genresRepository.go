package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type genresRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewGenresRepository(db *sqlx.DB, logger *logrus.Logger) *genresRepository {
	return &genresRepository{db: db, logger: logger}
}

func (r *genresRepository) Shutdown() error {
	return r.db.Close()
}

const (
	genresTableName = "genres"
)

func (r *genresRepository) GetGenre(ctx context.Context, id int32) (Genre, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetGenre")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", genresTableName)

	var genre Genre
	err = r.db.GetContext(ctx, &genre, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return Genre{}, ErrNotFound
	} else if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, id)
		return Genre{}, err
	}

	return genre, nil
}

func (r *genresRepository) GetGenreByName(ctx context.Context, name string) (Genre, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetGenreByName")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", genresTableName)

	var genre Genre
	err = r.db.GetContext(ctx, &genre, query, name)

	if errors.Is(err, sql.ErrNoRows) {
		return Genre{}, ErrNotFound
	} else if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, name)
		return Genre{}, err
	}

	return genre, nil
}

func (r *genresRepository) CreateGenre(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.CreateGenre")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING id;", genresTableName)

	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return int32(0), ErrNotFound
	} else if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, name)
		return 0, err
	}

	return id, nil
}

func (r *genresRepository) DeleteGenre(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.DeleteCity")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING id", genresTableName)
	var deletedID int32
	err = r.db.GetContext(ctx, &deletedID, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *genresRepository) IsGenreExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.IsGenreExist")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", genresTableName)
	var findedID int32
	err = r.db.GetContext(ctx, &findedID, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return false, nil
	} else if err != nil {
		return true, err
	}
	return true, nil
}

func (r *genresRepository) IsGenreAlreadyExists(ctx context.Context, name string) (bool, int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.IsGenreAlreadyExists")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1", genresTableName)
	var findedID int32
	err = r.db.GetContext(ctx, &findedID, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, 0, nil
	} else if err != nil {
		return true, 0, err
	}
	return true, findedID, nil
}

func (r *genresRepository) GetGenres(ctx context.Context) ([]Genre, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetGenres")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s", genresTableName)
	var genres []Genre
	err = r.db.SelectContext(ctx, &genres, query)
	if err != nil {
		return []Genre{}, err
	} else if len(genres) == 0 {
		return []Genre{}, ErrNotFound
	}
	return genres, nil
}

func (r *genresRepository) IsGenresExists(ctx context.Context, ids []int32) (bool, []int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.IsGenresExists")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT id FROM %s WHERE id=ANY($1)", genresTableName)
	var foundIDs []int32
	err = r.db.SelectContext(ctx, &foundIDs, query, ids)
	if err != nil {
		return false, []int32{}, err
	}

	IDsMap := make(map[int32]struct{}, len(ids))
	for _, id := range foundIDs {
		if _, ok := IDsMap[id]; !ok {
			IDsMap[id] = struct{}{}
		}
	}

	var notFoundIDs = make([]int32, 0, len(ids))
	for _, id := range ids {
		if _, ok := IDsMap[id]; !ok {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return len(notFoundIDs) == 0, notFoundIDs, nil
}

func (r *genresRepository) UpdateGenre(ctx context.Context, name string, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.UpdateGenre")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2", genresTableName)
	_, err = r.db.ExecContext(ctx, query, name, id)
	return err
}

func (r *genresRepository) GetGenresForMovies(ctx context.Context, ids []int32) (map[int32][]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetGenresForMovies")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT movie_id, ARRAY_AGG(name) FROM %s JOIN %s ON genre_id=id "+
		"WHERE movie_id=ANY($1) GROUP BY movie_id",
		moviesGenresTableName, genresTableName)
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return map[int32][]string{}, err
	}

	genres := make(map[int32][]string, len(ids))
	for rows.Next() {
		var id int32
		var names string

		rows.Scan(&id, &names)
		genres[id] = strings.Split(strings.Trim(names, "{}"), ",")
	}
	return genres, nil
}

func (r *genresRepository) GetGenresForMovie(ctx context.Context, id int32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetGenresForMovie")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT name FROM %s JOIN %s ON genre_id=id WHERE movie_id=$1",
		moviesGenresTableName, genresTableName)

	var genres []string
	err = r.db.SelectContext(ctx, &genres, query, id)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return []string{}, err
	}

	return genres, nil
}
