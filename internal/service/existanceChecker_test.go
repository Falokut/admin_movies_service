package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Falokut/admin_movies_service/internal/service"
	mock_service "github.com/Falokut/admin_movies_service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var defaultCheckerBehavior = mock_service.CheckerBehavior{
	Returns: struct {
		Exists       bool
		NotExistsIDs []int32
		Err          error
	}{Exists: true, NotExistsIDs: []int32{}, Err: nil},
	Expected: struct {
		IDs       []int32
		CallTimes int32
	}{CallTimes: 1},
	SleepTime: time.Millisecond,
}

func TestCheckExistance(t *testing.T) {
	defer goleak.VerifyNone(t)
	type Args struct {
		directorsIDs []int32
		countriesIDs []int32
		genresIDs    []int32
	}

	testCases := []struct {
		ExpectedCode      codes.Code
		args              Args
		CountriesBehavior mock_service.CheckerBehavior
		GenresBehavior    mock_service.CheckerBehavior
		msg               string
	}{
		{
			ExpectedCode: codes.OK,
			args: Args{
				directorsIDs: []int32{1, 2, 3},
				countriesIDs: []int32{6, 3, 4},
				genresIDs:    []int32{546, 1231, 5},
			},
			CountriesBehavior: defaultCheckerBehavior,
			GenresBehavior:    defaultCheckerBehavior,
			msg:               "Test case num %d, mustn't return error, if all checked without errors",
		},
		{
			ExpectedCode: codes.InvalidArgument,
			args: Args{
				directorsIDs: []int32{1, 2, 3},
				countriesIDs: []int32{1, 3, 4},
				genresIDs:    []int32{546, 1239, 5},
			},
			CountriesBehavior: mock_service.CheckerBehavior{
				Expected: struct {
					IDs       []int32
					CallTimes int32
				}{CallTimes: 1},
				SleepTime: defaultCheckerBehavior.SleepTime * 10,
				Returns: struct {
					Exists       bool
					NotExistsIDs []int32
					Err          error
				}{
					Exists: false, NotExistsIDs: []int32{1239}, Err: nil,
				},
			},
			GenresBehavior: defaultCheckerBehavior,
			msg:            "Test case num %d, must return invalid argument code, if country not found",
		},
		{
			ExpectedCode: codes.Internal,
			args: Args{
				directorsIDs: []int32{1, 2, 3},
				countriesIDs: []int32{1, 3, 4},
				genresIDs:    []int32{546, 1239, 5},
			},
			CountriesBehavior: mock_service.CheckerBehavior{
				Expected: struct {
					IDs       []int32
					CallTimes int32
				}{CallTimes: 1},
				SleepTime: defaultCheckerBehavior.SleepTime * 10,
				Returns: struct {
					Exists       bool
					NotExistsIDs []int32
					Err          error
				}{
					Exists: false, NotExistsIDs: []int32{}, Err: fmt.Errorf("some errors"),
				},
			},
			GenresBehavior: defaultCheckerBehavior,
			msg:            "Test case num %d, must return internal code, if geo checker returns error",
		},
		{
			ExpectedCode: codes.Internal,
			args: Args{
				directorsIDs: []int32{1654, 2, 343},
				countriesIDs: []int32{1, 3123, 4},
				genresIDs:    []int32{56, 1239, 5},
			},
			GenresBehavior: mock_service.CheckerBehavior{
				Expected: struct {
					IDs       []int32
					CallTimes int32
				}{CallTimes: 1},
				SleepTime: time.Nanosecond,
				Returns: struct {
					Exists       bool
					NotExistsIDs []int32
					Err          error
				}{
					Exists: false, NotExistsIDs: []int32{}, Err: fmt.Errorf("some errors"),
				},
			},
			CountriesBehavior: defaultCheckerBehavior,
			msg:               "Test case num %d, must return iternal code, if genres checker returns error",
		},
		{
			ExpectedCode: codes.InvalidArgument,
			args: Args{
				directorsIDs: []int32{1, 2, 3},
				countriesIDs: []int32{1, 3, 4},
				genresIDs:    []int32{46, 129, 523},
			},
			GenresBehavior: mock_service.CheckerBehavior{
				Expected: struct {
					IDs       []int32
					CallTimes int32
				}{CallTimes: 1},
				SleepTime: time.Nanosecond * 10,
				Returns: struct {
					Exists       bool
					NotExistsIDs []int32
					Err          error
				}{
					Exists: false, NotExistsIDs: []int32{4}, Err: nil,
				},
			},
			CountriesBehavior: defaultCheckerBehavior,
			msg:               "Test case num %d, must return invalid argument code, if genre not found",
		},
	}

	for i, testCase := range testCases {
		ctr := mock_service.NewCheckerController(t, nil)
		defer ctr.Check()

		testCase.CountriesBehavior.Expected.IDs = testCase.args.countriesIDs
		countriesService := mock_service.NewGeoCheckerMock(testCase.CountriesBehavior, ctr)

		testCase.GenresBehavior.Expected.IDs = testCase.args.genresIDs
		genresService := mock_service.NewGenresCheckerMock(testCase.GenresBehavior, ctr)
		checker := service.NewExistanceChecker(countriesService, genresService, getNullLogger())
		err := checker.CheckExistance(context.Background(), testCase.args.countriesIDs, testCase.args.genresIDs)

		msg := fmt.Sprint(testCase.msg, i+1)
		assert.Equal(t, status.Code(err), testCase.ExpectedCode, msg)
	}
}
