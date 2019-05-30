package appsflyer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	af := AppsFlyer{}
	assert.Error(t, af.SendEvent(&Event{}))
}
