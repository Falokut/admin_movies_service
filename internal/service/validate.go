package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
)

var ErrInvalidFilter = errors.New("invalid filter value, filter must contain only digits and commas")

func validateFilter(filter *admin_movies_service.GetMoviesRequest) error {
	if filter.GetGenresIDs() != "" {
		if err := checkParam(*filter.GenresIDs); err != nil {
			return err
		}
	}
	if filter.GetCountriesIDs() != "" {
		if err := checkParam(*filter.CountriesIDs); err != nil {
			return err
		}
	}
	if filter.GetMoviesIDs() != "" {
		if err := checkParam(*filter.MoviesIDs); err != nil {
			return err
		}
	}
	if filter.GetAgeRatings() != "" {
		if strings.Contains(*filter.AgeRatings, "'") {
			return ErrInvalidFilter
		}
	}
	return nil
}

func validateCreateRequest(in *admin_movies_service.CreateMovieRequest) error {
	if in.Description == "" {
		return errors.New("description mustn't be empty")
	}

	if in.ShortDescription == "" {
		return errors.New("short description mustn't be empty")
	}

	if in.TitleRU == "" {
		return errors.New("titleRU mustn't be empty")
	}

	if len(in.GenresIDs) == 0 {
		return errors.New("genres ids mustn't be empty")
	}

	if in.Duration <= 0 {
		return errors.New("movie duration must be bigger than zero")
	}

	if len(in.Poster) == 0 {
		return errors.New("movie must have poster")
	}

	if len(in.PreviewPoster) == 0 {
		return errors.New("movie must have preview poster")
	}

	if in.ReleaseYear < 1700 || in.ReleaseYear > int32(time.Now().Year()) {
		return errors.New("movie release year must be bigger than 1700 and not to be more than this year")
	}

	if in.AgeRating == "" {
		return errors.New("age rating mustn't be emtpy")
	}

	return nil
}

func validateUpdateRequest(in *admin_movies_service.UpdateMovieRequest) error {
	if in.GetDuration() < 0 {
		return errors.New("movie duration must be bigger than zero")
	}

	if in.GetReleaseYear() != 0 && (in.GetReleaseYear() < 1700 || in.GetReleaseYear() > int32(time.Now().Year())) {
		return errors.New("movie release year must be bigger than 1700 and not to be more than this year")
	}

	return nil
}

func checkParam(val string) error {
	exp := regexp.MustCompile("^[!-&!+,0-9]+$")
	if !exp.Match([]byte(val)) {
		return ErrInvalidFilter
	}

	return nil
}
