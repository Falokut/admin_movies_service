package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Falokut/admin_movies_service/internal/repository"
	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	image_processing_service "github.com/Falokut/image_processing_service/pkg/image_processing_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type moviesRepositoryWrapper struct {
	countriesRepo    repository.CountriesRepository
	genresRepo       repository.GenresRepository
	moviesRepo       repository.MoviesRepository
	existanceChecker ExistanceChecker

	logger *logrus.Logger

	imagesService ImagesService

	postersConfig        PictureConfig
	previewPostersConfig PictureConfig
	backgroundConfig     PictureConfig
}

func NewMoviesRepositoryWrapper(countriesRepo repository.CountriesRepository,
	genresRepo repository.GenresRepository,
	moviesRepo repository.MoviesRepository,
	imagesService ImagesService, postersConfig, previewPostersConfig, backgroundConfig PictureConfig,
	existanceChecker ExistanceChecker,
	logger *logrus.Logger) *moviesRepositoryWrapper {

	return &moviesRepositoryWrapper{countriesRepo: countriesRepo,
		genresRepo:           genresRepo,
		moviesRepo:           moviesRepo,
		logger:               logger,
		imagesService:        imagesService,
		postersConfig:        postersConfig,
		previewPostersConfig: previewPostersConfig,
		backgroundConfig:     backgroundConfig,
		existanceChecker:     existanceChecker,
	}
}

func (r *moviesRepositoryWrapper) GetMovie(ctx context.Context, in *admin_movies_service.GetMovieRequest) (*admin_movies_service.Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	movieCh := make(chan repository.Movie, 1)
	errCh := make(chan error, 1)
	genresAndCountriesCh := make(chan struct{ genres, countries []string }, 1)

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		select {
		case <-reqCtx.Done():
			errCh <- reqCtx.Err()
			return
		default:
			mov, err := r.moviesRepo.GetMovie(reqCtx, in.MovieID)
			if err != nil {
				errCh <- err
				return
			}
			movieCh <- mov
		}

	}()
	go func() {
		select {
		case <-reqCtx.Done():
			errCh <- reqCtx.Err()
			return
		default:
			genres, countries, err := r.getGenresAndCountriesForMovie(reqCtx, in.MovieID)
			if err != nil {
				errCh <- err
				return
			}
			genresAndCountriesCh <- struct {
				genres    []string
				countries []string
			}{genres: genres, countries: countries}
		}
	}()

	var movie admin_movies_service.Movie
	var movieDone, genresAndCountriesDone bool
	for !genresAndCountriesDone || !movieDone {
		select {
		case <-reqCtx.Done():
			return nil, reqCtx.Err()
		case str := <-genresAndCountriesCh:
			if genresAndCountriesDone {
				break
			}
			movie.Genres = str.genres
			movie.Countries = str.countries
			genresAndCountriesDone = true
			r.logger.Debug("genres and countries done")
		case res := <-movieCh:
			if movieDone {
				break
			}

			movie.TitleRU = res.TitleRU
			movie.TitleEN = res.TitleEN.String
			movie.Description = res.Description
			movie.ShortDescription = res.ShortDescription
			movie.Duration = res.Duration
			movie.PosterURL = r.imagesService.GetPictureURL(res.PosterID.String, r.postersConfig.Category)
			movie.PreviewPosterURL = r.imagesService.GetPictureURL(res.PreviewPosterID.String, r.previewPostersConfig.Category)
			movie.BackgroundURL = r.imagesService.GetPictureURL(res.BackgroundPictureID.String, r.backgroundConfig.Category)
			movie.ReleaseYear = res.ReleaseYear
			movie.AgeRating = res.AgeRating
			movieDone = true
			r.logger.Debug("movies done")
		case err = <-errCh:
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return nil, ErrNotFound
				}
				return nil, err
			}
		}
	}
	return &movie, nil
}

