package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Falokut/admin_movies_service/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

const (
	moviesTableName     = "movies"
	ageRatingsTableName = "age_ratings"
)

func getNullLogger() *logrus.Logger {
	l, _ := test.NewNullLogger()
	return l
}

func getRowsForGetMovie(m repository.Movie) *sqlxmock.Rows {
	var columns = []string{"id", "age_rating",
		"title_ru", "title_en", "description", "short_description",
		"duration", "poster_picture_id", "background_picture_id",
		"preview_poster_picture_id", "release_year",
	}

	val := fmt.Sprintf("%d,%s,%s,%s,%s,%s,%d,%s,%s,%s,%d", m.ID,
		m.AgeRating, m.TitleRU, m.TitleEN.String, m.Description, m.ShortDescription,
		m.Duration, m.PosterID.String,
		m.BackgroundPictureID.String, m.PreviewPosterID.String, m.ReleaseYear)

	return sqlxmock.NewRows(columns).FromCSVString(val)
}

func TestGetMovie(t *testing.T) {
	var getMovieQuery string = fmt.Sprintf("^SELECT (.+) FROM %s LEFT JOIN %s (.+)$", moviesTableName, ageRatingsTableName)

	testCases := []struct {
		msg            string
		QueryError     error
		ID             int32
		ExpectedError  error
		ExpectedResult repository.Movie
	}{
		{
			QueryError:    sql.ErrNoRows,
			ID:            1,
			ExpectedError: repository.ErrNotFound,
			msg:           "test case num %d, repository must returns err not found, if driver returns error sql.ErrNoRows",
		},
		{
			QueryError:    sql.ErrTxDone,
			ID:            199123,
			ExpectedError: sql.ErrTxDone,
			msg:           "test case num %d, repository must returns same err,that driver returns, if error differ than sql.ErrNoRows",
		},
		{
			ID: 10,
			ExpectedResult: repository.Movie{
				ID:      10,
				TitleRU: "ruTitle",
				TitleEN: sql.NullString{String: "", Valid: false},
			},
			msg: "test case num %d, repository must returns expected responce",
		},
	}

	for i, testCase := range testCases {
		db, mock, err := sqlxmock.Newx()
		assert.NoError(t, err)
		defer db.Close()

		q := mock.ExpectQuery(getMovieQuery).WithArgs(testCase.ID)
		if testCase.QueryError != nil {
			q.WillReturnError(testCase.QueryError)
		} else {
			rows := getRowsForGetMovie(testCase.ExpectedResult)
			assert.NotNil(t, rows)
			q.WillReturnRows(rows)
		}
		repo := repository.NewMoviesRepository(db, getNullLogger())
		res, err := repo.GetMovie(context.Background(), testCase.ID)

		msg := fmt.Sprintf(testCase.msg, i+1)
		assert.ErrorIs(t, err, testCase.ExpectedError, msg)

		var comp assert.Comparison = func() (success bool) {
			return isMoviesFromDatabaseEqual(t, testCase.ExpectedResult, res)
		}

		assert.Condition(t, comp, msg)
		if err := mock.ExpectationsWereMet(); err != nil {
			assert.NoError(t, err, msg)
		}
	}
}

func isMoviesFromDatabaseEqual(t *testing.T, expected, result repository.Movie) bool {
	return assert.Equal(t, expected.Description, result.Description, "descriptions not equals") &&
		assert.Equal(t, expected.ShortDescription, result.ShortDescription, "short descriptions not equals") &&
		assert.Equal(t, expected.TitleRU, result.TitleRU, "ru titles not equals") &&
		assert.Equal(t, expected.TitleEN.String, result.TitleEN.String, "en titles not equals") &&
		assert.Equal(t, expected.Duration, result.Duration, "duration not equals") &&
		assert.Equal(t, expected.PosterID.String, result.PosterID.String, "posters not equals") &&
		assert.Equal(t, expected.PreviewPosterID.String, result.PreviewPosterID.String, "preview posters not equals") &&
		assert.Equal(t, expected.BackgroundPictureID.String, result.BackgroundPictureID.String, "backgrounds not equals") &&
		assert.Equal(t, expected.ReleaseYear, result.ReleaseYear, "release years not equals")
}
