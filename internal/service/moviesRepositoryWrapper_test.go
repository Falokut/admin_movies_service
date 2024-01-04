package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/Falokut/admin_movies_service/internal/repository"
	mock_repository "github.com/Falokut/admin_movies_service/internal/repository/mocks"
	"github.com/Falokut/admin_movies_service/internal/service"
	mock_service "github.com/Falokut/admin_movies_service/internal/service/mocks"
	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

type GetMovieRepositoryBehavior func(r *mock_repository.MockMoviesRepository, movieId int32, expectedMovie repository.Movie)
type GetMoviesRepositoryBehavior func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
	expectedMovies []repository.Movie, limit, offset uint32)
type GetPictureURLMultipleBehavior func(s *mock_service.MockImagesService, times int)
type GetCountriesForMovieBehavior func(s *mock_repository.MockCountriesRepository, movieID int32, namesToReturn []string)
type GetCountriesForMoviesBehavior func(s *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string)
type GetGenresForMovieBehavior func(s *mock_repository.MockGenresRepository, movieID int32, namesToReturn []string)
type GetGenresForMoviesBehavior func(s *mock_repository.MockGenresRepository, namesToReturn map[int32][]string)

var AnyGetPictureURLMultipleBehavior func(s *mock_service.MockImagesService) = func(s *mock_service.MockImagesService) {
	s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Return("someurl").AnyTimes()
}

func getNullExistanceCheckerWithZeroCallTimes(ctrl *gomock.Controller) service.ExistanceChecker {
	serv := mock_service.NewMockExistanceChecker(ctrl)
	serv.EXPECT().CheckExistance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	return serv
}

func getDefaultMoviesRepositoryWrapper(t *testing.T,
	moviesRepo repository.MoviesRepository,
	countriesRepo repository.CountriesRepository,
	genresRepo repository.GenresRepository,
	imagesService service.ImagesService,
	existanceChecker service.ExistanceChecker) service.MoviesRepository {
	return service.NewMoviesRepositoryWrapper(countriesRepo,
		genresRepo, moviesRepo,
		imagesService,
		service.PictureConfig{},
		service.PictureConfig{},
		service.PictureConfig{},
		existanceChecker,
		getNullLogger())
}

