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

type countriesRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewCountriesRepository(db *sqlx.DB, logger *logrus.Logger) *countriesRepository {
	return &countriesRepository{db: db, logger: logger}
}

func (r *countriesRepository) Shutdown() error {
	return r.db.Close()
}

const (
	countriesTableName = "countries"
)

func (r *countriesRepository) GetCountry(ctx context.Context, id int32) (Country, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetCountry")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", countriesTableName)

	var country Country
	err = r.db.GetContext(ctx, &country, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return Country{}, ErrNotFound
	} else if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, id)
		return Country{}, err
	}

	return country, nil
}

func (r *countriesRepository) GetCountryByName(ctx context.Context, name string) (Country, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetCountryByName")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s WHERE name_ru=$1 OR name_en=$1", countriesTableName)

	var country Country
	err = r.db.GetContext(ctx, &country, query, name)

	if errors.Is(err, sql.ErrNoRows) {
		return Country{}, ErrNotFound
	} else if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, name)
		return Country{}, err
	}

	return country, nil
}

func (r *countriesRepository) CreateCountry(ctx context.Context, countryName string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.CreateCountry")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING id;", countriesTableName)

	var id int32
	err = r.db.GetContext(ctx, &id, query, countryName)
	if errors.Is(err, sql.ErrNoRows) {
		return int32(0), ErrNotFound
	} else if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, countryName)
		return 0, err
	}

	return id, nil
}

func (r *countriesRepository) DeleteCountry(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.DeleteCity")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING id", countriesTableName)
	var deletedID int32
	err = r.db.GetContext(ctx, &deletedID, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *countriesRepository) UpdateCountry(ctx context.Context, countryName string, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.UpdateCountry")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2", countriesTableName)
	_, err = r.db.ExecContext(ctx, query, countryName, id)
	return err
}

func (r *countriesRepository) IsCountryExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.IsCountryExist")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", countriesTableName)
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

func (r *countriesRepository) IsCountryAlreadyExists(ctx context.Context, name string) (bool, int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.IsCountryAlreadyExists")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1", countriesTableName)
	var findedID int32
	err = r.db.GetContext(ctx, &findedID, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, 0, nil
	} else if err != nil {
		return true, 0, err
	}
	return true, findedID, nil
}

func (r *countriesRepository) GetCountries(ctx context.Context) ([]Country, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetCountries")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s", countriesTableName)
	var countries []Country
	err = r.db.SelectContext(ctx, &countries, query)
	if err != nil {
		return []Country{}, err
	}

	return countries, nil
}

func (r *countriesRepository) IsCountriesExists(ctx context.Context, ids []int32) (bool, []int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.IsCountriesExists")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT id FROM %s WHERE id=ANY($1)", countriesTableName)
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

func (r *countriesRepository) GetCountriesForMovies(ctx context.Context, ids []int32) (map[int32][]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetCountriesForMovies")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT movie_id, ARRAY_AGG(name) FROM %s JOIN %s ON country_id=id "+
		"WHERE movie_id=ANY($1) GROUP BY movie_id",
		moviesCountriesTableName, countriesTableName)
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return map[int32][]string{}, err
	}

	countries := make(map[int32][]string, len(ids))
	for rows.Next() {
		var id int32
		var names string

		rows.Scan(&id, &names)
		countries[id] = strings.Split(strings.Trim(names, "{}"), ",")
	}
	return countries, nil
}

func (r *countriesRepository) GetCountriesForMovie(ctx context.Context, id int32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetCountriesForMovie")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT name FROM %s JOIN %s ON country_id=id WHERE movie_id=$1",
		moviesCountriesTableName, countriesTableName)

	var countries []string
	err = r.db.SelectContext(ctx, &countries, query, id)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return []string{}, err
	}

	return countries, nil
}
