package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	assert := require.New(t)
	_, err := Validate(map[string]interface{}{})
	assert.Equal(err, nil)
}
