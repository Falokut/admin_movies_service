package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type GeoChecker interface {
	IsCountriesExists(ctx context.Context, ids []int32) (exists bool, notExistsIDs []int32, err error)
}

type GenresChecker interface {
	IsGenresExists(ctx context.Context, ids []int32) (exists bool, notExistsIDs []int32, err error)
}

type existanceChecker struct {
	geoCheck     GeoChecker
	genresCheck  GenresChecker
	errorHandler errorHandler
	logger       *logrus.Logger
}

func NewExistanceChecker(geoCheck GeoChecker, genresCheck GenresChecker, logger *logrus.Logger) *existanceChecker {
	errorHandler := newErrorHandler(logger)
	return &existanceChecker{
		logger:       logger,
		errorHandler: errorHandler,
		geoCheck:     geoCheck,
		genresCheck:  genresCheck,
	}
}

func convertIntSliceIntoString(nums []int32) string {
	var str = make([]string, len(nums))
	for i, num := range nums {
		str[i] = fmt.Sprint(num)
	}
	return strings.Join(str, ",")
}

func (c *existanceChecker) CheckExistance(ctx context.Context, countriesIDs, genresIDs []int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "existanceChecker.CheckExistance")
	defer span.Finish()
	if len(countriesIDs) == 0 && len(genresIDs) == 0 {
		return nil
	}

	c.logger.Info("start checking existance")
	var wg sync.WaitGroup

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var errCh = make(chan error, 2)
	f := func(f func(context.Context, []int32) (bool, []int32, error),
		ids []int32, iternalMsg, notfoundFormat string) {
		defer wg.Done()

		select {
		case <-reqCtx.Done():
			return
		default:
			exists, notFound, err := f(reqCtx, ids)
			if err != nil {
				errCh <- c.errorHandler.createErrorResponceWithSpan(span, ErrInternal, iternalMsg+" "+err.Error())
			} else if !exists {
				errCh <- c.errorHandler.createExtendedErrorResponceWithSpan(span, ErrInvalidArgument, "",
					fmt.Sprintf(notfoundFormat, convertIntSliceIntoString(notFound)))
			}
		}
	}

	if len(countriesIDs) > 0 {
		wg.Add(1)
		go f(c.geoCheck.IsCountriesExists, countriesIDs,
			"can't check countries existence", "countries with ids %s not found")
	}
	if len(genresIDs) > 0 {
		wg.Add(1)
		go f(c.genresCheck.IsGenresExists, genresIDs,
			"can't check genres existence", "genres with ids %s not found")

	}

	go func() {
		wg.Wait()
		c.logger.Debug("channel closed")
		close(errCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, open := <-errCh:
			if !open {
				return nil
			} else if err != nil {
				return err
			}
		}
	}
}
