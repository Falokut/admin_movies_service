package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type moviesRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	moviesTableName          = "movies"
	ageRatingsTableName      = "age_ratings"
	moviesGenresTableName    = "movies_genres"
	moviesCountriesTableName = "movies_countries"
)

func NewMoviesRepository(db *sqlx.DB, logger *logrus.Logger) *moviesRepository {
	return &moviesRepository{db: db, logger: logger}
}

func (r *moviesRepository) Shutdown() {
	r.db.Close()
}

const movieFields = `title_ru, title_en, description, short_description, duration,
poster_picture_id, background_picture_id, preview_poster_picture_id, release_year`

func (r *moviesRepository) GetMovie(ctx context.Context, movieID int32) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT %[1]s.id, %[2]s, COALESCE(%[3]s.name,'') AS age_rating "+
		"FROM %[1]s LEFT JOIN %[3]s ON age_rating_id=%[3]s.id WHERE %[1]s.id=$1",
		moviesTableName, movieFields, ageRatingsTableName)
	var movie Movie
	err = r.db.GetContext(ctx, &movie, query, movieID)
	if errors.Is(err, sql.ErrNoRows) {
		return Movie{}, ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: movie_id: %d", err.Error(), query, movieID)
		return Movie{}, err
	}

	return movie, nil
}

