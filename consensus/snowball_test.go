package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Verify(t *testing.T) {
	invalidAlphal := Config{
		Alphal: 10,
		K:      3,
		Beta:   10,
	}
	err := invalidAlphal.Verify()
	assert.NotNil(t, err)

	invalidBeta := Config{
		K:      3,
		Alphal: 2,
		Beta:   0,
	}
	err = invalidBeta.Verify()
	assert.NotNil(t, err)

	validConfig := Config{
		K:      3,
		Alphal: 2,
		Beta:   10,
	}
	err = validConfig.Verify()
	assert.NoError(t, err)
}
