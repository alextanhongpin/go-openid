package testdata

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/alextanhongpin/go-openid"
)

type clientHelper struct {
	mock.Mock
	duration time.Duration
	now      time.Time
}

func NewClientHelper(duration time.Duration, now time.Time) *clientHelper {
	return &clientHelper{
		duration: duration,
		now:      now,
	}
}

func (c *clientHelper) NewDuration() time.Duration {
	return c.duration
}

func (c *clientHelper) NewTime() time.Time {
	return c.now
}

func (c *clientHelper) NewClientID() string {
	args := c.Called()
	return args.String(0)
}

func (c *clientHelper) NewClientSecret(clientID string, duration time.Duration) string {
	args := c.Called(clientID, duration)
	return args.String(0)
}

type clientValidator struct {
	mock.Mock
}

func NewClientValidator() *clientValidator {
	return &clientValidator{}
}

func (c *clientValidator) Validate(client *oidc.Client) error {
	args := c.Called(client)
	return args.Error(0)
}
