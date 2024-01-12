package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Falokut/admin_movies_service/internal/events"
	"github.com/Falokut/admin_movies_service/internal/repository"
	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	image_processing_service "github.com/Falokut/image_processing_service/pkg/image_processing_service/v1/protos"
	"golang.org/x/exp/maps"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go
type MoviesRepository interface {
	GetMovie(ctx context.Context, in *admin_movies_service.GetMovieRequest) (*admin_movies_service.Movie, error)
	GetMovies(ctx context.Context, in *admin_movies_service.GetMoviesRequest) (*admin_movies_service.Movies, error)
	GetAgeRatings(ctx context.Context) (*admin_movies_service.AgeRatings, error)
	CreateAgeRating(ctx context.Context, in *admin_movies_service.CreateAgeRatingRequest) error
	DeleteAgeRating(ctx context.Context, in *admin_movies_service.DeleteAgeRatingRequest) error
	CreateMovie(ctx context.Context,
		in *admin_movies_service.CreateMovieRequest) (*admin_movies_service.CreateMovieResponce, error)
	IsMovieExists(ctx context.Context,
		in *admin_movies_service.IsMovieExistsRequest) (*admin_movies_service.IsMovieExistsResponce, error)
	UpdateMoviePictures(ctx context.Context, in *admin_movies_service.UpdateMoviePicturesRequest) error
	UpdateMovie(ctx context.Context, in *admin_movies_service.UpdateMovieRequest) error
	DeleteMovie(ctx context.Context, in *admin_movies_service.DeleteMovieRequest) error
}

//go:generate mockgen -source=service.go -destination=mocks/service.go
type ImagesService interface {
	GetPictureURL(pictureID, category string) string
	UploadPicture(ctx context.Context, picture UploadPictureParam) (string, error)
	ReplacePicture(ctx context.Context, picture ReplacePicturesParam, createIfNotExist bool) (string, error)

	UploadPictures(ctx context.Context, pictures map[string]UploadPictureParam) (map[string]string, error)
	ReplacePictures(ctx context.Context, pictures map[string]ReplacePicturesParam) (map[string]string, error)

	DeletePicture(ctx context.Context, category, pictureID string) error
}

//go:generate mockgen -source=service.go -destination=mocks/service.go
type ExistanceChecker interface {
	CheckExistance(ctx context.Context, countriesIDs, genresIDs []int32) error
}

type PictureConfig struct {
	ValidateImage bool
	CheckImage    struct {
		MaxImageWidth  int32
		MaxImageHeight int32
		MinImageWidth  int32
		MinImageHeight int32
		AllowedTypes   []string
	}
	ResizeImage           bool
	ImageProcessingConfig struct {
		ImageHeight       int32
		ImageWidth        int32
		ImageResizeMethod image_processing_service.ResampleFilter
	}
	Category string
}

type MoviesService struct {
	admin_movies_service.UnimplementedMoviesServiceV1Server
	logger        *logrus.Logger
	moviesRepo    MoviesRepository
	genresRepo    repository.GenresRepository
	countriesRepo repository.CountriesRepository
	errorHandler  errorHandler
	moviesEvents  events.MoviesEventsMQ
}

func NewMoviesService(logger *logrus.Logger, moviesRepo MoviesRepository,
	countriesRepo repository.CountriesRepository,
	genresRepo repository.GenresRepository, moviesEvents events.MoviesEventsMQ) *MoviesService {
	errorHandler := newErrorHandler(logger)
	return &MoviesService{
		logger:        logger,
		moviesRepo:    moviesRepo,
		genresRepo:    genresRepo,
		countriesRepo: countriesRepo,
		errorHandler:  errorHandler,
		moviesEvents:  moviesEvents,
	}
}

func (s *MoviesService) GetMovie(ctx context.Context, in *admin_movies_service.GetMovieRequest) (*admin_movies_service.Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMovie")
	defer span.Finish()

	movie, err := s.moviesRepo.GetMovie(ctx, in)
	switch err {
	case ErrNotFound:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	default:
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}
		span.SetTag("grpc.status", codes.OK)
		return movie, nil
	}
}