func (r *moviesRepositoryWrapper) getGenresAndCountriesForMovie(ctx context.Context, id int32) (genres, countries []string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.getGenresAndCountriesForMovie")
	defer span.Finish()
	defer span.SetTag("error", err != nil)

	genresCh, countriesCh := make(chan []string, 1), make(chan []string, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := r.genresRepo.GetGenresForMovie(ctx, id)
		if err != nil {
			errCh <- err
			return
		}
		genresCh <- res
	}()
	go func() {
		res, err := r.countriesRepo.GetCountriesForMovie(ctx, id)
		if err != nil {
			errCh <- err
			return
		}
		countriesCh <- res
	}()

	var countriesDone, genresDone bool
	for !genresDone || !countriesDone {
		select {
		case <-ctx.Done():
			return []string{}, []string{}, ctx.Err()
		case g := <-genresCh:
			if genresDone {
				break
			}
			genres = g
			genresDone = true
		case c := <-countriesCh:
			if countriesDone {
				break
			}
			countries = c
			countriesDone = true
		case err = <-errCh:
			return
		}
	}

	return
}

func (r *moviesRepositoryWrapper) getGenresAndCountriesForMovies(ctx context.Context, ids []int32) (genres, countries map[int32][]string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.getGenresAndCountriesForMovie")
	defer span.Finish()

	genresCh, countriesCh := make(chan map[int32][]string, 1),
		make(chan map[int32][]string, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := r.genresRepo.GetGenresForMovies(ctx, ids)
		if err != nil {
			errCh <- err
			return
		}
		genresCh <- res
	}()
	go func() {
		res, err := r.countriesRepo.GetCountriesForMovies(ctx, ids)
		if err != nil {
			errCh <- err
			return
		}
		countriesCh <- res
	}()

	var countriesDone, genresDone bool
	for !genresDone || !countriesDone {
		select {
		case <-ctx.Done():
			return map[int32][]string{}, map[int32][]string{}, ctx.Err()
		case g := <-genresCh:
			if genresDone {
				break
			}
			genres = g
			genresDone = true
		case c := <-countriesCh:
			if countriesDone {
				break
			}
			countries = c
			countriesDone = true
		case err = <-errCh:
			return
		}
	}
	return genres, countries, nil
}

func (r *moviesRepositoryWrapper) GetMovies(ctx context.Context,
	in *admin_movies_service.GetMoviesRequest) (*admin_movies_service.Movies, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.GetMovies")
	defer span.Finish()

	if err := validateFilter(in); err != nil {
		return nil, err
	}

	filter := repository.MoviesFilter{
		MoviesIDs:    ReplaceAllDoubleQuates(in.GetMoviesIDs()),
		GenresIDs:    ReplaceAllDoubleQuates(in.GetGenresIDs()),
		CountriesIDs: ReplaceAllDoubleQuates(in.GetCountriesIDs()),
		AgeRating:    GetAgeRatingsFilter(in.GetAgeRatings()),
		Title:        ReplaceAllDoubleQuates(in.GetTitle()),
	}

	if in.Limit == 0 {
		in.Limit = 10
	} else if in.Limit > 100 {
		in.Limit = 100
	}

	movies, err := r.moviesRepo.GetMovies(ctx, filter, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}
	if len(movies) == 0 {
		return nil, ErrNotFound
	}

	var ids = make([]int32, 0, len(movies))
	for _, p := range movies {
		ids = append(ids, p.ID)
	}

	genres, countries, err := r.getGenresAndCountriesForMovies(ctx, ids)
	if err != nil {
		return nil, err
	}

	res := make(map[int32]*admin_movies_service.Movie, len(movies))
	for i, movie := range movies {
		res[movie.ID] = &admin_movies_service.Movie{
			TitleRU:          movie.TitleRU,
			TitleEN:          movie.TitleEN.String,
			Description:      movie.Description,
			ShortDescription: movie.ShortDescription,
			BackgroundURL:    r.imagesService.GetPictureURL(movie.BackgroundPictureID.String, r.backgroundConfig.Category),
			Genres:           genres[movies[i].ID],
			Countries:        countries[movies[i].ID],
			Duration:         movie.Duration,
			PosterURL:        r.imagesService.GetPictureURL(movie.PosterID.String, r.postersConfig.Category),
			PreviewPosterURL: r.imagesService.GetPictureURL(movie.PreviewPosterID.String, r.previewPostersConfig.Category),
			ReleaseYear:      movie.ReleaseYear,
			AgeRating:        movie.AgeRating,
		}
	}
	return &admin_movies_service.Movies{Movies: res}, nil
}

