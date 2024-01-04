package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("entity not found")
var ErrInvalidArgument = errors.New("invalid input data")

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

type Movie struct {
	ID                  int32          `db:"id"`
	TitleRU             string         `db:"title_ru"`
	TitleEN             sql.NullString `db:"title_en"`
	PreviewPosterID     sql.NullString `db:"preview_poster_picture_id"`
	Description         string         `db:"description"`
	ShortDescription    string         `db:"short_description"`
	Duration            int32          `db:"duration"`
	PosterID            sql.NullString `db:"poster_picture_id"`
	BackgroundPictureID sql.NullString `db:"background_picture_id"`
	ReleaseYear         int32          `db:"release_year"`
	AgeRating           string         `db:"age_rating"`
}

type MoviesFilter struct {
	MoviesIDs    string
	GenresIDs    string
	CountriesIDs string
	Title        string
	AgeRating    string
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type CreateMovieParam struct {
	TitleRU             string  `db:"title_ru"`
	TitleEN             string  `db:"title_en"`
	Description         string  `db:"description"`
	Genres              []int32 `db:"genres"`
	Duration            int32   `db:"duration"`
	PosterID            string  `db:"poster_picture_id"`
	PreviewPosterID     string  `db:"preview_poster_picture_id"`
	BackgroundPictureID string  `db:"background_picture_id"`
	ShortDescription    string  `db:"short_description"`
	CountriesIDs        []int32 `db:"countries"`
	ReleaseYear         int32   `db:"release_year"`
	AgeRating           int32   `db:"age_rating_id"`
}

type Country struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type IsMovieAlreadyExistsParam struct {
	TitleRU          string `db:"title_ru"`
	TitleEN          string `db:"title_en"`
	Description      string `db:"description"`
	Duration         int32  `db:"duration"`
	ShortDescription string `db:"short_description"`
	ReleaseYear      int32  `db:"release_year"`
	AgeRating        int32  `db:"age_rating_id"`
}

type UpdateMovieParam struct {
	TitleRU          string `db:"title_ru"`
	TitleEN          string `db:"title_en"`
	Description      string `db:"description"`
	Duration         int32  `db:"duration"`
	ShortDescription string `db:"short_description"`
	ReleaseYear      int32  `db:"release_year"`
	AgeRating        int32  `db:"age_rating_id"`
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go
type MoviesRepository interface {
	GetMovie(ctx context.Context, movieID int32) (Movie, error)
	GetMovies(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]Movie, error)
	CreateMovie(ctx context.Context, param CreateMovieParam) (int32, error)
	IsMovieAlreadyExists(ctx context.Context, existParam IsMovieAlreadyExistsParam) (bool, []int32, error)
	UpdatePictures(ctx context.Context, id int32, posterNameID, previewPosterID, backgroundID string) error
	UpdateMovie(ctx context.Context, id int32, param UpdateMovieParam,
		genres, countries []int32, updateGenres, updateCountries bool) error
	GetPicturesIds(ctx context.Context, id int32) (poster, preview, background string, err error)

	IsMovieExists(ctx context.Context, id int32) (bool, error)
	DeleteMovie(ctx context.Context, id int32) error
	GetAgeRatings(ctx context.Context) ([]string, error)
	GetAgeRating(ctx context.Context, name string) (int32, error)
	IsAgeRatingAlreadyExists(ctx context.Context, name string) (bool, error)
	CreateAgeRating(ctx context.Context, name string) error
	DeleteAgeRating(ctx context.Context, name string) error
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go
type CountriesRepository interface {
	GetCountry(ctx context.Context, id int32) (Country, error)
	GetCountryByName(ctx context.Context, name string) (Country, error)

	GetCountries(ctx context.Context) ([]Country, error)
	GetCountriesForMovies(ctx context.Context, ids []int32) (map[int32][]string, error)
	GetCountriesForMovie(ctx context.Context, id int32) ([]string, error)

	CreateCountry(ctx context.Context, name string) (int32, error)
	DeleteCountry(ctx context.Context, id int32) error
	UpdateCountry(ctx context.Context, countryName string, id int32) error

	IsCountryExist(ctx context.Context, id int32) (bool, error)
	IsCountriesExists(ctx context.Context, ids []int32) (bool, []int32, error)
	IsCountryAlreadyExists(ctx context.Context, name string) (bool, int32, error)
}

type Genre struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go
type GenresRepository interface {
	GetGenre(ctx context.Context, id int32) (Genre, error)
	GetGenreByName(ctx context.Context, name string) (Genre, error)
	GetGenres(ctx context.Context) ([]Genre, error)

	GetGenresForMovies(ctx context.Context, ids []int32) (map[int32][]string, error)
	GetGenresForMovie(ctx context.Context, id int32) ([]string, error)

	CreateGenre(ctx context.Context, name string) (int32, error)
	DeleteGenre(ctx context.Context, id int32) error
	UpdateGenre(ctx context.Context, name string, id int32) error

	IsGenreExist(ctx context.Context, id int32) (bool, error)
	IsGenreAlreadyExists(ctx context.Context, nameRU string) (bool, int32, error)
	IsGenresExists(ctx context.Context, ids []int32) (bool, []int32, error)
}
