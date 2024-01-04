package mock_service

import (
	context "context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

type checkerResource interface {
	check(t *testing.T)
}

type checkerMockController struct {
	t         *testing.T
	logger    *logrus.Logger
	resources []checkerResource
}

func (c *checkerMockController) Check() {
	for _, r := range c.resources {
		r.check(c.t)
	}
}

func (c *checkerMockController) addResource(r checkerResource) {
	c.resources = append(c.resources, r)
}

func NewCheckerController(t *testing.T, l *logrus.Logger) *checkerMockController {
	if l == nil {
		l, _ = test.NewNullLogger()
	}
	assert.NotNil(t, l)
	return &checkerMockController{t: t, logger: l}
}

type CheckerBehavior struct {
	Expected struct {
		IDs       []int32
		CallTimes int32
	}

	Returns struct {
		Exists       bool
		NotExistsIDs []int32
		Err          error
	}

	SleepTime time.Duration
}

type geoCheckerMock struct {
	behavior  CheckerBehavior
	callTimes int32
	t         *testing.T
	logger    *logrus.Logger
}

func NewGeoCheckerMock(behavior CheckerBehavior, controller *checkerMockController) *geoCheckerMock {
	mock := &geoCheckerMock{
		behavior:  behavior,
		callTimes: 0,
		logger:    controller.logger,
	}

	controller.addResource(mock)
	mock.t = controller.t
	return mock
}

func (m *geoCheckerMock) check(t *testing.T) {
	assert.Equal(m.t, m.behavior.Expected.CallTimes, m.callTimes, "calls doesn't match")
}

func (m *geoCheckerMock) IsCountriesExists(ctx context.Context,
	ids []int32) (exists bool, notExistsIDs []int32, err error) {
	m.callTimes++
	assert.LessOrEqual(m.t, m.behavior.Expected.CallTimes, m.callTimes, "unexpected call")
	assert.Equal(m.t, m.behavior.Expected.IDs, ids, "unexpected ids")

	select {
	case <-ctx.Done():
		m.logger.Debug("geo checker context canceled")
		return m.behavior.Returns.Exists, m.behavior.Returns.NotExistsIDs, m.behavior.Returns.Err
	case <-time.After(m.behavior.SleepTime):
		return m.behavior.Returns.Exists, m.behavior.Returns.NotExistsIDs, m.behavior.Returns.Err
	}
}

type genresCheckerMock struct {
	behavior  CheckerBehavior
	callTimes int32
	t         *testing.T
	logger    *logrus.Logger
}

func NewGenresCheckerMock(behavior CheckerBehavior, controller *checkerMockController) *genresCheckerMock {
	mock := &genresCheckerMock{
		behavior:  behavior,
		callTimes: 0,
		logger:    controller.logger,
	}

	controller.addResource(mock)
	mock.t = controller.t
	return mock
}

func (m *genresCheckerMock) IsGenresExists(ctx context.Context,
	ids []int32) (exists bool, notExistsIDs []int32, err error) {

	m.callTimes++
	assert.LessOrEqual(m.t, m.behavior.Expected.CallTimes, m.callTimes, "unexpected call")
	assert.Equal(m.t, m.behavior.Expected.IDs, ids, "unexpected ids")

	select {
	case <-ctx.Done():
		m.logger.Debug("genres checker context canceled")
		return m.behavior.Returns.Exists, m.behavior.Returns.NotExistsIDs, m.behavior.Returns.Err
	case <-time.After(m.behavior.SleepTime):
		return m.behavior.Returns.Exists, m.behavior.Returns.NotExistsIDs, m.behavior.Returns.Err
	}
}

func (m *genresCheckerMock) check(t *testing.T) {
	assert.Equal(m.t, m.behavior.Expected.CallTimes, m.callTimes, "calls doesn't match")
}

type personsCheckerMock struct {
	behavior  CheckerBehavior
	callTimes int32
	t         *testing.T
	logger    *logrus.Logger
}

func NewPersonsCheckerMock(behavior CheckerBehavior, controller *checkerMockController) *personsCheckerMock {
	mock := &personsCheckerMock{
		behavior:  behavior,
		callTimes: 0,
		logger:    controller.logger,
	}

	controller.addResource(mock)
	mock.t = controller.t
	return mock
}

func (m *personsCheckerMock) IsPersonsExists(ctx context.Context,
	ids []int32) (exists bool, notExistsIDs []int32, err error) {
	m.callTimes++
	assert.LessOrEqual(m.t, m.behavior.Expected.CallTimes, m.callTimes, "unexpected call")
	assert.Equal(m.t, m.behavior.Expected.IDs, ids, "unexpected ids")

	select {
	case <-ctx.Done():
		m.logger.Debug("persons checker context canceled")
		return m.behavior.Returns.Exists, m.behavior.Returns.NotExistsIDs, m.behavior.Returns.Err
	case <-time.After(m.behavior.SleepTime):
		return m.behavior.Returns.Exists, m.behavior.Returns.NotExistsIDs, m.behavior.Returns.Err
	}
}

func (m *personsCheckerMock) check(t *testing.T) {
	assert.Equal(m.t, m.behavior.Expected.CallTimes, m.callTimes, "calls doesn't match")
}