func (r *moviesRepository) GetMovieDuration(ctx context.Context, movieID int32) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMovieDuration")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT duration FROM %s WHERE id=$1", moviesTableName)
	var duration uint32
	err = r.db.GetContext(ctx, &duration, query, movieID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: movie_id: %d", err.Error(), query, movieID)
		return 0, err
	}

	return duration, nil
}
func (r *moviesRepository) GetMovies(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMovies")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT %[1]s.id, %[2]s, COALESCE(%[3]s.name,'') AS age_rating 
		FROM %[1]s LEFT JOIN %[3]s ON age_rating_id=%[3]s.id
		%[4]s 
		ORDER BY id
		LIMIT %d OFFSET %d`,
		moviesTableName, movieFields, ageRatingsTableName, convertFilterToWhere(filter), limit, offset)
	var movies []Movie
	err = r.db.SelectContext(ctx, &movies, query)
	if errors.Is(err, sql.ErrNoRows) {
		return []Movie{}, ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s", err.Error(), query)
		return []Movie{}, err
	}

	return movies, nil
}

func (r *moviesRepository) IsAgeRatingAlreadyExists(ctx context.Context, name string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.IsAgeRatingAlreadyExists")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1", ageRatingsTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		r.logger.Errorf("%v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *moviesRepository) CreateAgeRating(ctx context.Context, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.CreateAgeRating")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1)", ageRatingsTableName)
	_, err = r.db.ExecContext(ctx, query, name)
	if err != nil {
		r.logger.Errorf("%v query: %s", err.Error(), query)
		return err
	}
	return nil
}

func (r *moviesRepository) DeleteAgeRating(ctx context.Context, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.DeleteAgeRating")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("DELETE FROM %s WHERE name=$1 RETURNING id", ageRatingsTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s", err.Error(), query)
		return err
	}
	return nil
}

type insertIntoMoviesDtParam struct {
	TitleRU             string `db:"title_ru"`
	TitleEN             string `db:"title_en"`
	Description         string `db:"description"`
	Duration            int32  `db:"duration"`
	PosterID            string `db:"poster_picture_id"`
	PreviewPosterID     string `db:"preview_poster_picture_id"`
	BackgroundPictureID string `db:"background_picture_id"`
	ShortDescription    string `db:"short_description"`
	ReleaseYear         int32  `db:"release_year"`
	AgeRating           int32  `db:"age_rating_id"`
}

func (r *moviesRepository) CreateMovie(ctx context.Context, movie CreateMovieParam) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.CreateMovie")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	args, fields, values := getInsertStatement(newInsertIntoMoviesDtParam(movie))
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) RETURNING id", moviesTableName, fields, values)
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()

	if err != nil {
		r.logger.Errorf("can't begin transaction err: %v", err.Error())
		return 0, err
	}
	var movieID int32
	err = r.db.GetContext(ctx, &movieID, query, args...)
	if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, args)
		return 0, err
	}

	errorCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		insertStatements := make([]string, 0, len(movie.Genres))
		for _, id := range movie.Genres {
			insertStatements = append(insertStatements, fmt.Sprintf("(%d,%d)", movieID, id))
		}
		query = fmt.Sprintf("INSERT INTO %s (movie_id, genre_id) VALUES %s ON CONFLICT DO NOTHING;",
			moviesGenresTableName, strings.Join(insertStatements, ","))
		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			r.logger.Errorf("%v query: %s", err.Error(), query)
			errorCh <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		insertStatements := make([]string, 0, len(movie.CountriesIDs))
		for _, id := range movie.CountriesIDs {
			insertStatements = append(insertStatements, fmt.Sprintf("(%d,%d)", movieID, id))
		}
		query = fmt.Sprintf("INSERT INTO %s (movie_id, country_id) VALUES %s ON CONFLICT DO NOTHING;",
			moviesCountriesTableName, strings.Join(insertStatements, ","))
		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			r.logger.Errorf("%v query: %s", err.Error(), query)
			errorCh <- err
		}
	}()

	go func() {
		wg.Wait()
		r.logger.Debug("channel closed")
		close(errorCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case err, open := <-errorCh:
			if err != nil {
				return 0, err
			}
			if !open {
				tx.Commit()
				return movieID, nil
			}
		}
	}
}

func newInsertIntoMoviesDtParam(p CreateMovieParam) insertIntoMoviesDtParam {
	return insertIntoMoviesDtParam{
		TitleRU:             p.TitleRU,
		TitleEN:             p.TitleEN,
		Description:         p.Description,
		Duration:            p.Duration,
		PosterID:            p.PosterID,
		PreviewPosterID:     p.PreviewPosterID,
		BackgroundPictureID: p.BackgroundPictureID,
		ShortDescription:    p.ShortDescription,
		ReleaseYear:         p.ReleaseYear,
		AgeRating:           p.AgeRating,
	}
}

func (r *moviesRepository) GetAgeRating(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetAgeRatings")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1", ageRatingsTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: name: %s", err.Error(), query, name)
		return 0, err
	}
	return id, err
}

func (r *moviesRepository) GetAgeRatings(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetAgeRatings")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT name FROM %s", ageRatingsTableName)
	var ratings = []string{}
	err = r.db.SelectContext(ctx, &ratings, query)
	return ratings, err
}

// if any filter param filled, return string with WHERE statement
func convertFilterToWhere(filter MoviesFilter) string {
	statement := ""
	var first = true
	if len(filter.MoviesIDs) > 0 {
		statement += fmt.Sprintf(" %s.id IN(%s) ", moviesTableName, filter.MoviesIDs)
		first = false
	}

	str, first := getGenresFilter(filter.GenresIDs, first)
	statement += str

	str, first = getCountriesFilter(filter.CountriesIDs, first)
	statement += str

	str, first = containsInArray(ageRatingsTableName+".name", filter.AgeRating, first)
	statement += str

	if filter.Title != "" {
		filter.Title = strings.ReplaceAll(strings.ToLower(filter.Title), "'", "''") + "%"
		str = fmt.Sprintf(" LOWER(title_ru) LIKE('%[1]s') OR LOWER(title_en) LIKE('%[1]s')", filter.Title)
		if !first {
			statement += " AND " + str
		} else {
			statement += str
		}
	}

	if statement == "" {
		return statement
	}

	return " WHERE " + statement
}

func getGenresFilter(ids string, first bool) (string, bool) {
	if len(ids) == 0 {
		return "", first
	}

	filter := fmt.Sprintf("%s.id=ANY(SELECT movie_id FROM %s WHERE genre_id=ANY(ARRAY[%s]))", moviesTableName, moviesGenresTableName, ids)
	if first {
		first = false
	} else {
		filter = " AND " + filter
	}
	return filter, first
}

func getCountriesFilter(ids string, first bool) (string, bool) {
	if len(ids) == 0 {
		return "", first
	}

	filter := fmt.Sprintf("%s.id=ANY(SELECT movie_id FROM %s WHERE country_id=ANY(ARRAY[%s]))", moviesTableName, moviesCountriesTableName, ids)
	if first {
		first = false
	} else {
		filter = " AND " + filter
	}
	return filter, first
}

func containsInArray(fieldname string, array string, first bool) (string, bool) {
	if fieldname == "" || array == "" {
		return "", first
	}
	str := fmt.Sprintf(" %s=ANY('{%s}')", fieldname, array)
	if first {
		return str, !first
	}
	return " AND " + str, false
}

func getInsertStatement[T comparable](val T) ([]any, string, string) {
	rv := reflect.ValueOf(val)
	rt := rv.Type()
	var fields = make([]string, 0, rt.NumField())
	var args = make([]any, 0, rt.NumField())
	var values = make([]string, 0, rt.NumField())
	index := 1

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		v := rv.Field(i).Interface()
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice, reflect.Array:
			if reflect.ValueOf(v).Len() == 0 {
				continue
			}
		default:
			if isDefaultValue(v) {
				continue
			}
		}

		fields = append(fields, field.Tag.Get("db"))
		values = append(values, fmt.Sprintf("$%d", index))
		args = append(args, v)
		index++
	}

	return args, strings.Join(fields, ", "), strings.Join(values, ", ")
}

func (r *moviesRepository) IsMovieAlreadyExists(ctx context.Context,
	existParam IsMovieAlreadyExistsParam) (bool, []int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.IsMovieAlreadyExists")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	whereStatement, args := getWhereStatement(existParam)
	query := fmt.Sprintf("SELECT id FROM %s %s", moviesTableName, whereStatement)
	if len(args) == 0 {
		return false, []int32{}, ErrInvalidArgument
	}
	var ids []int32
	err = r.db.SelectContext(ctx, &ids, query, args...)
	if err != nil {
		r.logger.Errorf("%v query: %s args: %v", err.Error(), query, args)
		return false, []int32{}, err
	} else if len(ids) == 0 {
		return false, []int32{}, nil
	}

	return true, ids, nil
}

func isDefaultValue(field interface{}) bool {
	fieldVal := reflect.ValueOf(field)

	return !fieldVal.IsValid() || fieldVal.Interface() == reflect.Zero(fieldVal.Type()).Interface()
}

func getWhereStatement[T comparable](val T) (string, []any) {
	rv := reflect.ValueOf(val)
	rt := rv.Type()

	statements := make([]string, 0, rt.NumField())
	args := make([]any, 0, rt.NumField())
	index := 1

	for i := 0; i < rt.NumField(); i++ {
		v := rv.Field(i).Interface()

		if isDefaultValue(v) {
			continue
		}

		statements = append(statements, fmt.Sprintf("%s=$%d", rt.Field(i).Tag.Get("db"), index))
		args = append(args, v)
		index++
	}

	return " WHERE " + strings.Join(statements, " AND "), args
}

func (r *moviesRepository) DeleteMovie(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.DeleteMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING id", moviesTableName)

	var deletedId int32
	err = r.db.GetContext(ctx, &deletedId, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: %d", err.Error(), query, id)
		return err
	}
	return nil
}

func (r *moviesRepository) IsMovieExists(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.IsMovieExists")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", moviesTableName)

	var gettedId int32
	err = r.db.GetContext(ctx, &gettedId, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: %d", err.Error(), query, id)
		return false, err
	}

	return true, nil
}

func (r *moviesRepository) UpdatePictures(ctx context.Context, id int32,
	posterNameID, previewPosterID, backgroundID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.UpdatePictures")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	var fields []string
	var args []any
	i := 1
	if posterNameID != "" {
		fields = append(fields, fmt.Sprintf("poster_picture_id=$%d", i))
		args = append(args, posterNameID)
		i++
	}
	if previewPosterID != "" {
		fields = append(fields, fmt.Sprintf("preview_poster_picture_id=$%d", i))
		args = append(args, previewPosterID)
		i++
	}
	if backgroundID != "" {
		fields = append(fields, fmt.Sprintf("background_picture_id=$%d", i))
		args = append(args, backgroundID)
		i++
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING id", moviesTableName, strings.Join(fields, ","), i)

	var updated int32
	err = r.db.GetContext(ctx, &updated, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: %d", err.Error(), query, args)
		return err
	}

	return nil
}

func (r *moviesRepository) UpdateMovie(ctx context.Context, id int32, param UpdateMovieParam, genres, countries []int32,
	updateGenres, updateCountries bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.UpdateMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		r.logger.Errorf("error while creating transaction: %v", err)
		return err
	}
	errorCh := make(chan error, 1)
	var wg sync.WaitGroup

	statement, args := getUpdateStatement(param)
	if len(args) > 0 {
		wg.Add(1)
		args = append(args, id)
		go func() {
			defer wg.Done()
			query := fmt.Sprintf("UPDATE %s %s WHERE id=$%d", moviesTableName, statement, len(args))
			r.logger.Errorf("query: %s args: %v", query, args)
			_, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				r.logger.Errorf("%v query: %s", err.Error(), query)
				errorCh <- err
			}
		}()
	}

	if updateGenres {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(genres) > 0 {
				insertStatements := make([]string, 0, len(genres))
				for _, genreId := range genres {
					insertStatements = append(insertStatements, fmt.Sprintf("(%d,%d)", id, genreId))
				}
				query := fmt.Sprintf("INSERT INTO %s (movie_id, genre_id) VALUES %s ON CONFLICT DO NOTHING;",
					moviesGenresTableName, strings.Join(insertStatements, ","))
				_, err = tx.ExecContext(ctx, query)
				if err != nil {
					r.logger.Errorf("%v query: %s", err.Error(), query)
					errorCh <- err
				}
			}

			query := fmt.Sprintf("DELETE FROM %s WHERE movie_id=$1 AND NOT genre_id=ANY($2)", moviesGenresTableName)
			_, err = tx.ExecContext(ctx, query, id, genres)
			if err != nil {
				r.logger.Errorf("%v query: %s", err.Error(), query)
				errorCh <- err
			}
		}()
	}

	if updateCountries {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if len(countries) > 0 {
				insertStatements := make([]string, 0, len(countries))
				for _, countryId := range countries {
					insertStatements = append(insertStatements, fmt.Sprintf("(%d,%d)", id, countryId))
				}
				query := fmt.Sprintf("INSERT INTO %s (movie_id, country_id) VALUES %s ON CONFLICT DO NOTHING;",
					moviesCountriesTableName, strings.Join(insertStatements, ","))
				_, err = tx.ExecContext(ctx, query)
				if err != nil {
					r.logger.Errorf("%v query: %s", err.Error(), query)
					errorCh <- err
				}
			}

			query := fmt.Sprintf("DELETE FROM %s WHERE movie_id=$1 AND NOT country_id=ANY($2)", moviesCountriesTableName)
			_, err = tx.ExecContext(ctx, query, id, countries)
			if err != nil {
				r.logger.Errorf("%v query: %s", err.Error(), query)
				errorCh <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		r.logger.Debug("channel closed")
		close(errorCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, open := <-errorCh:
			if err != nil {
				return err
			}
			if !open {
				tx.Commit()
				return nil
			}
		}
	}

}

type moviePictures struct {
	PosterID            string `db:"poster_picture_id"`
	PreviewPosterID     string `db:"preview_poster_picture_id"`
	BackgroundPictureID string `db:"background_picture_id"`
}

func (r *moviesRepository) GetPicturesIds(ctx context.Context, id int32) (poster, preview, background string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetPicturesIds")
	defer span.Finish()
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT poster_picture_id, preview_poster_picture_id,background_picture_id 
		FROM %s WHERE id=$1`, moviesTableName)
	var movie moviePictures
	err = r.db.GetContext(ctx, &movie, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", "", ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("%v query: %s args: movie_id: %d", err.Error(), query, id)
		return
	}

	return movie.PosterID, movie.PreviewPosterID, movie.BackgroundPictureID, nil
}

func getUpdateStatement[T comparable](val T) (string, []any) {
	rv := reflect.ValueOf(val)
	rt := rv.Type()

	statements := make([]string, 0, rt.NumField())
	args := make([]any, 0, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		v := rv.Field(i).Interface()

		if isDefaultValue(v) {
			continue
		}
		args = append(args, v)
		statements = append(statements, fmt.Sprintf("%s=$%d", rt.Field(i).Tag.Get("db"), len(args)))
	}

	return " SET " + strings.Join(statements, ","), args
}