func (s *MoviesService) GetAgeRatings(ctx context.Context, in *emptypb.Empty) (*admin_movies_service.AgeRatings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetAgeRatings")
	defer span.Finish()

	ratings, err := s.moviesRepo.GetAgeRatings(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return ratings, nil
}

func (s *MoviesService) CreateAgeRating(ctx context.Context, in *admin_movies_service.CreateAgeRatingRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.CreateAgeRating")
	defer span.Finish()
	err := s.moviesRepo.CreateAgeRating(ctx, in)
	switch err {
	case ErrAlreadyExists:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrAlreadyExists, "age ratings already exists")
	default:
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}

		span.SetTag("grpc.status", codes.OK)
		return &emptypb.Empty{}, nil
	}
}

func (s *MoviesService) DeleteAgeRating(ctx context.Context, in *admin_movies_service.DeleteAgeRatingRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.DeleteAgeRating")
	defer span.Finish()

	err := s.moviesRepo.DeleteAgeRating(ctx, in)
	switch err {
	case ErrNotFound:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, fmt.Sprintf("age rating with %s name not found", in.AgeRatingName))
	default:
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}

		span.SetTag("grpc.status", codes.OK)
		return &emptypb.Empty{}, nil
	}
}

func (s *MoviesService) GetMovies(ctx context.Context, in *admin_movies_service.GetMoviesRequest) (*admin_movies_service.Movies, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMovies")
	defer span.Finish()

	movies, err := s.moviesRepo.GetMovies(ctx, in)
	switch err {
	case ErrNotFound:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, err, "movies not found")
	case ErrInvalidFilter, ErrInvalidArgument:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, err, "")
	default:
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}
		span.SetTag("grpc.status", codes.OK)
		return movies, nil
	}
}

func (s *MoviesService) CreateMovie(ctx context.Context,
	in *admin_movies_service.CreateMovieRequest) (*admin_movies_service.CreateMovieResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.CreateMovie")
	defer span.Finish()
	if err := validateCreateRequest(in); err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}

	res, err := s.moviesRepo.CreateMovie(ctx, in)
	if errors.Is(err, ErrInvalidArgument) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}
	if errors.Is(err, ErrAlreadyExists) {
		return nil, s.errorHandler.createExtendedErrorResponceWithSpan(span, ErrAlreadyExists, "", err.Error())
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return res, nil

}

func formatSlice[T any](nums []T) string {
	var str = make([]string, 0, len(nums))
	for _, num := range nums {
		str = append(str, fmt.Sprint(num))
	}
	return strings.Join(str, ",")
}

func (s *MoviesService) DeleteMovie(ctx context.Context, in *admin_movies_service.DeleteMovieRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.DeleteMovie")
	defer span.Finish()

	err := s.moviesRepo.DeleteMovie(ctx, in)
	switch err {
	case ErrNotFound:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "movie with this id not found")
	default:
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}
		s.moviesEvents.MovieDeleted(ctx, in.MovieID)
		span.SetTag("grpc.status", codes.OK)
		return &emptypb.Empty{}, nil
	}
}

func (s *MoviesService) IsMovieExists(ctx context.Context,
	in *admin_movies_service.IsMovieExistsRequest) (*admin_movies_service.IsMovieExistsResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.IsMovieExists")
	defer span.Finish()

	exists, err := s.moviesRepo.IsMovieExists(ctx, in)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return exists, nil
}

func (s *MoviesService) UpdateMovie(ctx context.Context,
	in *admin_movies_service.UpdateMovieRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.UpdateMovie")
	defer span.Finish()

	err := s.moviesRepo.UpdateMovie(ctx, in)
	switch err {
	case ErrNotFound:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, fmt.Sprintf("movie with id: %d not found", in.MovieID))
	default:
		if status.Code(err) == codes.InvalidArgument {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, err, "")
		}
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}

		span.SetTag("grpc.status", codes.OK)
		return &emptypb.Empty{}, nil
	}
}

func (s *MoviesService) UpdateMoviePictures(ctx context.Context,
	in *admin_movies_service.UpdateMoviePicturesRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.UpdateMoviePictures")
	defer span.Finish()

	err := s.moviesRepo.UpdateMoviePictures(ctx, in)
	switch err {
	case ErrNotFound:
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, fmt.Sprintf("movie with id: %d not found", in.MovieID))
	default:
		if status.Code(err) == codes.InvalidArgument {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
		}
		if err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
		}

		span.SetTag("grpc.status", codes.OK)
		return &emptypb.Empty{}, nil
	}
}