func (s *moviesRepositoryWrapper) GetAgeRatings(ctx context.Context) (*admin_movies_service.AgeRatings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.GetAgeRatings")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	ratings, err := s.moviesRepo.GetAgeRatings(ctx)
	if err != nil {
		return nil, err
	}

	return &admin_movies_service.AgeRatings{Ratings: ratings}, nil
}

func (s *moviesRepositoryWrapper) CreateAgeRating(ctx context.Context, in *admin_movies_service.CreateAgeRatingRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.CreateAgeRating")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	exists, err := s.moviesRepo.IsAgeRatingAlreadyExists(ctx, in.AgeRatingName)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	err = s.moviesRepo.CreateAgeRating(ctx, in.AgeRatingName)
	if err != nil {
		return err
	}

	return nil
}

func (s *moviesRepositoryWrapper) DeleteAgeRating(ctx context.Context, in *admin_movies_service.DeleteAgeRatingRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.DeleteAgeRating")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	err = s.moviesRepo.DeleteAgeRating(ctx, in.AgeRatingName)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func convertProtoToCreateMovieParam(in *admin_movies_service.CreateMovieRequest,
	posterID, previewPosterID, backgroundID string, ageRating int32) repository.CreateMovieParam {
	return repository.CreateMovieParam{
		TitleRU:             in.TitleRU,
		TitleEN:             in.GetTitleEN(),
		Description:         in.Description,
		Genres:              in.GenresIDs,
		Duration:            in.Duration,
		PosterID:            posterID,
		PreviewPosterID:     previewPosterID,
		BackgroundPictureID: backgroundID,
		ShortDescription:    in.ShortDescription,
		CountriesIDs:        in.CountriesIDs,
		ReleaseYear:         in.ReleaseYear,
		AgeRating:           ageRating,
	}
}

func (s *moviesRepositoryWrapper) GetMovieDuration(ctx context.Context, id int32) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"moviesRepositoryWrapper.CreateMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	duration, err := s.moviesRepo.GetMovieDuration(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return 0, fmt.Errorf("movie with this id not found.%w", ErrNotFound)
	}
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (s *moviesRepositoryWrapper) CreateMovie(ctx context.Context,
	in *admin_movies_service.CreateMovieRequest) (*admin_movies_service.CreateMovieResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.CreateMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	s.logger.Info("validating data")
	ageRating, err := s.moviesRepo.GetAgeRating(ctx, in.AgeRating)
	if errors.Is(err, repository.ErrNotFound) {
		ext.LogError(span, err)
		return nil, fmt.Errorf("age rating not found. %w", ErrInvalidArgument)
	}
	if err != nil {
		ext.LogError(span, err)
		return nil, fmt.Errorf("can't get age rating id. %w ", err)
	}
	exists, moviesIds, err := s.moviesRepo.IsMovieAlreadyExists(ctx, convertProtoToSearchParam(in, ageRating))
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("finded movies with ids: %s,"+
			"If this list does not contain the id of the movies you "+
			"want to add, add more information about the movies. %w ", formatSlice(moviesIds), ErrAlreadyExists)
	}

	in.CountriesIDs = cleanSlice(in.CountriesIDs)
	in.GenresIDs = cleanSlice(in.GenresIDs)

	err = s.existanceChecker.CheckExistance(ctx, in.CountriesIDs, in.GenresIDs)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	s.logger.Info("data is valid")

	s.logger.Info("Uploading pictures")
	uploadParams := map[string]UploadPictureParam{
		"poster":        convertPictureConfigToUploadParam(in.Poster, s.postersConfig),
		"previewPoster": convertPictureConfigToUploadParam(in.PreviewPoster, s.previewPostersConfig),
		"background":    convertPictureConfigToUploadParam(in.Background, s.backgroundConfig),
	}

	ids, err := s.imagesService.UploadPictures(ctx, uploadParams)

	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		for key, id := range ids {
			err := s.imagesService.DeletePicture(ctx, uploadParams[key].Category, id)
			if err != nil {
				s.logger.Errorf("image with %s id not deleted err: %s", id, err.Error())
			}
		}

		return nil, err
	}

	posterID := ids["poster"]
	previewPosterID := ids["previewPoster"]
	backgroundID := ids["background"]

	s.logger.Info("Creating movie")
	id, err := s.moviesRepo.CreateMovie(ctx,
		convertProtoToCreateMovieParam(in, posterID, previewPosterID, backgroundID, ageRating))
	if err != nil {
		return nil, err
	}

	s.logger.Info("Movie created")
	return &admin_movies_service.CreateMovieResponce{MovieID: id}, nil
}

