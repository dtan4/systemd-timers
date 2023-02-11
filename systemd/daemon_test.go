package systemd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchTimer(t *testing.T) {
	assert.True(t, matchTimer("system-test-01", []string{"system-test-*"}))
	assert.True(t, matchTimer("system-test-01", []string{"system-*-01"}))
	assert.True(t, matchTimer("system-test-01", []string{"system-*-01*", "system-test-02"}))
	assert.False(t, matchTimer("system-test-01", []string{"system-test-1"}))
	assert.True(t, matchTimer("system-test-01", []string{"system-test-1", "system-test-01"}))
	assert.False(t, matchTimer("", []string{"system-test-1", "system-test-01"}))
}
