package glog

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	l := NewDefault()
	ctx := WithContext(context.Background(), l)
	nl := FromContext(ctx)
	require.NotNil(t, nl)
	require.True(t, reflect.DeepEqual(l, nl))

	// update logger
	l = l.WithLevel(InfoLevel)
	ctx = WithContext(ctx, l)
	nl = FromContext(ctx)
	require.NotNil(t, nl)
	require.True(t, reflect.DeepEqual(l, nl))

	// new logger
	l1 := NewDefault()
	ctx = WithContext(ctx, l1)
	nl = FromContext(ctx)
	require.False(t, reflect.DeepEqual(l, nl))
}

func TestFromContext(t *testing.T) {
	l := FromContext(context.Background())
	require.Nil(t, l)
}

func TestFromContextDefault(t *testing.T) {
	l := FromContextDefault(context.Background())
	require.NotNil(t, l)
}
