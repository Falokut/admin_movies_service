package service_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Falokut/admin_movies_service/internal/events"
	mock_events "github.com/Falokut/admin_movies_service/internal/events/mocks"
	"github.com/Falokut/admin_movies_service/internal/repository"
	mock_repository "github.com/Falokut/admin_movies_service/internal/repository/mocks"
	"github.com/Falokut/admin_movies_service/internal/service"
	mock_service "github.com/Falokut/admin_movies_service/internal/service/mocks"
	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	gomock "github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func getNullLogger() *logrus.Logger {
	l, _ := test.NewNullLogger()
	return l
}

func newServer(t *testing.T, register func(srv *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	return conn
}

func newClient(t *testing.T, s *service.MoviesService) *grpc.ClientConn {
	return newServer(t, func(srv *grpc.Server) { admin_movies_service.RegisterMoviesServiceV1Server(srv, s) })
}

type GetMovieBehavior func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMovieRequest, out *admin_movies_service.Movie)
type GetMoviesBehavior func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMoviesRequest, expectedResponce *admin_movies_service.Movies)

func getNullEvents(ctrl *gomock.Controller) events.MoviesEventsMQ {
	serv := mock_events.NewMockMoviesEventsMQ(ctrl)
	serv.EXPECT().MovieDeleted(gomock.Any(), gomock.Any()).Times(0)
	return serv
}

func getDefaultMoviesService(t *testing.T, moviesRepo service.MoviesRepository,
	countriesRepo repository.CountriesRepository,
	genresRepo repository.GenresRepository,
	moviesEvents events.MoviesEventsMQ) (admin_movies_service.MoviesServiceV1Client, *grpc.ClientConn) {
	conn := newClient(t,
		service.NewMoviesService(getNullLogger(),
			moviesRepo,
			countriesRepo,
			genresRepo,
			moviesEvents))

	client := admin_movies_service.NewMoviesServiceV1Client(conn)
	assert.NotNil(t, client)
	return client, conn
}

type GetMovieRequestMatcher struct {
	in *admin_movies_service.GetMovieRequest
}

func (r *GetMovieRequestMatcher) Match(in *admin_movies_service.GetMovieRequest) bool {
	if r.in == nil && in == nil {
		return true
	}
	if (r.in != nil && in == nil) || (in != nil && r.in == nil) {
		return false
	}

	return r.in.MovieID == in.MovieID

}

func TestGetMovie(t *testing.T) {
	testCases := []struct {
		movieID          int32
		behavior         GetMovieBehavior
		expectedStatus   codes.Code
		expectedResponce *admin_movies_service.Movie
		expectedError    error
		msg              string
	}{
		{
			movieID: 1,
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMovieRequest, out *admin_movies_service.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Return(out, service.ErrNotFound).Times(1)
			},
			expectedStatus:   codes.NotFound,
			expectedError:    service.ErrNotFound,
			expectedResponce: nil,
			msg:              "Test case num %d, must return not found error, if movie not found",
		},
		{
			movieID: 1,
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMovieRequest, out *admin_movies_service.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Return(out, context.Canceled).Times(1)
			},
			expectedStatus:   codes.Canceled,
			expectedError:    context.Canceled,
			expectedResponce: nil,
			msg:              "Test case num %d, must return context error, if repo return context error",
		},
		{
			movieID: 10,
			expectedResponce: &admin_movies_service.Movie{
				TitleRU:          "TitleRU",
				PosterURL:        "someurl",
				PreviewPosterURL: "someurl",
				BackgroundURL:    "someurl",
				Description:      "Description",
				ShortDescription: "ShortDescription",
				Duration:         100,
				ReleaseYear:      2000,
			},
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMovieRequest, out *admin_movies_service.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Return(out, nil).Times(1)
			},
			expectedStatus: codes.OK,
			msg:            "Test case num %d, must return expected responce, if repo doesn't return error, service shouldn't change data, except for the link to the poster",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		repo := mock_service.NewMockMoviesRepository(ctrl)

		req := &admin_movies_service.GetMovieRequest{
			MovieID: testCase.movieID,
		}

		testCase.behavior(repo, req, testCase.expectedResponce)
		client, conn := getDefaultMoviesService(t,
			repo, mock_repository.NewMockCountriesRepository(ctrl),
			mock_repository.NewMockGenresRepository(ctrl), getNullEvents(ctrl))
		defer conn.Close()

		res, err := client.GetMovie(context.Background(), req)

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
		assert.Equal(t, testCase.expectedStatus, status.Code(err), testCase.msg)
	}
}

type GetMoviesRequestMatcher struct {
	in                 *admin_movies_service.GetMoviesRequest
	maxLimit, minLimit uint32
}