func convertProtoToSearchParam(in *admin_movies_service.CreateMovieRequest, ageRating int32) repository.IsMovieAlreadyExistsParam {
	return repository.IsMovieAlreadyExistsParam{
		TitleRU:          in.TitleRU,
		TitleEN:          in.GetTitleEN(),
		Description:      in.Description,
		Duration:         in.Duration,
		ShortDescription: in.ShortDescription,
		ReleaseYear:      in.ReleaseYear,
		AgeRating:        ageRating,
	}
}

func convertPictureConfigToUploadParam(img []byte, cfg PictureConfig) UploadPictureParam {
	return UploadPictureParam{
		Image:         img,
		ValidateImage: cfg.ValidateImage,
		CheckImage: CheckImageConfig{
			MaxImageWidth:  cfg.CheckImage.MaxImageWidth,
			MaxImageHeight: cfg.CheckImage.MaxImageHeight,
			MinImageWidth:  cfg.CheckImage.MinImageWidth,
			MinImageHeight: cfg.CheckImage.MinImageHeight,
		},
		ResizeImage: cfg.ResizeImage,
		ImageProcessingParam: struct {
			ImageHeight       int32
			ImageWidth        int32
			ImageResizeMethod image_processing_service.ResampleFilter
		}{
			ImageHeight:       cfg.ImageProcessingConfig.ImageHeight,
			ImageWidth:        cfg.ImageProcessingConfig.ImageWidth,
			ImageResizeMethod: cfg.ImageProcessingConfig.ImageResizeMethod,
		},
		Category: cfg.Category,
	}
}

func convertPictureConfigToReplaceParam(img []byte, cfg PictureConfig, pictureName string) ReplacePicturesParam {
	return ReplacePicturesParam{
		Image:         img,
		ValidateImage: cfg.ValidateImage,
		CheckImage: CheckImageConfig{
			MaxImageWidth:  cfg.CheckImage.MaxImageWidth,
			MaxImageHeight: cfg.CheckImage.MaxImageHeight,
			MinImageWidth:  cfg.CheckImage.MinImageWidth,
			MinImageHeight: cfg.CheckImage.MinImageHeight,
		},
		ResizeImage: cfg.ResizeImage,
		ImageProcessingParam: struct {
			ImageHeight       int32
			ImageWidth        int32
			ImageResizeMethod image_processing_service.ResampleFilter
		}{
			ImageHeight:       cfg.ImageProcessingConfig.ImageHeight,
			ImageWidth:        cfg.ImageProcessingConfig.ImageWidth,
			ImageResizeMethod: cfg.ImageProcessingConfig.ImageResizeMethod,
		},
		Category:  cfg.Category,
		ImageName: pictureName,
	}
}

func (r *moviesRepositoryWrapper) IsMovieExists(ctx context.Context, in *admin_movies_service.IsMovieExistsRequest) (*admin_movies_service.IsMovieExistsResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.IsMovieExists")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	exists, err := r.moviesRepo.IsMovieExists(ctx, in.MovieID)
	if err != nil {
		return nil, err
	}

	return &admin_movies_service.IsMovieExistsResponce{Exists: exists}, nil
}