func cleanSlice(nums []int32) []int32 {
	m := make(map[int32]struct{})
	for _, num := range nums {
		_, ok := m[num]
		if !ok {
			m[num] = struct{}{}
		}
	}
	return maps.Keys(m)
}

func GetAgeRatingsFilter(ageRating string) string {
	ageRating = ReplaceAllDoubleQuates(strings.ReplaceAll(ageRating, " ", ""))
	str := strings.Split(ageRating, ",")
	for i := 0; i < len(str); i++ {
		if num, err := strconv.Atoi(str[i]); err == nil {
			str[i] = fmt.Sprintf("%d+", num)
		}
	}
	return strings.Join(str, ",")
}

func ReplaceAllDoubleQuates(s string) string {
	return strings.ReplaceAll(s, `"`, "")
}

// Countries section

func (s *MoviesService) GetCountry(ctx context.Context,
	in *admin_movies_service.GetCountryRequest) (*admin_movies_service.Country, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetCountry")
	defer span.Finish()

	if in.CountryId < 0 {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "id must be positive")
	}

	country, err := s.countriesRepo.GetCountry(ctx, in.CountryId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return convertCountryToProto(country), nil
}

func (s *MoviesService) GetCountryByName(ctx context.Context, in *admin_movies_service.GetCountryByNameRequest) (*admin_movies_service.Country, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetCountryByName")
	defer span.Finish()

	if in.Name == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "country name mustn't be empty")
	}

	country, err := s.countriesRepo.GetCountryByName(ctx, in.Name)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return convertCountryToProto(country), nil
}

func (s *MoviesService) CreateCountry(ctx context.Context,
	in *admin_movies_service.CreateCountryRequest) (*admin_movies_service.CreateCountryResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.CreateCountry")
	defer span.Finish()
	if in.Name == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "fields mustn't be empty")
	}

	exists, id, err := s.countriesRepo.IsCountryAlreadyExists(ctx, in.Name)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	} else if exists {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrAlreadyExists,
			fmt.Sprintf("check country with %d id", id))
	}

	id, err = s.countriesRepo.CreateCountry(ctx, in.Name)
	if errors.Is(err, repository.ErrInvalidArgument) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "fields mustn't be empty")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &admin_movies_service.CreateCountryResponce{CountryID: id}, nil
}

func (s *MoviesService) GetCountries(ctx context.Context, in *emptypb.Empty) (*admin_movies_service.Countries, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetCountries")
	defer span.Finish()

	var countries []repository.Country
	var err error

	countries, err = s.countriesRepo.GetCountries(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}
	if len(countries) == 0 {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "no countries was found")
	}

	span.SetTag("grpc.status", codes.OK)
	return convertCountriesToProto(countries), nil
}

func convertStringsSlice(str []string) []int32 {
	var nums = make([]int32, 0, len(str))
	for _, s := range str {
		num, err := strconv.Atoi(s)
		if err == nil {
			nums = append(nums, int32(num))
		}
	}
	return nums
}

func (s *MoviesService) UpdateCountry(ctx context.Context, in *admin_movies_service.UpdateCountryRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.UpdateCountry")
	defer span.Finish()

	exist, err := s.countriesRepo.IsCountryExist(ctx, in.CountryID)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	} else if !exist {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}

	err = s.countriesRepo.UpdateCountry(ctx, in.Name, in.CountryID)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *MoviesService) DeleteCountry(ctx context.Context, in *admin_movies_service.DeleteCountryRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.DeleteCountry")
	defer span.Finish()

	err := s.countriesRepo.DeleteCountry(ctx, in.CountryId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "country with this id not found")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *MoviesService) IsCountriesExists(ctx context.Context,
	in *admin_movies_service.IsCountriesExistsRequest) (*admin_movies_service.ExistsResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.IsCountriesExists")
	defer span.Finish()

	in.CountriesIDs = strings.ReplaceAll(in.CountriesIDs, `"`, "")
	if in.CountriesIDs == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "countries_ids mustn't be empty")
	} else if err := checkParam(in.CountriesIDs); err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}

	ids := convertStringsSlice(strings.Split(in.CountriesIDs, ","))
	exists, notFoundIDs, err := s.countriesRepo.IsCountriesExists(ctx, ids)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &admin_movies_service.ExistsResponce{Exists: exists, NotExistIDs: notFoundIDs}, nil
}

