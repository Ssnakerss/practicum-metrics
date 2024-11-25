package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPadding(t *testing.T) {
	s1 := "test"
	s2 := "test"
	len := 100
	pad := "_"
	PadRight(&s1, pad, int64(len))
	PadRight2(&s2, pad, len)
	require.Equal(t, s1, s2)
}