//			if testCase.expectedError != nil {
//				if assert.Error(t, err) {
//					assert.Contains(t, err.Error(), testCase.expectedError.Error())
//				}
//			} else {
//				assert.NotNil(t, res)
//				var comp assert.Comparison = func() (success bool) {
//					return isProtoMoviesEqual(t, testCase.expectedResponce, res)
//				}
//				assert.Condition(t, comp, testCase.msg)
//			}
//		}
//	}
func TestMoviesWrapperGetMovie(t *testing.T) {
	defer goleak.VerifyNone(t)

	testCases := []struct {
		movieID               int32
		moviesRepoBehavior    GetMovieRepositoryBehavior
		countriesRepoBehavior GetCountriesForMovieBehavior
		genresRepoBehavior    GetGenresForMovieBehavior
		genres, countries     []string
		expectedResponce      *admin_movies_service.Movie
		movie                 repository.Movie
		expectedError         error
		urlRequestTimes       int
		msg                   string
	}{
		{
			movieID: 1,
			moviesRepoBehavior: func(r *mock_repository.MockMoviesRepository, movieId int32, expectedMovie repository.Movie) {
				r.EXPECT().GetMovie(gomock.Any(), movieId).Return(repository.Movie{}, repository.ErrNotFound).Times(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetGenresForMovie(gomock.Any(), movieID).MaxTimes(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetCountriesForMovie(gomock.Any(), movieID).MaxTimes(1)
			},
			expectedError:    service.ErrNotFound,
			expectedResponce: nil,
			msg:              "Test case num %d, must return not found error, if movie not found",
		},
		{
			movieID: 1,
			moviesRepoBehavior: func(r *mock_repository.MockMoviesRepository, movieId int32, expectedMovie repository.Movie) {
				r.EXPECT().GetMovie(gomock.Any(), movieId).Return(repository.Movie{}, context.Canceled).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetGenresForMovie(gomock.Any(), movieID).MaxTimes(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetCountriesForMovie(gomock.Any(), movieID).MaxTimes(1)
			},
			expectedError:    context.Canceled,
			expectedResponce: nil,
			msg:              "Test case num %d, must same error that repo",
		},
		{
			movieID: 10,
			movie: repository.Movie{
				ID:                  10,
				TitleRU:             "TitleRU",
				PosterID:            sql.NullString{String: "1012", Valid: true},
				BackgroundPictureID: sql.NullString{String: "1012", Valid: true},
				Description:         "Description",
				ShortDescription:    "ShortDescription",
				Duration:            100,
				ReleaseYear:         2000,
			},
			genres:    []string{"Genre1", "Genre2"},
			countries: []string{"Country1", "Country2"},
			expectedResponce: &admin_movies_service.Movie{
				TitleRU:          "TitleRU",
				PosterURL:        "someurl",
				PreviewPosterURL: "someurl",
				BackgroundURL:    "someurl",
				Description:      "Description",
				ShortDescription: "ShortDescription",
				Genres:           []string{"Genre1", "Genre2"},
				Countries:        []string{"Country1", "Country2"},
				Duration:         100,
				ReleaseYear:      2000,
			},
			moviesRepoBehavior: func(r *mock_repository.MockMoviesRepository, movieId int32, expectedMovie repository.Movie) {
				r.EXPECT().GetMovie(gomock.Any(), movieId).Return(expectedMovie, nil).Times(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetGenresForMovie(gomock.Any(), movieID).Return(namesToReturn, nil).Times(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetCountriesForMovie(gomock.Any(), movieID).Return(namesToReturn, nil).Times(1)
			},
			urlRequestTimes: 2,
			msg:             "Test case num %d, must return expected responce, if repo doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			movieID: 10,
			movie: repository.Movie{
				ID:                  10,
				TitleRU:             "TitleRU",
				PosterID:            sql.NullString{String: "1012", Valid: true},
				BackgroundPictureID: sql.NullString{String: "1012", Valid: true},
				Description:         "Description",
				ShortDescription:    "ShortDescription",
				Duration:            100,
				ReleaseYear:         2000,
			},
			expectedResponce: nil,
			expectedError:    errors.New("any error"),
			genres:           []string{"Genre1", "Genre2"},
			moviesRepoBehavior: func(r *mock_repository.MockMoviesRepository, movieId int32, expectedMovie repository.Movie) {
				r.EXPECT().GetMovie(gomock.Any(), movieId).Return(expectedMovie, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetGenresForMovie(gomock.Any(), movieID).Return(namesToReturn, nil).MaxTimes(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetCountriesForMovie(gomock.Any(), movieID).Return([]string{}, errors.New("any error")).Times(1)
			},
			urlRequestTimes: 0,
			msg:             "Test case num %d, must return expected error, if countriesRepo return error",
		},
		{
			movieID: 10,
			movie: repository.Movie{
				ID:                  10,
				TitleRU:             "TitleRU",
				PosterID:            sql.NullString{String: "1012", Valid: true},
				BackgroundPictureID: sql.NullString{String: "1012", Valid: true},
				Description:         "Description",
				ShortDescription:    "ShortDescription",
				Duration:            100,
				ReleaseYear:         2000,
			},
			expectedResponce: nil,
			countries:        []string{"Country1", "Country2"},
			moviesRepoBehavior: func(r *mock_repository.MockMoviesRepository, movieId int32, expectedMovie repository.Movie) {
				r.EXPECT().GetMovie(gomock.Any(), movieId).Return(expectedMovie, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetGenresForMovie(gomock.Any(), movieID).Return([]string{}, errors.New("any error")).Times(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, movieID int32, namesToReturn []string) {
				r.EXPECT().GetCountriesForMovie(gomock.Any(), movieID).Return(namesToReturn, nil).MaxTimes(1)
			},
			urlRequestTimes: 0,
			expectedError:   errors.New("any error"),
			msg:             "Test case num %d, must return expected error, if countriesRepo return error",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		moviesRepo := mock_repository.NewMockMoviesRepository(ctrl)
		genresRepo := mock_repository.NewMockGenresRepository(ctrl)
		countriesRepo := mock_repository.NewMockCountriesRepository(ctrl)

		imgServ := mock_service.NewMockImagesService(ctrl)

		AnyGetPictureURLMultipleBehavior(imgServ)
		existanceChecker := getNullExistanceCheckerWithZeroCallTimes(ctrl)
		testCase.moviesRepoBehavior(moviesRepo, testCase.movieID, testCase.movie)
		testCase.countriesRepoBehavior(countriesRepo, testCase.movieID, testCase.countries)
		testCase.genresRepoBehavior(genresRepo, testCase.movieID, testCase.genres)

		wrapper := getDefaultMoviesRepositoryWrapper(t, moviesRepo, countriesRepo, genresRepo, imgServ, existanceChecker)

		res, err := wrapper.GetMovie(context.Background(),
			&admin_movies_service.GetMovieRequest{
				MovieID: testCase.movieID,
			})

		testCase.msg = fmt.Sprintf(testCase.msg, i+1)
		if testCase.expectedError != nil {
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
			}
		} else if assert.NotNil(t, res) {
			var comp assert.Comparison = func() (success bool) {
				return isProtoMoviesEqual(t, testCase.expectedResponce, res)
			}
			assert.Condition(t, comp, testCase.msg)
		}
	}
}

func TestMoviesWrapperGetMovies(t *testing.T) {
	defer goleak.VerifyNone(t)
	type MoviesRequest struct {
		MoviesIDs    string
		GenresIDs    string
		CountriesIDs string
		Title        string
		Limit        uint32
		Offset       uint32
	}

	testCases := []struct {
		request               MoviesRequest
		urlRequestTimes       int
		moviesBehavior        GetMoviesRepositoryBehavior
		imgBehavior           GetPictureURLMultipleBehavior
		countriesRepoBehavior GetCountriesForMoviesBehavior
		genresRepoBehavior    GetGenresForMoviesBehavior
		genres, countries     map[int32][]string
		expectedResponce      *admin_movies_service.Movies
		movies                []repository.Movie
		expectedError         error
		msg                   string
	}{
		{
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Return([]repository.Movie{}, repository.ErrNotFound).MaxTimes(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedError:    service.ErrNotFound,
			expectedResponce: nil,
			msg:              "Test case num %d, must return not found error, if movies not found",
		},
		{
			movies: []repository.Movie{
				{ID: 1},
			},
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).MaxTimes(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, fmt.Errorf("any error")).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedError:    fmt.Errorf("any error"),
			expectedResponce: nil,
			msg:              "Test case num %d, must return expectedError, if repo return error != ErrNotFound",
		},
		{
			request: MoviesRequest{
				MoviesIDs: "10,12",
				Limit:     60,
			},
			movies: []repository.Movie{
				{
					ID:               10,
					TitleRU:          "TitleRU",
					ShortDescription: "ShortDescription",
					Duration:         100,
					ReleaseYear:      2000,
				},
				{
					ID:               12,
					TitleRU:          "TitleRU",
					Description:      "Description",
					TitleEN:          sql.NullString{String: "TitleEn", Valid: true},
					PreviewPosterID:  sql.NullString{String: "validString", Valid: true},
					ShortDescription: "ShortDescription",
					Duration:         150,
					ReleaseYear:      2200,
				},
			},

			expectedResponce: &admin_movies_service.Movies{
				Movies: map[int32]*admin_movies_service.Movie{
					10: {
						TitleRU:          "TitleRU",
						ShortDescription: "ShortDescription",
						Duration:         100,
						ReleaseYear:      2000,
					},

					12: {
						TitleEN:          "TitleEn",
						TitleRU:          "TitleRU",
						ShortDescription: "ShortDescription",
						Description:      "Description",
						Duration:         150,
						ReleaseYear:      2200,
					},
				},
			},
			urlRequestTimes: 6,
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).Times(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Return("").Times(times)
			},
			msg: "Test case num %d, must return expected responce, if repo doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			request: MoviesRequest{MoviesIDs: "99912,213", Limit: 190},
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Return([]repository.Movie{}, context.Canceled).Times(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedError:    context.Canceled,
			expectedResponce: nil,
			msg:              "Test case num %d, must return expected error, if repo return error != ErrNotFound",
		},
		{
			request: MoviesRequest{
				MoviesIDs: "10,12",
				GenresIDs: "10,2",
				Limit:     110,
			},

			movies:           []repository.Movie{},
			expectedResponce: nil,
			urlRequestTimes:  0,
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Return([]repository.Movie{},
					repository.ErrNotFound).Times(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedError: service.ErrNotFound,
			msg: "Test case num %d, must return expected error," +
				"if repo return empty Moves slice",
		},

		{
			request: MoviesRequest{
				MoviesIDs: "10,120",
				Limit:     110,
			},
			movies: []repository.Movie{
				{
					ID:               120,
					TitleRU:          "TitleRU",
					ShortDescription: "ShortDescription",
					Duration:         100,
					ReleaseYear:      2000,
				},
				{
					ID:               12,
					TitleRU:          "TitleRU",
					TitleEN:          sql.NullString{String: "TitleEn", Valid: true},
					ShortDescription: "ShortDescription",
					Duration:         150,
					ReleaseYear:      2200,
				},
			},

			expectedResponce: &admin_movies_service.Movies{
				Movies: map[int32]*admin_movies_service.Movie{
					120: {
						TitleRU:          "TitleRU",
						ShortDescription: "ShortDescription",
						Duration:         100,
						ReleaseYear:      2000,
					},

					12: {
						TitleEN:          "TitleEn",
						TitleRU:          "TitleRU",
						ShortDescription: "ShortDescription",
						Duration:         150,
						ReleaseYear:      2200,
					},
				},
			},
			urlRequestTimes: 6,
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).Times(1)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).MaxTimes(1)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Return("").Times(times)
			},
			msg: "Test case num %d, must return expected responce, limit should be in [10,100]," +
				"if repo doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			request:          MoviesRequest{MoviesIDs: "10-,2-12"},
			expectedResponce: nil,
			moviesBehavior: func(r *mock_repository.MockMoviesRepository, filter repository.MoviesFilter,
				expectedMovies []repository.Movie, limit, offset uint32) {
				r.EXPECT().GetMovies(gomock.Any(), filter, limit, offset).Times(0)
			},
			imgBehavior: func(s *mock_service.MockImagesService, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any()).Times(0)
			},
			countriesRepoBehavior: func(r *mock_repository.MockCountriesRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetCountriesForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).Times(0)
			},
			genresRepoBehavior: func(r *mock_repository.MockGenresRepository, namesToReturn map[int32][]string) {
				r.EXPECT().GetGenresForMovies(gomock.Any(), gomock.Any()).Return(namesToReturn, nil).Times(0)
			},
			expectedError: service.ErrInvalidFilter,
			msg:           "Test case num %d, must return expected error, if filter not valid",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		moviesRepo := mock_repository.NewMockMoviesRepository(ctrl)
		genresRepo := mock_repository.NewMockGenresRepository(ctrl)
		countriesRepo := mock_repository.NewMockCountriesRepository(ctrl)

		imgServ := mock_service.NewMockImagesService(ctrl)

		testCase.imgBehavior(imgServ, testCase.urlRequestTimes)
		existanceChecker := getNullExistanceCheckerWithZeroCallTimes(ctrl)
		if len(testCase.countries) == 0 {
			testCase.countries = make(map[int32][]string)
		}
		testCase.countriesRepoBehavior(countriesRepo, testCase.countries)
		if len(testCase.genres) == 0 {
			testCase.genres = make(map[int32][]string)
		}
		testCase.genresRepoBehavior(genresRepo, testCase.genres)

		client := getDefaultMoviesRepositoryWrapper(t, moviesRepo, countriesRepo, genresRepo, imgServ, existanceChecker)

		var limit = testCase.request.Limit
		if limit == 0 {
			limit = 10
		} else if limit > 100 {
			limit = 100
		}
		testCase.moviesBehavior(moviesRepo, repository.MoviesFilter{
			MoviesIDs:    testCase.request.MoviesIDs,
			GenresIDs:    testCase.request.GenresIDs,
			CountriesIDs: testCase.request.CountriesIDs,
			Title:        testCase.request.Title,
		}, testCase.movies,
			limit, testCase.request.Offset)

		res, err := client.GetMovies(context.Background(),
			&admin_movies_service.GetMoviesRequest{
				MoviesIDs:    &testCase.request.MoviesIDs,
				GenresIDs:    &testCase.request.GenresIDs,
				CountriesIDs: &testCase.request.CountriesIDs,
				Title:        &testCase.request.Title,
				Limit:        testCase.request.Limit,
				Offset:       testCase.request.Offset,
			})

		testCase.msg = fmt.Sprintf(testCase.msg, i+1)
		if testCase.expectedError != nil {
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
			}
		} else if assert.NotNil(t, res) && assert.Equal(t, len(testCase.expectedResponce.Movies), len(res.Movies)) {
			var comp assert.Comparison = func() (success bool) {
				for key, Expectedmovie := range testCase.expectedResponce.Movies {
					if !isProtoMoviesEqual(t, Expectedmovie, res.Movies[key]) {
						return false
					}
				}
				return true
			}
			assert.Condition(t, comp, testCase.msg)
		}
	}
}