func convertCountriesToProto(countries []repository.Country) *admin_movies_service.Countries {
	res := &admin_movies_service.Countries{}
	res.Countries = make([]*admin_movies_service.Country, 0, len(countries))

	for _, country := range countries {
		res.Countries = append(res.Countries, &admin_movies_service.Country{
			CountryID: country.ID,
			Name:      country.Name,
		})
	}

	return res
}

func convertCountryToProto(country repository.Country) *admin_movies_service.Country {
	return &admin_movies_service.Country{
		CountryID: country.ID,
		Name:      country.Name,
	}
}

func (s *MoviesService) GetGenre(ctx context.Context, in *admin_movies_service.GetGenreRequest) (*admin_movies_service.Genre, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetGenre")
	defer span.Finish()

	if in.GenreId < 0 {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "id must be positive")
	}

	country, err := s.genresRepo.GetGenre(ctx, in.GenreId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return convertGenreToProto(country), nil
}

func (s *MoviesService) GetGenreByName(ctx context.Context, in *admin_movies_service.GetGenreByNameRequest) (*admin_movies_service.Genre, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetGenreByName")
	defer span.Finish()

	if in.Name == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "genre name mustn't be empty")
	}

	genre, err := s.genresRepo.GetGenreByName(ctx, in.Name)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return convertGenreToProto(genre), nil
}

func (s *MoviesService) CreateGenre(ctx context.Context,
	in *admin_movies_service.CreateGenreRequest) (*admin_movies_service.CreateGenreResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.CreateGenre")
	defer span.Finish()
	if in.Name == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "fields mustn't be empty")
	}

	exists, id, err := s.genresRepo.IsGenreAlreadyExists(ctx, in.Name)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	} else if exists {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrAlreadyExists,
			fmt.Sprintf("check genre with %d id", id))
	}

	id, err = s.genresRepo.CreateGenre(ctx, in.Name)
	if errors.Is(err, repository.ErrInvalidArgument) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "fields mustn't be empty")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &admin_movies_service.CreateGenreResponce{GenreID: id}, nil
}

func (s *MoviesService) GetGenres(ctx context.Context, in *emptypb.Empty) (*admin_movies_service.Genres, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetGenres")
	defer span.Finish()

	genres, err := s.genresRepo.GetGenres(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}
	if len(genres) == 0 {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "no genres was found")
	}

	span.SetTag("grpc.status", codes.OK)
	return convertGenresToProto(genres), nil
}

func (s *MoviesService) UpdateGenre(ctx context.Context, in *admin_movies_service.UpdateGenreRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.UpdateGenre")
	defer span.Finish()

	exist, err := s.genresRepo.IsGenreExist(ctx, in.GenreID)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	} else if !exist {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}

	err = s.genresRepo.UpdateGenre(ctx, in.Name, in.GenreID)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *MoviesService) DeleteGenre(ctx context.Context, in *admin_movies_service.DeleteGenreRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.DeleteGenre")
	defer span.Finish()

	err := s.genresRepo.DeleteGenre(ctx, in.GenreId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "genre with this id not found")
	} else if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *MoviesService) IsGenresExists(ctx context.Context,
	in *admin_movies_service.IsGenresExistsRequest) (*admin_movies_service.ExistsResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.IsGenresExists")
	defer span.Finish()

	in.GenresIDs = strings.ReplaceAll(in.GenresIDs, `"`, "")
	if in.GenresIDs == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "genres_ids mustn't be empty")
	} else if err := checkParam(in.GenresIDs); err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}

	ids := convertStringsSlice(strings.Split(in.GenresIDs, ","))
	exists, notFoundIDs, err := s.genresRepo.IsGenresExists(ctx, ids)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &admin_movies_service.ExistsResponce{Exists: exists, NotExistIDs: notFoundIDs}, nil
}

func convertGenresToProto(genres []repository.Genre) *admin_movies_service.Genres {
	res := &admin_movies_service.Genres{}
	res.Genres = make([]*admin_movies_service.Genre, 0, len(genres))

	for _, genre := range genres {
		res.Genres = append(res.Genres, &admin_movies_service.Genre{
			GenreID: genre.ID,
			Name:    genre.Name,
		})
	}

	return res
}

func convertGenreToProto(city repository.Genre) *admin_movies_service.Genre {
	return &admin_movies_service.Genre{
		GenreID: city.ID,
		Name:    city.Name,
	}
}
