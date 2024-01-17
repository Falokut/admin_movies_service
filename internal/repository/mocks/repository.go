// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	repository "github.com/Falokut/admin_movies_service/internal/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockMoviesRepository is a mock of MoviesRepository interface.
type MockMoviesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMoviesRepositoryMockRecorder
}

// MockMoviesRepositoryMockRecorder is the mock recorder for MockMoviesRepository.
type MockMoviesRepositoryMockRecorder struct {
	mock *MockMoviesRepository
}

// NewMockMoviesRepository creates a new mock instance.
func NewMockMoviesRepository(ctrl *gomock.Controller) *MockMoviesRepository {
	mock := &MockMoviesRepository{ctrl: ctrl}
	mock.recorder = &MockMoviesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMoviesRepository) EXPECT() *MockMoviesRepositoryMockRecorder {
	return m.recorder
}

// CreateAgeRating mocks base method.
func (m *MockMoviesRepository) CreateAgeRating(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAgeRating", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAgeRating indicates an expected call of CreateAgeRating.
func (mr *MockMoviesRepositoryMockRecorder) CreateAgeRating(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAgeRating", reflect.TypeOf((*MockMoviesRepository)(nil).CreateAgeRating), ctx, name)
}

// CreateMovie mocks base method.
func (m *MockMoviesRepository) CreateMovie(ctx context.Context, param repository.CreateMovieParam) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovie", ctx, param)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMovie indicates an expected call of CreateMovie.
func (mr *MockMoviesRepositoryMockRecorder) CreateMovie(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovie", reflect.TypeOf((*MockMoviesRepository)(nil).CreateMovie), ctx, param)
}

// DeleteAgeRating mocks base method.
func (m *MockMoviesRepository) DeleteAgeRating(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAgeRating", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAgeRating indicates an expected call of DeleteAgeRating.
func (mr *MockMoviesRepositoryMockRecorder) DeleteAgeRating(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAgeRating", reflect.TypeOf((*MockMoviesRepository)(nil).DeleteAgeRating), ctx, name)
}

// DeleteMovie mocks base method.
func (m *MockMoviesRepository) DeleteMovie(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMovie", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMovie indicates an expected call of DeleteMovie.
func (mr *MockMoviesRepositoryMockRecorder) DeleteMovie(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMovie", reflect.TypeOf((*MockMoviesRepository)(nil).DeleteMovie), ctx, id)
}

// GetAgeRating mocks base method.
func (m *MockMoviesRepository) GetAgeRating(ctx context.Context, name string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgeRating", ctx, name)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgeRating indicates an expected call of GetAgeRating.
func (mr *MockMoviesRepositoryMockRecorder) GetAgeRating(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgeRating", reflect.TypeOf((*MockMoviesRepository)(nil).GetAgeRating), ctx, name)
}

// GetAgeRatings mocks base method.
func (m *MockMoviesRepository) GetAgeRatings(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgeRatings", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgeRatings indicates an expected call of GetAgeRatings.
func (mr *MockMoviesRepositoryMockRecorder) GetAgeRatings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgeRatings", reflect.TypeOf((*MockMoviesRepository)(nil).GetAgeRatings), ctx)
}

// GetMovie mocks base method.
func (m *MockMoviesRepository) GetMovie(ctx context.Context, movieID int32) (repository.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovie", ctx, movieID)
	ret0, _ := ret[0].(repository.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovie indicates an expected call of GetMovie.
func (mr *MockMoviesRepositoryMockRecorder) GetMovie(ctx, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovie", reflect.TypeOf((*MockMoviesRepository)(nil).GetMovie), ctx, movieID)
}

// GetMovieDuration mocks base method.
func (m *MockMoviesRepository) GetMovieDuration(ctx context.Context, id int32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovieDuration", ctx, id)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovieDuration indicates an expected call of GetMovieDuration.
func (mr *MockMoviesRepositoryMockRecorder) GetMovieDuration(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovieDuration", reflect.TypeOf((*MockMoviesRepository)(nil).GetMovieDuration), ctx, id)
}

// GetMovies mocks base method.
func (m *MockMoviesRepository) GetMovies(ctx context.Context, Filter repository.MoviesFilter, limit, offset uint32) ([]repository.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies", ctx, Filter, limit, offset)
	ret0, _ := ret[0].([]repository.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockMoviesRepositoryMockRecorder) GetMovies(ctx, Filter, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockMoviesRepository)(nil).GetMovies), ctx, Filter, limit, offset)
}

// GetPicturesIds mocks base method.
func (m *MockMoviesRepository) GetPicturesIds(ctx context.Context, id int32) (string, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPicturesIds", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetPicturesIds indicates an expected call of GetPicturesIds.
func (mr *MockMoviesRepositoryMockRecorder) GetPicturesIds(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPicturesIds", reflect.TypeOf((*MockMoviesRepository)(nil).GetPicturesIds), ctx, id)
}

// IsAgeRatingAlreadyExists mocks base method.
func (m *MockMoviesRepository) IsAgeRatingAlreadyExists(ctx context.Context, name string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAgeRatingAlreadyExists", ctx, name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAgeRatingAlreadyExists indicates an expected call of IsAgeRatingAlreadyExists.
func (mr *MockMoviesRepositoryMockRecorder) IsAgeRatingAlreadyExists(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAgeRatingAlreadyExists", reflect.TypeOf((*MockMoviesRepository)(nil).IsAgeRatingAlreadyExists), ctx, name)
}

// IsMovieAlreadyExists mocks base method.
func (m *MockMoviesRepository) IsMovieAlreadyExists(ctx context.Context, existParam repository.IsMovieAlreadyExistsParam) (bool, []int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMovieAlreadyExists", ctx, existParam)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsMovieAlreadyExists indicates an expected call of IsMovieAlreadyExists.
func (mr *MockMoviesRepositoryMockRecorder) IsMovieAlreadyExists(ctx, existParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMovieAlreadyExists", reflect.TypeOf((*MockMoviesRepository)(nil).IsMovieAlreadyExists), ctx, existParam)
}

// IsMovieExists mocks base method.
func (m *MockMoviesRepository) IsMovieExists(ctx context.Context, id int32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMovieExists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMovieExists indicates an expected call of IsMovieExists.
func (mr *MockMoviesRepositoryMockRecorder) IsMovieExists(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMovieExists", reflect.TypeOf((*MockMoviesRepository)(nil).IsMovieExists), ctx, id)
}

// UpdateMovie mocks base method.
func (m *MockMoviesRepository) UpdateMovie(ctx context.Context, id int32, param repository.UpdateMovieParam, genres, countries []int32, updateGenres, updateCountries bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMovie", ctx, id, param, genres, countries, updateGenres, updateCountries)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMovie indicates an expected call of UpdateMovie.
func (mr *MockMoviesRepositoryMockRecorder) UpdateMovie(ctx, id, param, genres, countries, updateGenres, updateCountries interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMovie", reflect.TypeOf((*MockMoviesRepository)(nil).UpdateMovie), ctx, id, param, genres, countries, updateGenres, updateCountries)
}

// UpdatePictures mocks base method.
func (m *MockMoviesRepository) UpdatePictures(ctx context.Context, id int32, posterNameID, previewPosterID, backgroundID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePictures", ctx, id, posterNameID, previewPosterID, backgroundID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePictures indicates an expected call of UpdatePictures.
func (mr *MockMoviesRepositoryMockRecorder) UpdatePictures(ctx, id, posterNameID, previewPosterID, backgroundID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePictures", reflect.TypeOf((*MockMoviesRepository)(nil).UpdatePictures), ctx, id, posterNameID, previewPosterID, backgroundID)
}

// MockCountriesRepository is a mock of CountriesRepository interface.
type MockCountriesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCountriesRepositoryMockRecorder
}

// MockCountriesRepositoryMockRecorder is the mock recorder for MockCountriesRepository.
type MockCountriesRepositoryMockRecorder struct {
	mock *MockCountriesRepository
}

// NewMockCountriesRepository creates a new mock instance.
func NewMockCountriesRepository(ctrl *gomock.Controller) *MockCountriesRepository {
	mock := &MockCountriesRepository{ctrl: ctrl}
	mock.recorder = &MockCountriesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCountriesRepository) EXPECT() *MockCountriesRepositoryMockRecorder {
	return m.recorder
}

// CreateCountry mocks base method.
func (m *MockCountriesRepository) CreateCountry(ctx context.Context, name string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCountry", ctx, name)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCountry indicates an expected call of CreateCountry.
func (mr *MockCountriesRepositoryMockRecorder) CreateCountry(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCountry", reflect.TypeOf((*MockCountriesRepository)(nil).CreateCountry), ctx, name)
}

// DeleteCountry mocks base method.
func (m *MockCountriesRepository) DeleteCountry(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCountry", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCountry indicates an expected call of DeleteCountry.
func (mr *MockCountriesRepositoryMockRecorder) DeleteCountry(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCountry", reflect.TypeOf((*MockCountriesRepository)(nil).DeleteCountry), ctx, id)
}

// GetCountries mocks base method.
func (m *MockCountriesRepository) GetCountries(ctx context.Context) ([]repository.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountries", ctx)
	ret0, _ := ret[0].([]repository.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountries indicates an expected call of GetCountries.
func (mr *MockCountriesRepositoryMockRecorder) GetCountries(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountries", reflect.TypeOf((*MockCountriesRepository)(nil).GetCountries), ctx)
}

// GetCountriesForMovie mocks base method.
func (m *MockCountriesRepository) GetCountriesForMovie(ctx context.Context, id int32) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountriesForMovie", ctx, id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountriesForMovie indicates an expected call of GetCountriesForMovie.
func (mr *MockCountriesRepositoryMockRecorder) GetCountriesForMovie(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountriesForMovie", reflect.TypeOf((*MockCountriesRepository)(nil).GetCountriesForMovie), ctx, id)
}

// GetCountriesForMovies mocks base method.
func (m *MockCountriesRepository) GetCountriesForMovies(ctx context.Context, ids []int32) (map[int32][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountriesForMovies", ctx, ids)
	ret0, _ := ret[0].(map[int32][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountriesForMovies indicates an expected call of GetCountriesForMovies.
func (mr *MockCountriesRepositoryMockRecorder) GetCountriesForMovies(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountriesForMovies", reflect.TypeOf((*MockCountriesRepository)(nil).GetCountriesForMovies), ctx, ids)
}

// GetCountry mocks base method.
func (m *MockCountriesRepository) GetCountry(ctx context.Context, id int32) (repository.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountry", ctx, id)
	ret0, _ := ret[0].(repository.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountry indicates an expected call of GetCountry.
func (mr *MockCountriesRepositoryMockRecorder) GetCountry(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountry", reflect.TypeOf((*MockCountriesRepository)(nil).GetCountry), ctx, id)
}

// GetCountryByName mocks base method.
func (m *MockCountriesRepository) GetCountryByName(ctx context.Context, name string) (repository.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountryByName", ctx, name)
	ret0, _ := ret[0].(repository.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountryByName indicates an expected call of GetCountryByName.
func (mr *MockCountriesRepositoryMockRecorder) GetCountryByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountryByName", reflect.TypeOf((*MockCountriesRepository)(nil).GetCountryByName), ctx, name)
}

// IsCountriesExists mocks base method.
func (m *MockCountriesRepository) IsCountriesExists(ctx context.Context, ids []int32) (bool, []int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCountriesExists", ctx, ids)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsCountriesExists indicates an expected call of IsCountriesExists.
func (mr *MockCountriesRepositoryMockRecorder) IsCountriesExists(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCountriesExists", reflect.TypeOf((*MockCountriesRepository)(nil).IsCountriesExists), ctx, ids)
}

// IsCountryAlreadyExists mocks base method.
func (m *MockCountriesRepository) IsCountryAlreadyExists(ctx context.Context, name string) (bool, int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCountryAlreadyExists", ctx, name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsCountryAlreadyExists indicates an expected call of IsCountryAlreadyExists.
func (mr *MockCountriesRepositoryMockRecorder) IsCountryAlreadyExists(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCountryAlreadyExists", reflect.TypeOf((*MockCountriesRepository)(nil).IsCountryAlreadyExists), ctx, name)
}

// IsCountryExist mocks base method.
func (m *MockCountriesRepository) IsCountryExist(ctx context.Context, id int32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCountryExist", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsCountryExist indicates an expected call of IsCountryExist.
func (mr *MockCountriesRepositoryMockRecorder) IsCountryExist(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCountryExist", reflect.TypeOf((*MockCountriesRepository)(nil).IsCountryExist), ctx, id)
}

// UpdateCountry mocks base method.
func (m *MockCountriesRepository) UpdateCountry(ctx context.Context, countryName string, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCountry", ctx, countryName, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCountry indicates an expected call of UpdateCountry.
func (mr *MockCountriesRepositoryMockRecorder) UpdateCountry(ctx, countryName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCountry", reflect.TypeOf((*MockCountriesRepository)(nil).UpdateCountry), ctx, countryName, id)
}

// MockGenresRepository is a mock of GenresRepository interface.
type MockGenresRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGenresRepositoryMockRecorder
}

// MockGenresRepositoryMockRecorder is the mock recorder for MockGenresRepository.
type MockGenresRepositoryMockRecorder struct {
	mock *MockGenresRepository
}

// NewMockGenresRepository creates a new mock instance.
func NewMockGenresRepository(ctrl *gomock.Controller) *MockGenresRepository {
	mock := &MockGenresRepository{ctrl: ctrl}
	mock.recorder = &MockGenresRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenresRepository) EXPECT() *MockGenresRepositoryMockRecorder {
	return m.recorder
}

// CreateGenre mocks base method.
func (m *MockGenresRepository) CreateGenre(ctx context.Context, name string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGenre", ctx, name)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGenre indicates an expected call of CreateGenre.
func (mr *MockGenresRepositoryMockRecorder) CreateGenre(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGenre", reflect.TypeOf((*MockGenresRepository)(nil).CreateGenre), ctx, name)
}

// DeleteGenre mocks base method.
func (m *MockGenresRepository) DeleteGenre(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGenre", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGenre indicates an expected call of DeleteGenre.
func (mr *MockGenresRepositoryMockRecorder) DeleteGenre(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGenre", reflect.TypeOf((*MockGenresRepository)(nil).DeleteGenre), ctx, id)
}

// GetGenre mocks base method.
func (m *MockGenresRepository) GetGenre(ctx context.Context, id int32) (repository.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenre", ctx, id)
	ret0, _ := ret[0].(repository.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenre indicates an expected call of GetGenre.
func (mr *MockGenresRepositoryMockRecorder) GetGenre(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenre", reflect.TypeOf((*MockGenresRepository)(nil).GetGenre), ctx, id)
}

// GetGenreByName mocks base method.
func (m *MockGenresRepository) GetGenreByName(ctx context.Context, name string) (repository.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenreByName", ctx, name)
	ret0, _ := ret[0].(repository.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenreByName indicates an expected call of GetGenreByName.
func (mr *MockGenresRepositoryMockRecorder) GetGenreByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenreByName", reflect.TypeOf((*MockGenresRepository)(nil).GetGenreByName), ctx, name)
}

// GetGenres mocks base method.
func (m *MockGenresRepository) GetGenres(ctx context.Context) ([]repository.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres", ctx)
	ret0, _ := ret[0].([]repository.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockGenresRepositoryMockRecorder) GetGenres(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockGenresRepository)(nil).GetGenres), ctx)
}

// GetGenresForMovie mocks base method.
func (m *MockGenresRepository) GetGenresForMovie(ctx context.Context, id int32) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenresForMovie", ctx, id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenresForMovie indicates an expected call of GetGenresForMovie.
func (mr *MockGenresRepositoryMockRecorder) GetGenresForMovie(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenresForMovie", reflect.TypeOf((*MockGenresRepository)(nil).GetGenresForMovie), ctx, id)
}

// GetGenresForMovies mocks base method.
func (m *MockGenresRepository) GetGenresForMovies(ctx context.Context, ids []int32) (map[int32][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenresForMovies", ctx, ids)
	ret0, _ := ret[0].(map[int32][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenresForMovies indicates an expected call of GetGenresForMovies.
func (mr *MockGenresRepositoryMockRecorder) GetGenresForMovies(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenresForMovies", reflect.TypeOf((*MockGenresRepository)(nil).GetGenresForMovies), ctx, ids)
}

// IsGenreAlreadyExists mocks base method.
func (m *MockGenresRepository) IsGenreAlreadyExists(ctx context.Context, nameRU string) (bool, int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsGenreAlreadyExists", ctx, nameRU)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsGenreAlreadyExists indicates an expected call of IsGenreAlreadyExists.
func (mr *MockGenresRepositoryMockRecorder) IsGenreAlreadyExists(ctx, nameRU interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsGenreAlreadyExists", reflect.TypeOf((*MockGenresRepository)(nil).IsGenreAlreadyExists), ctx, nameRU)
}

// IsGenreExist mocks base method.
func (m *MockGenresRepository) IsGenreExist(ctx context.Context, id int32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsGenreExist", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsGenreExist indicates an expected call of IsGenreExist.
func (mr *MockGenresRepositoryMockRecorder) IsGenreExist(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsGenreExist", reflect.TypeOf((*MockGenresRepository)(nil).IsGenreExist), ctx, id)
}

// IsGenresExists mocks base method.
func (m *MockGenresRepository) IsGenresExists(ctx context.Context, ids []int32) (bool, []int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsGenresExists", ctx, ids)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsGenresExists indicates an expected call of IsGenresExists.
func (mr *MockGenresRepositoryMockRecorder) IsGenresExists(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsGenresExists", reflect.TypeOf((*MockGenresRepository)(nil).IsGenresExists), ctx, ids)
}

// UpdateGenre mocks base method.
func (m *MockGenresRepository) UpdateGenre(ctx context.Context, name string, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGenre", ctx, name, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGenre indicates an expected call of UpdateGenre.
func (mr *MockGenresRepositoryMockRecorder) UpdateGenre(ctx, name, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGenre", reflect.TypeOf((*MockGenresRepository)(nil).UpdateGenre), ctx, name, id)
}
