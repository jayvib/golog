package golog

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLogrus_Print(t *testing.T) {
	t.Run("Global state level is lower than current level", func(t *testing.T){
		var out bytes.Buffer
		SetLevel(DebugLevel)
		l := NewLogrusLogger(InfoLevel)
		l.logger.SetOutput(&out)

		l.Print("hello world")
		assert.True(t, strings.Contains(out.String(), "hello world"))
		assert.True(t, strings.Contains(out.String(), "info"))
	})

	t.Run("Global state level is higher than the current level", func(t *testing.T){
		var out bytes.Buffer
		SetLevel(DisabledLevel)
		l := NewLogrusLogger(InfoLevel)
		l.logger.SetOutput(&out)
		l.Print("empty")
		assert.Empty(t, out.String())
	})
}

