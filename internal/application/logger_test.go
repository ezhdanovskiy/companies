package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger_ProductionLevel(t *testing.T) {
	logger, err := newLogger("info", "json")
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewLogger_DebugLevel(t *testing.T) {
	logger, err := newLogger("debug", "json")
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewLogger_DebugLevelCaseInsensitive(t *testing.T) {
	logger, err := newLogger("DEBUG", "json")
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewLogger_ConsoleEncoding(t *testing.T) {
	logger, err := newLogger("info", "console")
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewLogger_EmptyEncoding(t *testing.T) {
	logger, err := newLogger("info", "")
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewLogger_DefaultLevel(t *testing.T) {
	logger, err := newLogger("warn", "json")
	require.NoError(t, err)
	assert.NotNil(t, logger)
}