func (r *GetMoviesRequestMatcher) Match(in *admin_movies_service.GetMoviesRequest) bool {
	if r.in == nil && in == nil {
		return true
	}
	if (r.in != nil && in == nil) || (in != nil && r.in == nil) {
		return false
	}

	if r.in.Limit <= r.minLimit {
		r.in.Limit = r.minLimit
	} else if r.in.Limit > r.maxLimit {
		r.in.Limit = r.maxLimit
	}

	return r.in.GetMoviesIDs() == in.GetMoviesIDs() &&
		r.in.GetCountriesIDs() == in.GetCountriesIDs() &&
		r.in.GetGenresIDs() == in.GetGenresIDs() &&
		r.in.GetTitle() == in.GetTitle() &&
		r.in.GetAgeRatings() == in.GetAgeRatings() &&
		r.in.Offset == in.Offset &&
		r.in.Limit == in.GetLimit()

}

func TestGetMovies(t *testing.T) {
	type MoviesRequest struct {
		MoviesIDs    string
		GenresIDs    string
		CountriesIDs string
		Title        string
		Limit        uint32
		Offset       uint32
	}

	testCases := []struct {
		request          MoviesRequest
		behavior         GetMoviesBehavior
		expectedStatus   codes.Code
		expectedResponce *admin_movies_service.Movies
		expectedError    error
		msg              string
	}{
		{
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMoviesRequest, expectedResponce *admin_movies_service.Movies) {
				m.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(expectedResponce, service.ErrNotFound).Times(1)
			},
			expectedStatus:   codes.NotFound,
			expectedError:    service.ErrNotFound,
			expectedResponce: nil,
			msg:              "Test case num %d, must return not found error, if movies not found",
		},
		{
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMoviesRequest, expectedResponce *admin_movies_service.Movies) {
				m.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(expectedResponce, context.Canceled).Times(1)
			},
			expectedStatus:   codes.Canceled,
			expectedError:    context.Canceled,
			expectedResponce: nil,
			msg:              "Test case num %d, must return context error, if repo return context error",
		},
		{
			request: MoviesRequest{
				MoviesIDs: "10,12",
				Limit:     60,
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
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMoviesRequest, expectedResponce *admin_movies_service.Movies) {
				m.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(expectedResponce, nil).Times(1)
			},
			expectedStatus: codes.OK,
			msg:            "Test case num %d, must return expected responce, if repo doesn't return error, service shouldn't change data, except for the link to the poster",
		},

		{
			request: MoviesRequest{
				MoviesIDs: "10,120",
				Limit:     110,
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
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMoviesRequest, expectedResponce *admin_movies_service.Movies) {
				m.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(expectedResponce, nil).Times(1)
			},
			expectedStatus: codes.OK,
			msg: "Test case num %d, must return expected responce, limit should be in [10,100]," +
				"if repo doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			request:          MoviesRequest{MoviesIDs: "10-,2-12"},
			expectedResponce: nil,
			behavior: func(m *mock_service.MockMoviesRepository, in *admin_movies_service.GetMoviesRequest, expectedResponce *admin_movies_service.Movies) {
				m.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(expectedResponce, service.ErrInvalidArgument).Times(1)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrInvalidArgument,
			msg:            "Test case num %d, must return expected error, if receive InvalidArgument error",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		repo := mock_service.NewMockMoviesRepository(ctrl)
		client, conn := getDefaultMoviesService(t,
			repo, mock_repository.NewMockCountriesRepository(ctrl),
			mock_repository.NewMockGenresRepository(ctrl), getNullEvents(ctrl))
		defer conn.Close()

		req := &admin_movies_service.GetMoviesRequest{
			MoviesIDs:    &testCase.request.MoviesIDs,
			GenresIDs:    &testCase.request.GenresIDs,
			CountriesIDs: &testCase.request.CountriesIDs,
			Title:        &testCase.request.Title,
			Limit:        testCase.request.Limit,
			Offset:       testCase.request.Offset,
		}

		testCase.behavior(repo, req, testCase.expectedResponce)
		res, err := client.GetMovies(context.Background(), req)
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
		assert.Equal(t, testCase.expectedStatus, status.Code(err), testCase.msg)
	}
}

func isProtoMoviesEqual(t *testing.T, expected, result *admin_movies_service.Movie) bool {
	if expected == nil && result == nil {
		return true
	} else if expected == nil && result != nil ||
		expected != nil && result == nil {
		return false
	}
	return assert.Equal(t, expected.Description, result.Description, "descriptions not equals") &&
		assert.Equal(t, expected.ShortDescription, result.ShortDescription, "short descriptions not equals") &&
		assert.Equal(t, expected.TitleRU, result.TitleRU, "ru titles not equals") &&
		assert.Equal(t, expected.TitleEN, result.TitleEN, "en titles not equals") &&
		assert.Equal(t, expected.Genres, result.Genres, "genres not equals") &&
		assert.Equal(t, expected.Duration, result.Duration, "duration not equals") &&
		assert.Equal(t, expected.Countries, result.Countries, "countries not equals") &&
		assert.Equal(t, expected.PosterURL, result.PosterURL, "posters urls not equals") &&
		assert.Equal(t, expected.PreviewPosterURL, result.PreviewPosterURL, "preview posters urls not equals") &&
		assert.Equal(t, expected.BackgroundURL, result.BackgroundURL, "backgrounds not equals") &&
		assert.Equal(t, expected.ReleaseYear, result.ReleaseYear, "release years not equals")
}