func (r *moviesRepositoryWrapper) UpdateMoviePictures(ctx context.Context,
	in *admin_movies_service.UpdateMoviePicturesRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.UpdateMoviePictures")
	defer span.Finish()

	poster, preview, background, err := r.moviesRepo.GetPicturesIds(ctx, in.MovieID)
	if errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("movie with this id not found.%w", ErrNotFound)
	}
	if err != nil {
		return err
	}

	replaceParams := map[string]ReplacePicturesParam{
		"poster":        convertPictureConfigToReplaceParam(in.Poster, r.postersConfig, poster),
		"previewPoster": convertPictureConfigToReplaceParam(in.PreviewPoster, r.backgroundConfig, preview),
		"background":    convertPictureConfigToReplaceParam(in.Background, r.backgroundConfig, background),
	}

	ids, err := r.imagesService.ReplacePictures(ctx, replaceParams)
	if err != nil {
		return err
	}

	if err != nil {
		for key, id := range ids {
			err := r.imagesService.DeletePicture(ctx, replaceParams[key].Category, id)
			if err != nil {
				r.logger.Errorf("image with %s id not deleted err: %s", id, err.Error())
			}
		}
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return err
	}

	err = r.moviesRepo.UpdatePictures(ctx, in.MovieID, ids["poster"], ids["previewPoster"], ids["background"])
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func getUpdateParamFromProto(in *admin_movies_service.UpdateMovieRequest, ageRating int32) repository.UpdateMovieParam {
	return repository.UpdateMovieParam{
		TitleRU:          in.GetTitleRU(),
		TitleEN:          in.GetTitleEN(),
		Description:      in.GetDescription(),
		Duration:         in.GetDuration(),
		ShortDescription: in.GetShortDescription(),
		ReleaseYear:      in.GetReleaseYear(),
		AgeRating:        ageRating,
	}
}

func (r *moviesRepositoryWrapper) UpdateMovie(ctx context.Context,
	in *admin_movies_service.UpdateMovieRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.UpdateMovie")
	defer span.Finish()

	if err := validateUpdateRequest(in); err != nil {
		return ErrInvalidArgument
	}

	exists, err := r.moviesRepo.IsMovieExists(ctx, in.MovieID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	if in.UpdateCountries || in.UpdateGenres {
		err := r.existanceChecker.CheckExistance(ctx, in.CountriesIDs, in.GenresIDs)
		if err != nil {
			ext.LogError(span, err)
			span.SetTag("grpc.status", status.Code(err))
			return err
		}
	}

	var ageRating int32
	if in.GetAgeRating() != "" {
		ageRating, err = r.moviesRepo.GetAgeRating(ctx, in.GetAgeRating())
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("can't find age rating with name: %d.%w", in.AgeRating, ErrInvalidArgument)
		}
		if err != nil {
			return fmt.Errorf("can't get age rating id. %w" + err.Error())
		}
	}

	r.logger.Info("Updating movie")
	err = r.moviesRepo.UpdateMovie(ctx, in.MovieID,
		getUpdateParamFromProto(in, ageRating), in.GenresIDs, in.CountriesIDs, in.UpdateGenres, in.UpdateCountries)
	if err != nil {
		return err
	}

	return nil
}

func (r *moviesRepositoryWrapper) DeleteMovie(ctx context.Context, in *admin_movies_service.DeleteMovieRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepositoryWrapper.DeleteMovie")
	defer span.Finish()

	poster, preview, background, err := r.moviesRepo.GetPicturesIds(ctx, in.MovieID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	pictures := []struct{ category, name string }{
		struct {
			category string
			name     string
		}{r.postersConfig.Category, poster},
		struct {
			category string
			name     string
		}{r.previewPostersConfig.Category, preview},
		struct {
			category string
			name     string
		}{r.backgroundConfig.Category, background},
	}

	for _, picture := range pictures {
		err = r.imagesService.DeletePicture(ctx, picture.category, picture.name)
		if err != nil {
			return err
		}
	}

	err = r.moviesRepo.DeleteMovie(ctx, in.MovieID)
	if err != nil {
		return err
	}

	span.SetTag("grpc.status", codes.OK)
	return nil
}
