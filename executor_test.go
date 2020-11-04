package glog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultipleExecutor(t *testing.T) {
	var b1, b2 bytes.Buffer

	e1 := MatchExecutor(&b1, MatchGELevel(DebugLevel))
	e2 := MatchExecutor(&b2, MatchGELevel(ErrorLevel))

	l := NewDefault().WithExecutor(MultipleExecutor(e1, e2))

	l.Debug().Msg("DebugMessage").Fire()
	l.Info().Msg("InfoMessage").Fire()
	l.Error().Msg("ErrorMessage").Fire()
	l.Fatal().Msg("FatalMessage").Fire()

	s1 := b1.String()
	require.Contains(t, s1, "DebugMessage")
	require.Contains(t, s1, "InfoMessage")
	require.Contains(t, s1, "ErrorMessage")
	require.Contains(t, s1, "FatalMessage")

	s2 := b2.String()
	require.NotContains(t, s2, "DebugMessage")
	require.NotContains(t, s2, "InfoMessage")
	require.Contains(t, s2, "ErrorMessage")
	require.Contains(t, s2, "FatalMessage")
}
