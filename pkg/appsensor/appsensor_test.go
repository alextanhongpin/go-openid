package appsensor_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/stretchr/testify/assert"
)

func TestAppSensor(t *testing.T) {
	assert := assert.New(t)

	aps := appsensor.NewLoginDetector()
	id := "john@mail.com"
	aps.Increment(id)
	aps.Increment(id)
	attempt := aps.Stat(id)
	assert.Equal(int64(2), attempt.Count, "should have 2 attempts")

	aps.Increment(id)
	aps.Increment(id)

	attempt2 := aps.Stat(id)
	assert.Equal(int64(4), attempt.Count, "should have 4 attempts")
	assert.Equal(int64(4), attempt2.Count, "should have 4 attempts")
	assert.Equal(attempt, attempt2, "should have the same memory address")

	locked := aps.IsLocked(id)
	assert.True(locked)
}
